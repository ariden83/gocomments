import requests
from pygments import lex
from pygments.lexers import GoLexer
from pygments.token import Token
from dotenv import load_dotenv
import json
import os
import subprocess

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
    print(f"Parsing file: {file_path}")
    result = subprocess.run(['go', 'run', 'parse_functions.go', file_path], capture_output=True, text=True)
    if result.returncode == 0:
        return json.loads(result.stdout)
    else:
        print(f"Failed to parse file {file_path}: {result.stderr}")
        return []

def should_ignore_comment(comment, ignored_annotations):
    # Vérifie si l'un des mots-clés est présent dans le commentaire
    for annotation in ignored_annotations:
        if annotation in comment:
            return True
    return False

# Example usage:
# go_code = fetch_file_content("repository_owner", "repository_name", "path_to_file.go")
# functions_with_comments = extract_functions_and_comments(go_code)

# Charge les variables d'environnement du fichier .env
load_dotenv()

# Liste des annotations à ignorer
ignored_annotations = ["FIXME", "NOTE", "go:embed", "TODO", "BUG"]

# Récupère le token GitHub à partir des variables d'environnement
GITHUB_TOKEN = os.getenv('GITHUB_TOKEN')

headers = {
    "Authorization": f"token {GITHUB_TOKEN}",
    "Accept": "application/vnd.github.v3+json"
}

repos = search_go_repositories()
dataset = []
functions_with_comments = []

for repo in repos.get('items', []):
    owner = repo['owner']['login']
    name = repo['name']

    files = list_files_in_repo(owner, name)
    if files:
        for file in files:
            if file['name'].endswith('.go'):
                go_code = fetch_file_content(owner, name, file['path'])
                if go_code:
                    functions = extract_functions_and_comments(go_code)
                    dataset.extend(functions)

# Example dataset structure:
# [{'name': 'MyFunction', 'comment': '// This is a function comment'}, ...]

dataset_path = "functions_dataset.json"
with open(dataset_path, 'w') as file:
    json.dump(dataset, file, indent=4)

print(f"Dataset saved to {dataset_path}")
