import requests
from pygments import lex
from pygments.lexers import GoLexer
from pygments.token import Token
from dotenv import load_dotenv
from pathlib import Path
import shutil
import json
import os
import subprocess
import tempfile

def search_go_repositories(query="language:go", sort="stars", per_page=10, page=1):
    url = f"https://api.github.com/search/repositories"
    params = {
        "q": query,
        "sort": sort,
        "per_page": per_page,
        "page": page
    }
    response = requests.get(url, params=params, headers=headers)
    if response.status_code == 200:
        return response.json()
    else:
        print(f"Failed to fetch repositories: {response.status_code}")
        return None

def list_files_in_repo(repo_owner, repo_name, path=""):
    url = f"https://api.github.com/repos/{repo_owner}/{repo_name}/contents/{path}"
    response = requests.get(url, headers=headers)
    if response.status_code == 200:
        return response.json()
    else:
        print(f"Failed to list files in repo: {response.status_code}")
        return None

def fetch_file_content(repo_owner, repo_name, file_path):
    url = f"https://api.github.com/repos/{repo_owner}/{repo_name}/contents/{file_path}"
    print(f"fetch File Content saved to {url}")

    response = requests.get(url, headers=headers)
    if response.status_code == 200:
        file_info = response.json()
        if file_info.get("encoding") == "base64":
            import base64
            content = base64.b64decode(file_info["content"]).decode("utf-8")
            return content
    print(f"Failed to fetch file content: {response.status_code}")
    return None

def extract_functions_and_comments(go_code):
    lexer = GoLexer()
    tokens = lex(go_code, lexer)

    functions = []
    current_function = None
    current_comment = None

    for token_type, token_value in tokens:
        if token_type in Token.Comment:
            current_comment = token_value.strip()
        elif token_type in Token.Keyword and token_value == "func":
            current_function = token_value
        elif token_type in Token.Name and current_function == "func":
            function_name = token_value
            if current_comment is not None and not should_ignore_comment(current_comment, ignored_annotations):
                functions.append({
                    "name": function_name,
                    "comment": current_comment
                })
            current_function = None
            current_comment = None

    return functions

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

def should_ignore_comment(comment, ignored_annotations):
    # Vérifie si l'un des mots-clés est présent dans le commentaire
    for annotation in ignored_annotations:
        if annotation in comment:
            return True
    return False

def fetch_and_parse(repo_owner, repo_name, file_path):
    content = fetch_file_content(repo_owner, repo_name, file_path)
    if content:
        # Créer un fichier temporaire
        with tempfile.NamedTemporaryFile(delete=False, suffix=".go") as temp_file:
            temp_file.write(content.encode('utf-8'))
            temp_file_path = temp_file.name

        print(f"Temp file created at: {temp_file_path}")

        # Créer un répertoire de travail temporaire
        with tempfile.TemporaryDirectory() as work_dir:
            temp_filename = os.path.basename(temp_file_path)
            temp_work_path = os.path.join(work_dir, temp_filename)

            # Copier le fichier dans le répertoire de travail
            shutil.copy(temp_file_path, temp_work_path)
            print(f"File copied to work directory: {temp_work_path}")

            try:
                # Analyser le fichier temporaire dans le répertoire de travail
                functions_and_comments = parse_go_file(temp_work_path)
                return functions_and_comments
            finally:
                # Supprimer le fichier temporaire
                os.remove(temp_file_path)
                print(f"Temp file deleted: {temp_file_path}")
    return None

# Example usage:
# go_code = fetch_file_content("repository_owner", "repository_name", "path_to_file.go")
# functions_with_comments = extract_functions_and_comments(go_code)

# Charge les variables d'environnement du fichier .env
load_dotenv()

# Liste des annotations à ignorer
ignored_annotations = ["FIXME", "NOTE", "go:embed", "TODO", "BUG", "deprecated"]

# Récupère le token GitHub à partir des variables d'environnement
GITHUB_TOKEN = os.getenv('GITHUB_TOKEN')
LOCAL_REPO_PATH = os.getenv('LOCAL_REPO_PATH')

headers = {
    "Authorization": f"token {GITHUB_TOKEN}",
    "Accept": "application/vnd.github.v3+json"
}

repos = search_go_repositories()
dataset = []
functions_with_comments = []

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


for repo in repos.get('items', []):
   owner = repo['owner']['login']
   name = repo['name']

   files = list_files_in_repo(owner, name)
   if files:
      for file in files:
         print("file", file)
         if file['name'].endswith('.go'):
            functions = fetch_and_parse(owner, name, file['path'])
            dataset.extend(functions)

# Example dataset structure:
# [{'name': 'MyFunction', 'comment': '// This is a function comment'}, ...]

dataset_path = "functions_dataset_github.json"
with open(dataset_path, 'w') as file:
    json.dump(dataset, file, indent=4)

print(f"Dataset saved to {dataset_path}")
