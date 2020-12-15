all: cross

cross: prepare
	test -d dist || mkdir dist
	test -d dist/static || mkdir dist/static
	export GIN_MODE=release
	GOOS=linux GOARCH=arm go build -v -o dist/web-config-server .

test: prepare
	test -d test || mkdir test
	export GIN_MODE=debug
	go build -o test/web-config-server .

clean:
	rm -rf dist test

prepare:
	go get -v -t -d ./...