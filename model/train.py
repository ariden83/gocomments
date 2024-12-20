import os
import time
import math
import random
import datetime
from pathlib import Path
import tensorflow as tf
import glob

os.environ["TF_CPP_MIN_LOG_LEVEL"] = "1"  # reduce the amount of console output from TF
os.environ['CUDA_VISIBLE_DEVICES'] = "0"

from transformers import *
from datasets import load_dataset

logging.set_verbosity_warning()
logging.set_verbosity_error()

import logging

print('TF version', tf.__version__)
print("Num GPUs Available: ", len(tf.config.list_physical_devices('GPU'))) # check GPU available

from matplotlib.pylab import plt

def setup_strategy(xla, fp16, no_cuda):
    print(" Tensorflow: setting up strategy")

    if xla:
        tf.config.optimizer.set_jit(True)
        print("XLA Enabled")
    if fp16:
        policy = tf.keras.mixed_precision.Policy('mixed_float16')
        tf.keras.mixed_precision.set_global_policy(policy)
        print("Mixed Precision Training Enabled")

    # setup distribution strategy
    gpus = tf.config.list_physical_devices("GPU")
    if no_cuda:
        strategy = tf.distribute.OneDeviceStrategy(device="/cpu:0")
    else:
        if len(gpus) == 0:
            print(" One Device Strategy [CPU] Enabled")
            strategy = tf.distribute.OneDeviceStrategy(device="/cpu:0")
        elif len(gpus) == 1:
            print(" One Device Strategy [GPU] Enabled")
            strategy = tf.distribute.OneDeviceStrategy(device="/gpu:0")
        elif len(gpus) > 1:
            print(" Mirrored Strategy Enabled")
            # If only want to use a specific subset of GPUs use CUDA_VISIBLE_DEVICES=0`
            strategy = tf.distribute.MirroredStrategy()
        else:
            strategy = tf.distribute.get_strategy()

    return strategy

def check_gpu_memory():
    gpus = tf.config.experimental.list_physical_devices('GPU')
    if gpus:
        try:
            for gpu in gpus:
                details = tf.config.experimental.get_memory_info(gpu)
                print(f"GPU {gpu}: {details['current']} MB used out of {details['peak']} MB")
        except RuntimeError as e:
            print(e)

def n_replicas(strategy):
    # return number of devices
    return strategy.num_replicas_in_sync

def get_current_checkpoint_epoch(ckpt_dir):
    checkpoint_files = [f for f in os.listdir(ckpt_dir) if f.startswith('ckpt') and f.endswith('.index')]
    if not checkpoint_files:
        return 0
    latest_ckpt = max(checkpoint_files, key=lambda x: int(x.split('-')[1].split('.')[0]))
    current_epoch = int(latest_ckpt.split('-')[1].split('.')[0])
    return current_epoch

def adjust_epochs(args, current_epoch):
    args.epochs = args.epochs - current_epoch
    print(f"Adjusted number of epochs to: {args.epochs}")


# note:
# huggingface TF-T5 implementation has issues when mixed precision is enabled
# we will disable FP16 for this but can be used for training any other model
strategy = setup_strategy(xla=True, fp16=False, no_cuda=False)
check_gpu_memory()


def download_dataset(cache_dir):
    # download data using a keras utility
    _url = "https://raw.githubusercontent.com/google-research/google-research/master/mbpp/mbpp.jsonl" # download mbpp dataset
    dataset_path = tf.keras.utils.get_file("mbpp.jsonl", origin=_url, cache_dir=cache_dir, cache_subdir=cache_dir)
    return dataset_path


def download_local_dataset(cache_dir):
    # Spécifiez le répertoire où les fichiers sont stockés
    dataset_directory = "/app/scripts/dataset"

    # Utilisez glob pour lister tous les fichiers .jsonl dans le répertoire
    jsonl_files = glob.glob(os.path.join(dataset_directory, "*.jsonl"))

    # Vérifiez s'il y a des fichiers .jsonl
    if not jsonl_files:
        raise FileNotFoundError(f"No .jsonl files found in directory {dataset_directory}.")

    # Trouvez le fichier le plus récent basé sur le temps de modification
    latest_file = max(jsonl_files, key=os.path.getmtime)

    print(f"Using local file: {latest_file}")
    return latest_file


def convert_examples_to_features(examples, tokenizer, args):
    # encode text-code pairs
    names = examples['name']
    file_name = examples['file_name']
    comments = examples['comment']
    analyzed_comments = examples['comment_analysis']
    # tests = [" ".join(test) for test in examples['test_list']] # convert list of test cases to single string

    # encode names by prepending the task for input sequence
    # inputs = [text for text in names]
    inputs = [f"{file_name} {name}" for file_name, name in zip(file_name, names)]
    model_inputs = tokenizer(inputs, max_length=args.max_input_length, padding="max_length", truncation=True)

    # encode names by prepending the task for input sequence and appending the test sequence
    # inputs = [text + " " + test for text, test in zip(names, tests)]
    # model_inputs = tokenizer(inputs, max_length=args.max_input_length, padding="max_length", truncation=True)


    # Combinaison des commentaires analysés avec les commentaires originaux pour la séquence cible
    targets = [f"{comment} POS: {analyzed_comment['pos_tags']}"
               for comment, analyzed_comment in zip(comments, analyzed_comments)]
    labels = tokenizer(targets, max_length=args.max_target_length, padding="max_length", truncation=True).input_ids

    # encode names by prepending the task for input sequence
    # labels = tokenizer(comments, max_length=args.max_target_length, padding="max_length", truncation=True).input_ids

    # we need to replace the index of the padding tokens by -100
    # such that they are not taken into account by the CrossEntropyLoss
    labels_with_ignore_index = []
    for labels_example in labels:
        labels_example = [label if label != 0 else -100 for label in labels_example]
        labels_with_ignore_index.append(labels_example)
    model_inputs["labels"] = labels_with_ignore_index

    # return features
    return model_inputs


def get_train_tfdataset(train_dataset, num_train_examples, args):
    # select feature columns
    columns = ['input_ids', 'attention_mask', 'labels']
    # set to tensorflow format
    train_dataset.set_format(type='tensorflow', columns=columns)

    # specify return types
    return_types = {'input_ids':tf.int32, 'attention_mask':tf.int32, 'labels':tf.int32}
    # specify return shapes
    return_shapes = {'input_ids': tf.TensorShape([None]),'attention_mask': tf.TensorShape([None]), 'labels': tf.TensorShape([None])}
    # initialize dataset
    tf_dataset = tf.data.Dataset.from_generator(lambda : train_dataset, return_types, return_shapes)

    # turn off auto-sharding
    options = tf.data.Options()
    options.experimental_distribute.auto_shard_policy = tf.data.experimental.AutoShardPolicy.OFF
    tf_dataset = tf_dataset.with_options(options)

    # repeat, shuffle, batch, prefetch
    ds = (
        tf_dataset.repeat()
        .shuffle(num_train_examples, seed=args.seed)
        .batch(args.train_batch_size)
        .prefetch(tf.data.AUTOTUNE)
    )

    # distribute dataset to devices
    return strategy.experimental_distribute_dataset(ds)

def get_validation_tfdataset(eval_dataset, num_validation_examples, args):
    # select feature columns
    columns = ['input_ids', 'attention_mask', 'labels']
    # set to tensorflow format
    eval_dataset.set_format(type='tensorflow', columns=columns)

    # specify return types
    return_types = {'input_ids':tf.int32, 'attention_mask':tf.int32, 'labels':tf.int32}
    # specify return shapes
    return_shapes = {'input_ids': tf.TensorShape([None]),'attention_mask': tf.TensorShape([None]), 'labels': tf.TensorShape([None])}
    # initialize dataset
    tf_dataset = tf.data.Dataset.from_generator(lambda : eval_dataset, return_types, return_shapes)

    # turn off auto-sharding
    options = tf.data.Options()
    options.experimental_distribute.auto_shard_policy = tf.data.experimental.AutoShardPolicy.OFF
    tf_dataset = tf_dataset.with_options(options)

    # repeat, batch, prefetch
    ds = (
        tf_dataset.repeat()
        .batch(args.validation_batch_size)
        .prefetch(tf.data.AUTOTUNE)
    )

    # distribute dataset to devices
    return strategy.experimental_distribute_dataset(ds)

def fix_all_seeds(seed):
    # set random seed
    os.environ['PYTHONHASHSEED'] = str(seed)
    random.seed(seed)
    tf.random.set_seed(seed)

def init_logger(log_file=None, log_file_level=logging.NOTSET):
    # initialize logger for tracking events and save in file
    if isinstance(log_file, Path):
        log_file = str(log_file)
    log_format = logging.Formatter(
        fmt='%(asctime)s - %(levelname)s - %(name)s -   %(message)s',
        datefmt='%m/%d/%Y %H:%M:%S'
    )
    logger = logging.getLogger()
    logger.setLevel(logging.INFO)
    console_handler = logging.StreamHandler()
    console_handler.setFormatter(log_format)
    logger.handlers = [console_handler]
    if log_file and log_file != '':
        file_handler = logging.FileHandler(log_file)
        file_handler.setLevel(log_file_level)
        # file_handler.setFormatter(log_format)
        logger.addHandler(file_handler)
    return logger

class ProgressBar(object):
    # custom progress bar
    def __init__(self, n_total,width=30,desc = 'Training'):
        self.width = width
        self.n_total = n_total
        self.start_time = time.time()
        self.desc = desc

    def __call__(self, step, info={}):
        now = time.time()
        current = step + 1
        recv_per = current / self.n_total
        bar = f'[{self.desc}] {current}/{self.n_total} ['
        if recv_per >= 1:
            recv_per = 1
        prog_width = int(self.width * recv_per)
        if prog_width > 0:
            bar += '=' * (prog_width - 1)
            if current < self.n_total:
                bar += ">"
            else:
                bar += '='
        bar += '.' * (self.width - prog_width)
        bar += ']'
        show_bar = f"\r{bar}"
        time_per_unit = (now - self.start_time) / current
        if current < self.n_total:
            eta = time_per_unit * (self.n_total - current)
            if eta > 3600:
                eta_format = ('%d:%02d:%02d' %
                              (eta // 3600, (eta % 3600) // 60, eta % 60))
            elif eta > 60:
                eta_format = '%d:%02d' % (eta // 60, eta % 60)
            else:
                eta_format = '%ds' % eta
            time_info = f' - ETA: {eta_format}'
        else:
            if time_per_unit >= 1:
                time_info = f' {time_per_unit:.1f}s/step'
            elif time_per_unit >= 1e-3:
                time_info = f' {time_per_unit * 1e3:.1f}ms/step'
            else:
                time_info = f' {time_per_unit * 1e6:.1f}us/step'

        show_bar += time_info
        if len(info) != 0:
            show_info = f'{show_bar} ' + \
                        "-".join([f' {key}: {value:.4f} ' if key != "learning_rate" else f' {key}: {value:.8f} ' for key, value in info.items()])
            print(show_info, end='')
        else:
            print(show_bar, end='')

class Trainer:
    def __init__(
            self, model, args, train_dataset, validation_dataset,
            num_train_examples, num_validation_examples
    ):
        self.model = model
        self.args = args

        self.train_dataset = train_dataset
        self.num_train_examples = num_train_examples

        self.validation_dataset = validation_dataset
        self.num_validation_examples = num_validation_examples

        self.global_step = 0
        self.eval_loss = tf.keras.metrics.Sum()
        self.train_loss_dict = {}
        self.val_loss_dict = []

    def create_optimizer_and_scheduler(self, num_training_steps):
        # creates an optimizer with a learning rate schedule using a warmup phase followed by a linear decay.
        num_warmup_steps = math.ceil(num_training_steps * self.args.warmup_ratio)
        self.optimizer, self.lr_scheduler = create_optimizer(
            init_lr=self.args.learning_rate,
            num_train_steps=num_training_steps,
            num_warmup_steps=num_warmup_steps,
            weight_decay_rate=self.args.weight_decay,
            adam_epsilon=self.args.adam_epsilon
        )

    def evaluation_step(self, features, labels, nb_instances_in_global_batch):
        # forward pass
        outputs = self.model(input_ids=features['input_ids'], attention_mask=features['attention_mask'], labels=labels, training=False)[:2]
        loss, logits = outputs[:2]
        # loss scaling
        scaled_loss = loss / tf.cast(nb_instances_in_global_batch, dtype=loss.dtype)
        # add current batch loss
        self.eval_loss.update_state(scaled_loss)

    @tf.function
    def distributed_evaluation_steps(self, batch):
        features = {k: v for k, v in batch.items() if 'labels' not in k}
        labels = batch['labels']
        nb_instances = tf.reduce_sum(tf.cast(labels != -100, dtype=tf.int32))
        # strategy.run() expects args to be a list or tuple
        inputs = (features, labels, nb_instances)
        # `run` replicates the provided computation and runs with the distributed input
        strategy.run(self.evaluation_step, inputs)

    def evaluate(self):
        # calculate total validation steps
        steps = math.ceil(self.num_validation_examples / self.args.validation_batch_size)
        # reset eval loss after every epoch
        self.eval_loss.reset_states()
        logs = {}
        pbar = ProgressBar(n_total=steps, desc='Evaluating')
        # iterate over validation dataset
        for step, batch in enumerate(self.validation_dataset):
            # distributed evaluation step
            self.distributed_evaluation_steps(batch)
            logs["eval_loss"] = self.eval_loss.result() / (step + 1)
            pbar(step=step, info=logs)
            if step == steps - 1:
                break

        self.val_loss_dict.append(self.eval_loss.result() / (step + 1))
        print("\n------------- validation result -----------------")

    def apply_gradients(self, features, labels, nb_instances_in_global_batch):
        # forward pass
        outputs = self.model(input_ids=features['input_ids'], attention_mask=features['attention_mask'], labels=labels, training=True)[:2]
        loss, logits = outputs[:2]
        # loss scaling
        scaled_loss = loss / tf.cast(nb_instances_in_global_batch, dtype=loss.dtype)
        # calculate gradients
        gradients = tf.gradients(scaled_loss, self.model.trainable_variables)
        # convert gradients with nan value
        gradients = [g if g is not None else tf.zeros_like(v) for g, v in zip(gradients, self.model.trainable_variables)]
        # optimize the model
        self.optimizer.apply_gradients(list(zip(gradients, self.model.trainable_variables)))
        # add current batch loss
        self.train_loss.update_state(scaled_loss)

    @tf.function
    def distributed_training_steps(self, batch):
        with strategy.scope():
            features = {k: v for k, v in batch.items() if 'labels' not in k}
            labels = batch['labels']
            nb_instances = tf.reduce_sum(tf.cast(labels != -100, dtype=tf.int32))
            # strategy.run() expects args to be a list or tuple
            inputs = (features, labels, nb_instances)
            # `run` replicates the provided computation and runs with the distributed input.
            strategy.run(self.apply_gradients, inputs)

    def train(self):
        # calculate total training steps
        num_updates_per_epoch = self.num_train_examples // args.train_batch_size
        self.steps_per_epoch = num_updates_per_epoch
        t_total = self.steps_per_epoch * self.args.epochs

        with strategy.scope():
            # optimizer, and checkpoint must be created under `strategy.scope`
            # create optimizer and scheduler
            self.create_optimizer_and_scheduler(num_training_steps=t_total)

            # create checkpoint manager
            folder = os.path.join(self.args.output_dir, self.args.checkpoint_dir)
            ckpt = tf.train.Checkpoint(optimizer=self.optimizer, model=self.model)
            self.model.ckpt_manager = tf.train.CheckpointManager(ckpt, folder, max_to_keep=1)

            # restore checkpoint if available
            if self.model.ckpt_manager.latest_checkpoint:
                ckpt.restore(self.model.ckpt_manager.latest_checkpoint)
                logger.info(f"Restored from checkpoint: {self.model.ckpt_manager.latest_checkpoint}")
            else:
                logger.info("Starting training from scratch")

            iterations = self.optimizer.iterations

            logger.info("***** Running training *****")
            logger.info(f"  Num examples = {self.num_train_examples}")
            logger.info(f"  Num Epochs = {self.args.epochs}")
            logger.info(f"  Total train batch size (w. parallel & distributed) = {self.args.train_batch_size * n_replicas(strategy)}")
            logger.info(f"  Steps per epoch = {self.steps_per_epoch}")
            logger.info(f"  Total optimization steps = {t_total}")

            self.train_loss = tf.keras.metrics.Sum(name="training_loss")
            start_time = datetime.datetime.now()
            for epoch_iter in range(self.args.epochs):
                # training loop
                logger.info(f"Epoch {epoch_iter + 1}/{self.args.epochs}")

                pbar = ProgressBar(n_total=self.steps_per_epoch, desc='Training')
                # iterate over training dataset
                for step, batch in enumerate(self.train_dataset):
                    # distributed training step
                    self.distributed_training_steps(batch)

                    self.global_step = iterations.numpy()
                    training_loss = self.train_loss.result() / (step + 1)

                    logger.info(f"Step : {step} / {self.steps_per_epoch}, Training loss: {training_loss.numpy()}, Learning rate: {self.lr_scheduler(self.global_step).numpy()}")

                    if self.global_step % self.steps_per_epoch == 0:
                        print("\n------------- train result -----------------")
                        # call to evaluation loop
                        self.evaluate()
                        # save checkpoint
                        ckpt_save_path = self.model.ckpt_manager.save()
                        logger.info(f"Saving checkpoint at {ckpt_save_path}")

                        try:
                            # Extract the checkpoint number from ckpt_save_path
                            checkpoint_number = ckpt_save_path.split('-')[-1]

                            model_save_dir = os.path.join(self.args.output_dir, self.args.checkpoint_model_dir, f"checkpoint-{checkpoint_number}")
                            os.makedirs(model_save_dir, exist_ok=True)
                            self.model.save_pretrained(model_save_dir)

                            tokenizer = RobertaTokenizer.from_pretrained(args.tokenizer_name)
                            tokenizer.save_pretrained(model_save_dir)

                            model_save_dir = os.path.join(self.args.output_dir, self.args.checkpoint_model_dir, f"checkpoint-tf-{checkpoint_number}")
                            os.makedirs(model_save_dir, exist_ok=True)
                            self.model.save(model_save_dir, save_format='tf')

                        except OSError as e:
                            print(f"An error occurred while creating the directory or saving the model/tokenizer: {e}")
                        except Exception as e:
                            print(f"An unexpected error occurred: {e}")
                        self.train_loss.reset_states()
                        break

                # reset train loss after every epoch
                self.train_loss_dict[epoch_iter] = training_loss.numpy()
                print(self.train_loss_dict[epoch_iter])
                self.train_loss.reset_states()


            plt.plot(range(self.args.epochs), self.train_loss_dict.values(), label='Training Loss')
            print(self.train_loss_dict.values())
            print(self.val_loss_dict)
            plt.plot(range(self.args.epochs), self.val_loss_dict, label='Validation Loss')

            # Add in a title and axes labels
            plt.title('Training and Validation Loss')
            plt.xlabel('Epochs')
            plt.ylabel('Loss')

            # Set the tick locations
            # plt.xticks(arange(0, 21, 2))

            # Display the plot
            plt.legend(loc='best')
            plt.show()

            plt.savefig(f"{self.args.output_dir}/{self.args.checkpoint_dir}/training_curves.jpg")

            end_time = datetime.datetime.now()
            logger.info("Training Time cost: %s", str(end_time - start_time))


def run(args):
    logger.info(" Starting training / evaluation")

    # logger.info(" Downloading Data Files")
    # dataset_path = download_dataset(args.cache_dir)

    dataset_path = download_local_dataset(args.cache_dir)

    logger.info(" Loading Data Files")
    dataset = load_dataset('json', data_files=dataset_path, cache_dir='/app/scripts/dataset/')
    # train test split
    dataset = dataset['train'].train_test_split(0.1, shuffle=False)

    logger.info(" Initializing Tokenizer")
    tokenizer = RobertaTokenizer.from_pretrained(args.tokenizer_name)

    logger.info(" Preparing Features")
    dataset = dataset.map(convert_examples_to_features, batched=True, fn_kwargs={"tokenizer": tokenizer, "args": args})

    logger.info(" Intializing training and validation dataset ")
    train_dataset = dataset['train']
    num_train_examples = len(dataset['train'])
    # create tf train dataset
    tf_train_dataset = get_train_tfdataset(train_dataset, num_train_examples, args)

    validation_dataset = dataset['test']
    num_validation_examples = len(dataset['test'])
    # create tf validation dataset
    tf_validation_dataset = get_validation_tfdataset(train_dataset, num_validation_examples, args)

    logger.info(f' Intializing model | {args.model_type.upper()} ')
    with strategy.scope():
        # model must be created under `strategy.scope`
        model = TFT5ForConditionalGeneration.from_pretrained(args.model_name_or_path, from_pt=True)

    # Adjust epochs based on the checkpoint
    current_epoch = get_current_checkpoint_epoch(os.path.join(args.output_dir, args.checkpoint_dir))
    adjust_epochs(args, current_epoch)

    # custom training loop
    trainer = Trainer(model, args, tf_train_dataset, tf_validation_dataset, num_train_examples, num_validation_examples)
    trainer.train()

    # save pretrained model and tokenizer
    logger.info(f" Saving model in {args.save_dir}")
    trainer.model.save_pretrained(args.save_dir)
    tokenizer.save_pretrained(args.save_dir)


class Args:
    # define training arguments

    # MODEL
    model_type = 't5'
    tokenizer_name = 'Salesforce/codet5-base'
    model_name_or_path = 'Salesforce/codet5-base'

    # DATA
    train_batch_size = 8
    validation_batch_size = 8
    max_input_length = 48
    max_target_length = 128

    # OPTIMIZER
    learning_rate = 3e-4
    weight_decay = 1e-4
    warmup_ratio = 0.2
    adam_epsilon = 1e-8

    # TRAINING
    seed = 2022
    epochs = 60

    # timestamp = datetime.strftime("%Y%m%d_%H")

    # DIRECTORIES
    output_dir = "runs/"
    logging_dir = f"{output_dir}/logs/"
    checkpoint_dir = f"checkpoint"
    checkpoint_model_dir = f"checkpoint_model"
    save_dir = f"{output_dir}/saved_model_2/"
    cache_dir = '../working/'
    Path(output_dir).mkdir(parents=True, exist_ok=True)
    Path(logging_dir).mkdir(parents=True, exist_ok=True)
    Path(save_dir).mkdir(parents=True, exist_ok=True)


def run_predict(args, text):
    # load saved finetuned model
    model = TFT5ForConditionalGeneration.from_pretrained(args.save_dir)
    # load saved tokenizer
    tokenizer = RobertaTokenizer.from_pretrained(args.save_dir)

    # encode texts by prepending the task for input sequence and appending the test sequence
    query = text
    encoded_text = tokenizer(query, return_tensors='tf', padding='max_length', truncation=True, max_length=args.max_input_length)

    # inference
    generated_code = model.generate(
        encoded_text["input_ids"], attention_mask=encoded_text["attention_mask"],
        max_length=args.max_target_length, top_p=0.95, top_k=50, repetition_penalty=2, num_return_sequences=1
    )

    # decode generated tokens
    decoded_code = tokenizer.decode(generated_code.numpy()[0], skip_special_tokens=True)
    return decoded_code


def predict_from_dataset(args):
    # load using hf datasets
    dataset = load_dataset('json', data_files='../working/mbpp.jsonl')
    # train test split
    dataset = dataset['train'].train_test_split(0.1, shuffle=False)
    test_dataset = dataset['test']

    # randomly select an index from the validation dataset
    index = random.randint(0, len(test_dataset))
    text = test_dataset[index]['text']
    code = test_dataset[index]['code']

    # run-predict on text
    decoded_code = run_predict(args, text)

    print("#" * 25); print("QUERY: ", text)
    print()
    print('#' * 25); print("ORIGINAL: ")
    print("\n", code)
    print()
    print('#' * 25); print("GENERATED: ")
    print("\n", decoded_code)


def predict_from_text(args, text):
    # run-predict on text
    decoded_code = run_predict(args, text)
    print("#" * 25); print("QUERY: ", text)
    print()
    print('#' * 25); print("GENERATED: ")
    print("\n", decoded_code)


# initialize training arguments
args = Args()
# initialize logger
logger = init_logger(log_file=os.path.join(args.logging_dir, f"{args.model_type}-{time.strftime('%Y-%m-%d-%H-%M-%S', time.localtime())}.log"))
# fix all seeds
fix_all_seeds(args.seed)

if __name__ == "__main__":
    # run training and evaluation
     run(args)
