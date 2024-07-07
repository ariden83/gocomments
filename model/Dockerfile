FROM nvcr.io/nvidia/tensorflow:24.01-tf2-py3

# Désactiver les opérations oneDNN si besoin
ENV TF_ENABLE_ONEDNN_OPTS=0

# Installer les dépendances supplémentaires avec des versions spécifiques pour la compatibilité
RUN pip install --upgrade pip
RUN pip install fsspec==2023.9.2 transformers==4.28.1 tf-keras==2.14.1 keras==2.12.0 tokenizers==0.13.3 torch==1.13.1 python-dotenv requests pygments idna importlib-metadata libclang Markdown numpy tensorflow-estimator tensorflow-io-gcs-filesystem
RUN pip install -U datasets==2.10.1 scikit-learn==1.2.2 scipy==1.10.0 matplotlib==3.7.1 yahoo_fin==0.8.9.1

# Vérification de CUDA
RUN nvcc --version

# Vérification des GPU disponibles avec TensorFlow
RUN python -c "import tensorflow as tf; print('Num GPUs Available: ', len(tf.config.experimental.list_physical_devices('GPU')))"

# Vérification de la disponibilité de CUDA avec PyTorch
RUN python -c "import torch; print('Is CUDA available: ', torch.cuda.is_available()); print('Number of GPUs: ', torch.cuda.device_count())"