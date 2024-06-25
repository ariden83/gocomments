import requests
from pyspark.sql import SparkSession
from pyspark.sql.types import StructType, StructField, StringType, IntegerType
from pyspark.sql.functions import udf


# Configuration GitHub
GITHUB_TOKEN = 'votre_token_ici'
headers = {'Authorization': f'token {GITHUB_TOKEN}'}

def search_golang_repos(query, max_results=10):
    url = 'https://api.github.com/search/repositories'
    params = {'q': f'{query}+language:go', 'per_page': max_results}
    response = requests.get(url, headers=headers, params=params)
    return response.json()

def get_files_from_repo(owner, repo):
    url = f'https://api.github.com/repos/{owner}/{repo}/git/trees/main?recursive=1'
    response = requests.get(url, headers=headers)
    tree = response.json().get('tree', [])
    go_files = [file['path'] for file in tree if file['path'].endswith('.go')]
    return go_files

def get_file_content(owner, repo, path):
    url = f'https://raw.githubusercontent.com/{owner}/{repo}/main/{path}'
    response = requests.get(url)
    return response.text

# Lire les fichiers en mémoire
def read_go_files_from_github(query='golang', max_results=5):
    repos = search_golang_repos(query, max_results)
    files_content = []

    for repo in repos['items']:
        owner = repo['owner']['login']
        repo_name = repo['name']
        files = get_files_from_repo(owner, repo_name)

        for file in files:
            content = get_file_content(owner, repo_name, file)
            files_content.append((repo_name, file, content))

    return files_content

# Initialiser Spark
spark = SparkSession.builder.appName("GitHub Go Files Live Analysis").getOrCreate()

# Lire les fichiers Go
files_data = read_go_files_from_github()

# Définir le schéma du DataFrame
schema = StructType([
    StructField('repo', StringType(), True),
    StructField('path', StringType(), True),
    StructField('content', StringType(), True)
])

# Créer un DataFrame Spark
df = spark.createDataFrame(files_data, schema=schema)

def count_lines(content):
    return len(content.split('\n'))

count_lines_udf = udf(count_lines, IntegerType())
df = df.withColumn('line_count', count_lines_udf(df['content']))

# Afficher les résultats
df.show()

# Spécifiez le chemin où vous voulez sauvegarder le fichier JSON
output_path = '/warc/file/github_go_files.json'

# Écrire le DataFrame en JSON
df.write.mode('overwrite').json(output_path)

# Vérifiez que le fichier a bien été sauvegardé.
print(f'Les résultats ont été sauvegardés en JSON à : {output_path}')

