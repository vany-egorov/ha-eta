FROM debian:10 as build

ENV GO_VERSION="1.12.7"

RUN apt-get update
RUN apt-get install -y \
  curl \
  git-core

# <golang download>
WORKDIR /opt
RUN curl -O "https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz"
RUN mv ./go*tar.gz ./go.tar.gz
RUN tar -xf ./go.tar.gz
RUN rm -rf ./go.tar.gz
# </golang download>

ENV GOROOT="/opt/go"
ENV PATH="${GOROOT}/bin:${PATH}"

WORKDIR /usr/local/src/ha-eta

COPY ./src/go.mod ./
RUN go mod download
COPY ./src/ ./
RUN rm -rvf ./src/go.sum

RUN /bin/bash ./build.sh

ENV CGO_ENABLED=0
RUN go test ./apps/node -v

RUN ./ha-eta --version
RUN ./ha-eta --help

RUN ./ha-eta node --help
