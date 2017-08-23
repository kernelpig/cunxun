all: build

export GOPATH := $(CURDIR)/_project
export GOBIN := $(CURDIR)/bin

CURRENT_GIT_GROUP := wangqingang
CURRENT_GIT_REPO := cunxun

folder_dep:
	mkdir -p $(CURDIR)/_project/src/$(CURRENT_GIT_GROUP)
	test -d $(CURDIR)/_project/src/$(CURRENT_GIT_GROUP)/$(CURRENT_GIT_REPO) || ln -s $(CURDIR) $(CURDIR)/_project/src/$(CURRENT_GIT_GROUP)

glide:
	mkdir -p $(CURDIR)/vendor
	glide install

proto:
	protoc -I pb pb/$(CURRENT_GIT_REPO).proto --go_out=plugins=grpc:pb
	protoc -I $(CURRENT_GIT_REPO) $(CURRENT_GIT_REPO)/log.proto --go_out=plugins=grpc:$(CURRENT_GIT_REPO)

build: folder_dep
	go install $(CURRENT_GIT_GROUP)/$(CURRENT_GIT_REPO)

test:
	go test $(CURRENT_GIT_GROUP)/$(CURRENT_GIT_REPO)/model
	go test $(CURRENT_GIT_GROUP)/$(CURRENT_GIT_REPO)/test
	go test $(CURRENT_GIT_GROUP)/$(CURRENT_GIT_REPO)/utils/password

test-v:
	go test -v $(CURRENT_GIT_GROUP)/$(CURRENT_GIT_REPO)/test

clean:
	@rm -rf vendor bin _project

.PHONY: deps install test add_dep clean
