import argparse
import logging
import os
from transformers import T5Tokenizer, RobertaTokenizer, T5ForConditionalGeneration

logger = logging.getLogger(__name__)

def parse_args():
    parser = argparse.ArgumentParser(description="Load a T5 model checkpoint and save it to a new directory.")
    parser.add_argument(
        "--checkpoint_dir",
        default="/workspace/runs",
        type=str,
        help="Path to the checkpoint directory.",
    )
    parser.add_argument(
        "--output_dir_prefix",
        default="/workspace/runs/model_from_checkpoint",
        type=str,
        help="Prefix for the output directory where the new model will be saved.",
    )
    args = parser.parse_args()
    return args

def main(args):
    logging.basicConfig(level=logging.INFO)
    logger.info(f"Loading model from checkpoint: {args.checkpoint_dir}")

    # Load the tokenizer and model from the checkpoint directory
    tokenizer = RobertaTokenizer.from_pretrained('Salesforce/codet5-base')
    # tokenizer = RobertaTokenizer.from_pretrained(args.checkpoint_dir)
    model = T5ForConditionalGeneration.from_pretrained(args.checkpoint_dir)

    # Determine the checkpoint number from the checkpoint directory name
    checkpoint_number = os.path.basename(args.checkpoint_dir).split('-')[-1]
    output_dir = f"{args.output_dir_prefix}-checkpoint-{checkpoint_number}"

    # Ensure the output directory exists
    os.makedirs(output_dir, exist_ok=True)

    # Save the model and tokenizer to the new directory
    logger.info(f"Saving model to new directory: {output_dir}")
    model.save_pretrained(output_dir)
    tokenizer.save_pretrained(output_dir)

    logger.info("Model and tokenizer saved successfully.")

if __name__ == "__main__":
    args = parse_args()
    main(args)
