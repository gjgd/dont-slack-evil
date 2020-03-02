.PHONY: build clean deploy

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/hello hello/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/dynamodb-example dynamodb-example/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/messages messages/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/incoming-webhook incoming-webhook/main.go
	cp apphome/sample.json bin/sample.json

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
