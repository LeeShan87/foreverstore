build:
	@go build -o ./bin/fs

run: build
	@./bin/fs

test:
	@go test ./... -v --count=1

clean:
	@rm -rf ./bin