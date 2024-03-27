DEST := /tmp
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
	@rm -rf $(DEST)/$(NAME)
	@rm -rf $(BINARY)

install: build
	go install

.PHONY: test run
