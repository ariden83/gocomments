from flask import Flask, request, jsonify
from transformers import RobertaTokenizer, TFT5ForConditionalGeneration
import os
import time
from pathlib import Path
import logging

app = Flask(__name__)


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
    prefix = "Generate Python: "

    # OPTIMIZER
    learning_rate = 3e-4
    weight_decay = 1e-4
    warmup_ratio = 0.2
    adam_epsilon = 1e-8

    # TRAINING
    seed = 2022
    epochs = 20

    # DIRECTORIES
    output_dir = "/workspace/runs"
    logging_dir = f"{output_dir}/logs/"
    checkpoint_dir = f"checkpoint"
    save_dir = f"{output_dir}/checkpoint_model"
    cache_dir = '../working/'
    Path(output_dir).mkdir(parents=True, exist_ok=True)
    Path(logging_dir).mkdir(parents=True, exist_ok=True)
    Path(save_dir).mkdir(parents=True, exist_ok=True)


def tokenize(query, tokenizer_name, max_input_length):
    tokenizer = RobertaTokenizer.from_pretrained(tokenizer_name)
    return tokenizer(query, return_tensors='tf', padding='max_length', truncation=True, max_length=max_input_length)


def run_predict(args, query, version):

    checkpoint_path = os.path.join(args.save_dir, f"checkpoint-{version}")

    model = TFT5ForConditionalGeneration.from_pretrained(checkpoint_path)
    tokenizer = RobertaTokenizer.from_pretrained(checkpoint_path)

    encoded_text = tokenizer(query, return_tensors='tf', padding='max_length', truncation=True, max_length=args.max_input_length)

    # inference
    generated_code = model.generate(
        encoded_text["input_ids"], attention_mask=encoded_text["attention_mask"],
        max_length=args.max_target_length, top_p=0.95, top_k=50, repetition_penalty=2.0, num_return_sequences=1
    )

    # decode generated tokens
    decoded_code = tokenizer.decode(generated_code.numpy()[0], skip_special_tokens=True)
    return decoded_code


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


# initialize training arguments
args = Args()

@app.route('/ping', methods=['GET'])
def ping():
    return jsonify({'message': 'pong'}), 200

@app.route('/tokenize', methods=['POST'])
def tokenize_endpoint():
    data = request.json
    text = data.get('text')
    version = data.get('version', 9)

    if not text:
        return jsonify({'error': 'text and tokenizer_name are required'}), 400

    return jsonify({
        "comment": run_predict(args, text, version),
    })


if __name__ == '__main__':
    # initialize logger
    logger = init_logger(log_file=os.path.join(args.logging_dir, f"{args.model_type}-{time.strftime('%Y-%m-%d-%H-%M-%S', time.localtime())}.log"))
    app.run(host='0.0.0.0', port=5000)
