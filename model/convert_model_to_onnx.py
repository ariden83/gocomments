import os
import glob
import shutil
import subprocess
import time
import sys
from datetime import datetime

# Chemin source et destination.
source_path = 'runs'
dest_base_path = os.getenv('ONNX_DEST_PATH', 'onnx')

# Liste des dossiers dans le répertoire spécifié.
directories = [d for d in os.listdir(source_path) if os.path.isdir(os.path.join(source_path, d))]

# Récupération des informations sur les dossiers avec leur date de création.
dir_info = [(d, os.path.getctime(os.path.join(source_path, d))) for d in directories]

# Tri des dossiers par date de création (le plus récent en premier).
dir_info.sort(key=lambda x: x[1], reverse=True)

# Vérification et affichage
if dir_info:
    newest_directory, creation_time = dir_info[0]
    if newest_directory.startswith('saved_model_'):
        print(f"Le dernier dossier créé est : {newest_directory}")
        print(f"Date de création : {time.ctime(creation_time)}")
        latest_dir = newest_directory
    else:
        print(f"Aucun dossier ne commence par 'saved_model_' parmi les dossiers trouvés.")
else:
    print("Aucun dossier trouvé dans le répertoire spécifié.")
    sys.exit(1)  # Quitte le script avec un code de sortie non nul pour signaler l'erreur.

latest_dir_name = os.path.basename(os.path.normpath(latest_dir))

# Créer un sous-dossier dans le chemin de destination.
dest_path = os.path.join(dest_base_path, latest_dir_name)
os.makedirs(dest_path, exist_ok=True)

# Chemin du modèle TensorFlow et du fichier ONNX.
model_path = os.path.join(source_path, latest_dir)
output_onnx_path = os.path.join(dest_path, 'model.onnx')

try:
    convert_command = [
        'python', '-m', 'tf2onnx.convert',
        '--saved-model', model_path,
        '--output', output_onnx_path
    ]
    subprocess.run(convert_command, check=True)
    print(f"Conversion réussie. Modèle ONNX enregistré à : {output_onnx_path}")
except subprocess.CalledProcessError as e:
    print(f"Erreur lors de la conversion du modèle : {e}")
    sys.exit(1)
