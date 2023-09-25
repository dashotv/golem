
test:
	go test -v ./...

build: templates
	go build

gen: build
	cd test && ../golem generate

new: clean build
	cd .. && rm -rf blarg
	cd .. && ./golem/golem new blarg github.com/dashotv/blarg
	cd ../blarg && go mod init github.com/dashotv/blarg
	cd ../blarg && git init .
	cd ../blarg && ../golem/golem new model hello world:string count:int
	cd ../blarg && ../golem/golem new route /releases --rest
	cd ../blarg && ../golem/golem new route /hello id count:int
	cd ../blarg && ../golem/golem new route /hello/new
	cd ../blarg && ../golem/golem new model --struct foo bar:int baz:string
	cd ../blarg && ../golem/golem new model --struct download_file id:primitive.ObjectID medium_id:primitive.ObjectID medium:*Medium num:int
	cd ../blarg && ../golem/golem generate

clean:
	rm -rf ../blarg
	rm -rf golem

install: build
	go install

.PHONY: test run templates
