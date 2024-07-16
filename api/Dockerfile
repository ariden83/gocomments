FROM golang:1.21-bookworm

RUN curl -L https://storage.googleapis.com/tensorflow/versions/2.16.1/libtensorflow-cpu-linux-$(uname -m).tar.gz \
    | tar xz --directory /usr/local \
    && ldconfig

WORKDIR /app
