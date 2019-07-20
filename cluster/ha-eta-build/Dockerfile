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
ENV GO_LDFLAGS="-s -w -X main.buildDate=${GO_BUILD_DATE} -X \"main.version=0.0.1\" -extldflags -static"
ENV GOOS="linux"
ENV GOARCH="amd64"

COPY ./src/go.mod ./
RUN go mod download
COPY ./src/ ./
RUN rm -rvf ./src/go.sum

RUN GO_BUILD_DATE=$(date -u +%Y-%m-%dT%H:%M:%S) GOOS=${GOOS} GOARCH=${GOARCH} \
  go build \
    -v \
    -ldflags "${GO_LDFLAGS}" -o ./ha-eta ./common.go ./main.go
