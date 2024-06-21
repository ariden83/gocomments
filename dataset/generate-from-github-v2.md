Pour créer un dataset à partir des repositories GitHub contenant du code Go, avec des commentaires associés aux fonctions, et le tout sans télécharger les fichiers en local, vous pouvez suivre les étapes suivantes. L'objectif est de construire un dataset à utiliser avec TensorFlow pour une IA qui génère automatiquement des commentaires pour les fonctions Go. Voici un guide étape par étape pour accomplir cela en Python :

# Étapes Générales

1. **Utiliser l'API GitHub pour rechercher des repositories Go.**
2. **Récupérer les fichiers source contenant des fonctions Go.**
3. **Extraire les fonctions et leurs commentaires à partir du code.**
4. **Générer un dataset avec les fonctions et les commentaires associés.**
5. **Sauvegarder le dataset pour utilisation avec TensorFlow.**

## Pré-requis
Accès à l'API GitHub : Vous aurez besoin d'un token d'accès GitHub pour faire des requêtes à l'API. Créez un token via GitHub Developer Settings.

**Bibliothèques Python** : Installez les bibliothèques nécessaires.

```
pip install requests pygments
```

## Étape 1: Utiliser l'API GitHub pour Rechercher des Repositories Go

```
import requests

def search_go_repositories(query="language:go", sort="stars", per_page=10, page=1):
url = f"https://api.github.com/search/repositories"
params = {
"q": query,
"sort": sort,
"per_page": per_page,
"page": page
}
headers = {
"Accept": "application/vnd.github.v3+json",
"Authorization": f"token YOUR_GITHUB_TOKEN"  # Remplacez par votre token
}
response = requests.get(url, params=params, headers=headers)
if response.status_code == 200:
return response.json()
else:
print(f"Failed to fetch repositories: {response.status_code}")
return None

repos = search_go_repositories()
```

## Étape 2: Récupérer les Fichiers Source Contenant des Fonctions Go

Vous pouvez lister les fichiers d'un repository spécifique. Ensuite, vous pouvez lire leur contenu via l'API sans les télécharger localement.

```
def list_files_in_repo(repo_owner, repo_name, path=""):
url = f"https://api.github.com/repos/{repo_owner}/{repo_name}/contents/{path}"
headers = {
"Accept": "application/vnd.github.v3+json",
"Authorization": f"token YOUR_GITHUB_TOKEN"
}
response = requests.get(url, headers=headers)
if response.status_code == 200:
return response.json()
else:
print(f"Failed to list files in repo: {response.status_code}")
return None

def fetch_file_content(repo_owner, repo_name, file_path):
url = f"https://api.github.com/repos/{repo_owner}/{repo_name}/contents/{file_path}"
headers = {
"Accept": "application/vnd.github.v3+json",
"Authorization": f"token YOUR_GITHUB_TOKEN"
}
response = requests.get(url, headers=headers)
if response.status_code == 200:
file_info = response.json()
if file_info.get("encoding") == "base64":
import base64
content = base64.b64decode(file_info["content"]).decode("utf-8")
return content
print(f"Failed to fetch file content: {response.status_code}")
return None
```

## Étape 3: Extraire les Fonctions et leurs Commentaires

Utilisez une bibliothèque comme pygments pour analyser le code et extraire les fonctions Go et leurs commentaires.

```
from pygments import lex
from pygments.lexers import GoLexer
from pygments.token import Token

def extract_functions_and_comments(go_code):
lexer = GoLexer()
tokens = lex(go_code, lexer)

    functions = []
    comments = []
    current_function = None
    current_comment = None
    
    for token_type, token_value in tokens:
        if token_type in Token.Comment:
            current_comment = token_value.strip()
        elif token_type in Token.Keyword and token_value == "func":
            current_function = token_value
        elif token_type in Token.Name and current_function == "func":
            function_name = token_value
            functions.append({
                "name": function_name,
                "comment": current_comment
            })
            current_function = None
            current_comment = None
    
    return functions

# Example usage:
# go_code = fetch_file_content("repository_owner", "repository_name", "path_to_file.go")
# functions_with_comments = extract_functions_and_comments(go_code)
```

## Étape 4: Générer un Dataset
Créez une liste des fonctions et de leurs commentaires pour constituer le dataset.

```
dataset = []

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
```

## Étape 5: Sauvegarder le Dataset
Sauvegardez le dataset dans un fichier CSV ou JSON pour une utilisation future avec TensorFlow.

En CSV
```
import csv

with open('functions_dataset.csv', mode='w') as file:
writer = csv.DictWriter(file, fieldnames=['name', 'comment'])
writer.writeheader()
for data in dataset:
writer.writerow(data)
```

En JSON
```
import json

with open('functions_dataset.json', 'w') as file:
json.dump(dataset, file)
```

## Conclusion

En suivant ces étapes, vous pouvez collecter et structurer un dataset de fonctions Go et de leurs commentaires directement à partir de GitHub. Ce dataset peut ensuite être utilisé pour entraîner un modèle TensorFlow pour générer des commentaires de fonctions. Assurez-vous de respecter les conditions d'utilisation de l'API GitHub et de gérer le quota de requêtes pour éviter tout blocage.


##  Comment Configurer les Autorisations pour un Token GitHub

Go [here](https://github.com/settings/tokens) to generate a new token.

Pour parser les différents repositories publics sur GitHub, vous n’avez pas besoin de nombreuses autorisations puisque ces repositories sont accessibles publiquement. Cependant, il est crucial de configurer les autorisations correctement pour s’assurer que vous avez accès aux informations nécessaires tout en maintenant la sécurité. Voici une liste des autorisations minimales nécessaires pour effectuer cette tâche efficacement :

1. **Allez dans les Paramètres GitHub** : Connectez-vous à GitHub, cliquez sur votre avatar en haut à droite et sélectionnez "Settings".

2. **Allez dans Developer Settings** : Dans le menu de gauche, allez en bas et cliquez sur "Developer settings".

3. **Allez dans Personal Access Tokens** : Cliquez sur "Personal access tokens" puis "Tokens (classic)".

4. **Cliquez sur Generate New Token** : Donnez un nom à votre token, et sélectionnez une expiration si nécessaire.

5. **Sélectionnez les Autorisations** : Cochez uniquement les autorisations nécessaires :

**public_repo** : Pour accéder aux repositories publics.

6. **Générez et Copiez le Token** : Cliquez sur "Generate token", puis copiez le token généré et sauvegardez-le dans un endroit sûr.

