PROJECT_NAME=flibustadl

LOCAL_BIN := $(CURDIR)/bin

MAIN_FILE_PATH=cmd/main.go

LDFLAGS = "-w -s"

build:
	go build -ldflags ${LDFLAGS} -o ${LOCAL_BIN}/${PROJECT_NAME} ${MAIN_FILE_PATH}

tag:
	scripts/tag.sh

tidy:
	go mod tidy

.prep_bin:
	mkdir -p ${LOCAL_BIN}

GORELEASER_VERSION=v2.11.2

.install-goreleaser:
	curl -Ls https://github.com/goreleaser/goreleaser/releases/download/${GORELEASER_VERSION}/goreleaser_Linux_x86_64.tar.gz | tar xvz -C ${LOCAL_BIN} goreleaser

install-deps: \
	.prep_bin \
	.install-goreleaser

release: tag
	${LOCAL_BIN}/goreleaser release --clean