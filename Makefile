
test:
	go test -v ./...

build: templates
	go build

gen: build
	cd test && ../golem generate

new: build
	rm -rf blarg
	./golem new blarg github.com/dashotv/blarg
	cd blarg && ../golem new model hello world:string count:int

templates:
	go-bindata -pkg templates -o generators/templates/bindata.go generators/templates/*.tgo

.PHONY: test run templates
