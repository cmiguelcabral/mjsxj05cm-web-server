all: cross

cross: prepare
	test -d dist || mkdir dist
	test -d dist/static || mkdir dist/static
	GOOS=linux GOARCH=arm go build -v -o dist/web-config-server .

test: prepare
	test -d test || mkdir test
	go build -o test/web_control

clean:
	rm -rf dist test

prepare:
	go get -v -t -d ./...