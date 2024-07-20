import requests
from pygments import lex
from pygments.lexers import GoLexer
from pygments.token import Token
from dotenv import load_dotenv
from pathlib import Path
from datetime import datetime
import json
import os
import subprocess
import spacy

# Charger le modèle de langue anglais
nlp = spacy.load("en_core_web_sm")

# Termes techniques communs en programmation
TECH_TERMS = {
    "initialize", "compute", "return", "calculate", "set", "get", "create",
    "update", "delete", "process", "validate", "parse", "fetch", "store",
    "modify", "load", "save", "convert", "render", "transform", "execute",
    "query", "sync", "merge", "send", "receive", "configure", "build",
    "handle", "optimize", "resolve", "map", "filter", "sort", "aggregate",
    "serialize", "deserialize", "log", "debug", "trace", "flush", "connect",
    "disconnect", "authenticate", "authorize"
}


def parse_go_file(file_path):
    """Parse a Go file to extract functions as JSON."""
    script_dir = os.path.dirname(os.path.abspath(__file__))
    try:
        result = subprocess.run(
            ['go', 'run', 'parse_functions.go', '--', file_path],
            cwd=script_dir,
            capture_output=True,
            text=True,
            check=True
        )
        parsed_output = json.loads(result.stdout)
        return parsed_output if parsed_output else []
    except subprocess.CalledProcessError as e:
        print(f"Failed to parse file {file_path}: {e.stderr}")
        return []
    except json.JSONDecodeError as e:
        print(f"JSON decode error: {e}")
        print(f"Failed to parse JSON from output: '{result.stdout}'")
        return []


def list_go_files(directory):
    """List Go files in a directory, excluding test and example files."""
    try:
        directory_path = Path(directory)
        if not directory_path.is_dir():
            raise NotADirectoryError(f"{directory} is not a directory")

        exclude_patterns = ['_test.go', '_example.go']
        go_files = [
            str(file) for file in directory_path.rglob('*.go')
            if not any(file.name.endswith(pattern) for pattern in exclude_patterns) and 'vendor' not in file.parts
        ]
        return go_files
    except (FileNotFoundError, NotADirectoryError, PermissionError) as e:
        print(f"Error: {e}")
        return []


def analyze_comment(comment):
    """Analyze a comment for quality and return detailed analysis."""
    doc = nlp(comment)
    analysis = {
        "tokens": [token.text for token in doc],
        "lemmas": [token.lemma_ for token in doc],
        "pos_tags": [token.pos_ for token in doc],
        "entities": [(ent.text, ent.label_) for ent in doc.ents],
        "quality_score": evaluate_quality(doc)
    }
    return analysis


def evaluate_quality(doc):
    """Evaluate the quality of a comment based on various criteria."""
    length_threshold = 5
    min_verbs = 1
    max_stop_words_ratio = 0.5
    min_tech_terms = 1

    length = len(doc)
    num_verbs = sum(1 for token in doc if token.pos_ == "VERB")
    num_stop_words = sum(1 for token in doc if token.is_stop)
    tech_term_count = sum(1 for token in doc if token.lemma_.lower() in TECH_TERMS)

    if length < length_threshold:
        return "poor"
    if num_verbs < min_verbs:
        return "poor"
    if (num_stop_words / length) > max_stop_words_ratio:
        return "average"
    if tech_term_count < min_tech_terms:
        return "average"

    return "good"

# Charger les variables d'environnement
load_dotenv()

# Récupérer le chemin du dépôt local à partir des variables d'environnement
LOCAL_REPO_PATH = os.getenv('LOCAL_REPO_PATH')


def main():
    print(f"Loading files from {LOCAL_REPO_PATH}")
    files = list_go_files(LOCAL_REPO_PATH)

    timestamp = datetime.now().strftime("%Y%m%d_%H")
    dataset_directory = "./dataset/file"
    dataset_path = f"{dataset_directory}/functions_dataset_{timestamp}.jsonl"

    # Créer le répertoire s'il n'existe pas
    os.makedirs(dataset_directory, exist_ok=True)

    with open(dataset_path, 'a') as file:
        if files:
            for file_path in files:
                print(f"Processing file: {file_path}")
                functions = parse_go_file(file_path)
                for function in functions:
                    comment = function.get('comment', '')
                    analyzed_comment = analyze_comment(comment)

                    quality_score = analyzed_comment.pop('quality_score', None)
                    if quality_score in {'good'}:
                        function['comment_analysis'] = analyzed_comment
                        function['file_name'] = os.path.splitext(os.path.basename(file_path))[0]
                        file.write(json.dumps(function) + '\n')

    print(f"Dataset saved to {dataset_path}")


if __name__ == "__main__":
    main()
