
test:
	go test -v ./...

build: templates
	go build

gen: build
	cd test && ../golem generate

new: build
	./golem new blarg github.com/dashotv/blarg

templates:
	go-bindata -pkg templates -o generators/templates/bindata.go generators/templates/*.tgo

.PHONY: test run templates
