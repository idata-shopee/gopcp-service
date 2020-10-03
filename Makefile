clean:
	@go mod tidy

update:
	@go get -u

test:
	@go test -v -race

cover:
	@go test -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html

.PHONY:	test
