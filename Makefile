DEST := /tmp
NAME := blarg
BINARY := $(PWD)/golem

test:
	go test -v ./...

build:
	go build

new: clean build init
	cd $(DEST)/$(NAME) && $(BINARY) add group releases --rest
	cd $(DEST)/$(NAME) && $(BINARY) add route releases additional -m POST
	cd $(DEST)/$(NAME) && $(BINARY) add group hello
	cd $(DEST)/$(NAME) && $(BINARY) add route hello world -p funky/world id count:int
	cd $(DEST)/$(NAME) && $(BINARY) add route hello new -m POST
	cd $(DEST)/$(NAME) && $(BINARY) add model hello world:string count:int
	cd $(DEST)/$(NAME) && $(BINARY) add model --struct metric time:time.Time key value job:*Job
	cd $(DEST)/$(NAME) && $(BINARY) add model --struct job time:time.Time name external_id:primitive.ObjectID
	cd $(DEST)/$(NAME) && $(BINARY) add event jobs event id job:*Job
	cd $(DEST)/$(NAME) && $(BINARY) add event reporter metric:*Metric -c '$(NAME).summary.report'
	cd $(DEST)/$(NAME) && $(BINARY) add event metrics time:time.Time key value -c 'metrics.report' --receiver -p Metric -t 'reporter'
	cd $(DEST)/$(NAME) && $(BINARY) add event flame -r time:time.Time download:float64 upload:float64
	cd $(DEST)/$(NAME) && $(BINARY) plugin enable cache
	cd $(DEST)/$(NAME) && $(BINARY) add worker ProcessRelease id
	cd $(DEST)/$(NAME) && $(BINARY) add queue downloads -c 3
	cd $(DEST)/$(NAME) && $(BINARY) add worker process_download id -s '0 0 11 * * *' -q downloads
	cd $(DEST)/$(NAME) && $(BINARY) generate
	cd $(DEST)/$(NAME) && $(BINARY) readme
	cd $(DEST)/$(NAME) && $(BINARY) routes

init: clean build
	@mkdir -p $(DEST)
	cd $(DEST) && $(BINARY) init $(NAME) github.com/test/$(NAME)

clean:
	@echo "Cleaning... $(DEST)/$(NAME)"
	@rm -rf $(DEST)/$(NAME)
	@rm -rf golem

install: build
	go install

.PHONY: test run
