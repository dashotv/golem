
test:
	go test -v ./...

build: templates
	go build

gen: build
	cd test && ../golem generate

new: build
	cd .. && rm -rf blarg
	cd .. && ./golem/golem new blarg github.com/dashotv/blarg
	cd ../blarg && go mod init github.com/dashotv/blarg
	cd ../blarg && git init .
	cd ../blarg && ../golem/golem new model hello world:string count:int
	cd ../blarg && ../golem/golem new route /releases --rest
	cd ../blarg && ../golem/golem new route /hello id count:int
	cd ../blarg && ../golem/golem new route /hello/new
	cd ../blarg && ../golem/golem generate

#templates:
#	cd generators/templates/ && go-bindata -pkg templates -o bindata.go *.tgo

install: build
	go install

.PHONY: test run templates
