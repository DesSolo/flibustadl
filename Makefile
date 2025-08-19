PROJECT_NAME=flibustadl

LOCAL_BIN := $(CURDIR)/bin

MAIN_FILE_PATH=cmd/main.go

LDFLAGS = "-w -s"

build:
	go build -ldflags ${LDFLAGS} -o ${LOCAL_BIN}/${PROJECT_NAME} ${MAIN_FILE_PATH}

tag:
	scripts/tag.sh