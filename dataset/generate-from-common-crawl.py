import re
import json

def find_golang_code(content):
    golang_pattern = re.compile(r'```go(.*?)```', re.DOTALL)
    return golang_pattern.findall(content)

def collect_golang_code(file_path):
    code_snippets = []
    with warcio.ArchiveIterator(open(file_path, 'rb')) as archive:
        for record in archive:
            if record.rec_type == 'response':
                content = record.content_stream().read().decode('utf-8')
                snippets = find_golang_code(content)
                if snippets:
                    code_snippets.extend(snippets)
    return code_snippets

# Example usage:
golang_code_snippets = collect_golang_code('./warc/file')
for snippet in golang_code_snippets:
    print(snippet)

from pyspark.sql import SparkSession

# Initialize Spark session
spark = SparkSession.builder \
    .appName("CommonCrawlProcessing") \
    .getOrCreate()

# Load WARC files
warc_files = spark.read.format("com.databricks.spark.avro").load("/path/to/warc/files")

# Extract Golang code snippets
def extract_code(content):
    golang_pattern = re.compile(r'```go(.*?)```', re.DOTALL)
    return golang_pattern.findall(content)

# Apply extraction function
warc_files_rdd = warc_files.rdd.map(lambda x: extract_code(x['content']))


with open('golang_code_snippets.json', 'w') as f:
    json.dump(golang_code_snippets, f, indent=4)
