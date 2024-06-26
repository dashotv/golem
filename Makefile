DEST := /tmp/golem
NAME := blarg
BINARY := $(PWD)/golem

test:
	go test -v ./...

build:
	go build

new: clean build
	scripts/new.sh $(DEST) $(BINARY)

init: clean build
	@mkdir -p $(DEST)
	cd $(DEST) && $(BINARY) init $(NAME) github.com/test/$(NAME)

clean:
	@echo "Cleaning... $(DEST)/$(NAME)"
	@-find $(DEST)/$(NAME) -delete
	@rm -rf $(BINARY)

install: build
	go install

.PHONY: test run
