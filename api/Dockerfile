FROM model_tensorflow-container

# Désactiver les opérations oneDNN si besoin
ENV TF_ENABLE_ONEDNN_OPTS=0

# Installer Go
ENV GO_VERSION=1.21.6
RUN apt-get update && \
    apt-get install -y wget tar && \
    wget https://golang.org/dl/go$GO_VERSION.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go$GO_VERSION.linux-amd64.tar.gz && \
    rm go$GO_VERSION.linux-amd64.tar.gz && \
    apt-get clean

# Définir le GOPATH et mettre à jour le PATH
ENV GOPATH=/go
ENV PATH=$PATH:/usr/local/go/bin:$GOPATH/bin

# Créer un répertoire pour l'application
RUN mkdir -p /app
RUN mkdir -p /app/models
WORKDIR /app

# Installer TensorFlow C API
RUN curl -L "https://storage.googleapis.com/tensorflow/libtensorflow/libtensorflow-cpu-linux-x86_64-2.6.0.tar.gz" | tar -C /usr/local -xz

# Définir le chemin de la bibliothèque TensorFlow
ENV LIBRARY_PATH=$LIBRARY_PATH:/usr/local/lib
ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/local/lib
ENV C_INCLUDE_PATH=$C_INCLUDE_PATH:/usr/local/include
ENV CPLUS_INCLUDE_PATH=$CPLUS_INCLUDE_PATH:/usr/local/include

# Exposer le port de votre application, si nécessaire (ajustez selon vos besoins)
EXPOSE 8080

CMD ["sh", "-c", "go build -o myapp . && ./myapp"]
