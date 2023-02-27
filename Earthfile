VERSION 0.6

ARG ALPINE_VERSION=3.17
ARG GO_VERSION=1.19
ARG LINTER_VERSION=v1.51.2
FROM golang:$GO_VERSION-alpine$ALPINE_VERSION
WORKDIR /app

stage:
  COPY --dir go.mod go.sum ./
  RUN go mod download -x
  COPY --dir protoc cmd .
  SAVE ARTIFACT /app

vendor:
    FROM +stage
    RUN go mod vendor

lint:
    # Installs golangci-lint to ./bin
    RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s $LINTER_VERSION
    COPY +stage/app .
    RUN ./bin/golangci-lint run --skip-dirs=vendor --skip-dirs=./gen/ --deadline=5m --tests=true -E revive \
      -E gosec -E unconvert -E goconst -E gocyclo -E goimports

build:
    FROM +vendor
    # compile app binary, save as artifact
    ARG VERSION="dev"
    ARG GOOS
    ARG GOARCH
    RUN go build -ldflags="-s -w -X 'main.version=${VERSION}'" -mod=vendor -o bin/protoc-gen-fieldmask ./cmd/protoc-gen-fieldmask/...
    SAVE ARTIFACT ./bin/protoc-gen-fieldmask /protoc-gen-fieldmask

zip:
    RUN apk add zip
    WORKDIR /artifacts
    ARG VERSION
    ARG ZIP_FILE_NAME
    ARG EXT
    COPY (+build/protoc-gen-fieldmask) protoc-gen-fieldmask${EXT}
    RUN zip -m protoc-gen-fieldmask-${VERSION}-${ZIP_FILE_NAME}.zip protoc-gen-fieldmask${EXT}
    SAVE ARTIFACT /artifacts

build-all:
    WORKDIR /artifacts
    COPY (+zip/artifacts/*.zip --GOOS=darwin --GOARCH=amd64 --ZIP_FILE_NAME=osx-x86_64) .
    COPY (+zip/artifacts/*.zip --GOOS=linux --GOARCH=386 --ZIP_FILE_NAME=linux-x86_32) .
    COPY (+zip/artifacts/*.zip --GOOS=linux --GOARCH=amd64 --ZIP_FILE_NAME=linux-x86_64) .
    COPY (+zip/artifacts/*.zip --GOOS=windows --GOARCH=386 --ZIP_FILE_NAME=win32 --EXT=.exe) .
    COPY (+zip/artifacts/*.zip --GOOS=windows --GOARCH=amd64 --ZIP_FILE_NAME=win64 --EXT=.exe) .
    SAVE ARTIFACT /artifacts AS LOCAL bin

test-gen:
    ARG DOCKER_PROTOC_VERSION=1.51_1
    FROM namely/protoc-all:$DOCKER_PROTOC_VERSION
    RUN mkdir /plugins
    COPY +build/protoc-gen-fieldmask /usr/local/bin/.
    COPY --dir protos .
    RUN entrypoint.sh -i protos -d protos/cases -l go -o gen
    RUN entrypoint.sh -i protos -d protos/cases/thirdpartyimport -l go -o gen
    RUN protoc -I/opt/include -Iprotos --fieldmask_out=gen protos/cases/*.proto protos/cases/thirdpartyimport/*.proto
    SAVE ARTIFACT gen /gen AS LOCAL test/gen

test:
    FROM +vendor
    RUN apk add build-base
    COPY --dir test .
    COPY --dir +test-gen/gen test/.
    RUN go mod vendor
    RUN go test github.com/idodod/protoc-gen-fieldmask/test

all:
    BUILD +lint
    BUILD +test
    BUILD +build-all

clean:
    LOCALLY
    RUN rm -rf test/gen vendor

