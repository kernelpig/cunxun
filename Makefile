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

build:
	go install $(CURRENT_GIT_GROUP)/$(CURRENT_GIT_REPO)
	go install $(CURRENT_GIT_GROUP)/$(CURRENT_GIT_REPO)/xinit

test:
	go test $(CURRENT_GIT_GROUP)/$(CURRENT_GIT_REPO)/error
	go test $(CURRENT_GIT_GROUP)/$(CURRENT_GIT_REPO)/checkcode
	go test $(CURRENT_GIT_GROUP)/$(CURRENT_GIT_REPO)/login
	go test $(CURRENT_GIT_GROUP)/$(CURRENT_GIT_REPO)/password
	go test $(CURRENT_GIT_GROUP)/$(CURRENT_GIT_REPO)/phone
	go test $(CURRENT_GIT_GROUP)/$(CURRENT_GIT_REPO)/token/token_lib
	go test $(CURRENT_GIT_GROUP)/$(CURRENT_GIT_REPO)/token
	go test $(CURRENT_GIT_GROUP)/$(CURRENT_GIT_REPO)/model
	go test $(CURRENT_GIT_GROUP)/$(CURRENT_GIT_REPO)/middleware
	go test $(CURRENT_GIT_GROUP)/$(CURRENT_GIT_REPO)/handler
	go test $(CURRENT_GIT_GROUP)/$(CURRENT_GIT_REPO)/initial

clean:
	@rm -rf vendor bin _project

.PHONY: folder_dep glide test build clean
