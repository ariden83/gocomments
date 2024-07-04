import requests
from pygments import lex
from pygments.lexers import GoLexer
from pygments.token import Token
from dotenv import load_dotenv
from pathlib import Path
from datetime import datetime
import shutil
import json
import os
import subprocess
import tempfile

def parse_go_file(file_path):
    # print(f"Parsing file: {file_path}")
    script_dir = os.path.dirname(os.path.abspath(__file__))
    result = subprocess.run(
        ['go', 'run', 'parse_functions.go', '--', file_path],
        cwd=script_dir,
        capture_output=True,
        text=True
    )

    # Vérifiez si le code de retour est différent de 0
    if result.returncode != 0:
        print(f"Failed to parse file {file_path}: {result.stderr}")
        return []

    # Affichez la sortie brute pour le débogage
    print(f"Raw stdout: '{result.stdout}'")

    try:
        parsed_output = json.loads(result.stdout)
        # Check if parsed_output is "None"
        if parsed_output is None:
            print(f"Output is None, returning default value []")
            return []

    except json.JSONDecodeError as e:
        print(f"JSON decode error: {e}")
        print(f"Failed to parse JSON from output: '{result.stdout}'")
        return []

    print('Parsed output:', parsed_output)
    return parsed_output

# Charge les variables d'environnement du fichier .env
load_dotenv()

# Récupère le token GitHub à partir des variables d'environnement
LOCAL_REPO_PATH = os.getenv('LOCAL_REPO_PATH')

def list_go_files(directory):
    try:
        directory_path = Path(directory)
        if not directory_path.exists():
            raise FileNotFoundError(f"Le répertoire {directory} n'existe pas.")
        if not directory_path.is_dir():
            raise NotADirectoryError(f"{directory} n'est pas un répertoire.")

        exclude_patterns = ['_test.go', '_example.go']

        def is_excluded(file):
            if any(file.name.endswith(pattern) for pattern in exclude_patterns):
                return True
            # Check if any part of the path is named 'vendor'
            return 'vendor' in file.parts

        go_files = [str(file) for file in directory_path.rglob('*.go') if not is_excluded(file)]
        return go_files

    except (FileNotFoundError, NotADirectoryError, PermissionError) as e:
        print(f"Erreur: {e}")
        return []


print("load files from ", LOCAL_REPO_PATH)
files = list_go_files(LOCAL_REPO_PATH)

timestamp = datetime.now().strftime("%Y%m%d_%H")
dataset_path = f"./file/functions_dataset_{timestamp}.jsonl"

# Ouvrir le fichier en mode append
with open(dataset_path, 'a') as file:
    if files:
        for file_path in files:
            print("file", file_path)
            functions = parse_go_file(file_path)
            for function in functions:
                # Écrire chaque fonction comme une ligne JSON séparée
                file.write(json.dumps(function) + '\n')

print(f"Dataset saved to {dataset_path}")
