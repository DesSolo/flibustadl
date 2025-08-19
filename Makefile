PROJECT_NAME=flibustadl

LOCAL_BIN := $(CURDIR)/bin

MAIN_FILE_PATH=cmd/main.go

build:
	go build -o ${LOCAL_BIN}/${PROJECT_NAME} ${MAIN_FILE_PATH}