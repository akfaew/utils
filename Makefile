TEST_ARGS = -failfast

fmt:
	go fmt ./...

test: fmt
	go test $(TEST_ARGS) ./...

test-cover: fmt
	go test $(TEST_ARGS) -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out

update:
	go get -u
	go mod tidy
	go mod verify
