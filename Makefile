all: build

export GOPATH := $(GOPATH)
export GOBIN := $(CURDIR)/bin

CURRENT_GIT_GROUP := wangqingang
CURRENT_GIT_REPO := cunxun

proto:
	protoc -I pb pb/$(CURRENT_GIT_REPO).proto --go_out=plugins=grpc:pb
	protoc -I $(CURRENT_GIT_REPO) $(CURRENT_GIT_REPO)/log.proto --go_out=plugins=grpc:$(CURRENT_GIT_REPO)

build:
	go install $(CURRENT_GIT_GROUP)/$(CURRENT_GIT_REPO)

test:
	go test $(CURRENT_GIT_GROUP)/$(CURRENT_GIT_REPO)/checkcode
	go test $(CURRENT_GIT_GROUP)/$(CURRENT_GIT_REPO)/login
	go test $(CURRENT_GIT_GROUP)/$(CURRENT_GIT_REPO)/password
	go test $(CURRENT_GIT_GROUP)/$(CURRENT_GIT_REPO)/phone
	go test $(CURRENT_GIT_GROUP)/$(CURRENT_GIT_REPO)/token/token_lib
	go test $(CURRENT_GIT_GROUP)/$(CURRENT_GIT_REPO)/token
	go test $(CURRENT_GIT_GROUP)/$(CURRENT_GIT_REPO)/model
	go test $(CURRENT_GIT_GROUP)/$(CURRENT_GIT_REPO)/middleware
	go test $(CURRENT_GIT_GROUP)/$(CURRENT_GIT_REPO)/handler

clean:
	@rm -rf vendor bin _project

.PHONY: proto build test clean
