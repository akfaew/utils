TEST_ARGS = -failfast

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: updatelint
updatelint:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest


.PHONY: test
test: fmt lint
	go test $(TEST_ARGS) ./...

.PHONY: test-regen
test-regen: fmt
	mkdir -p testdata/output
	go test -regen $(TEST_ARGS) ./...

.PHONY: test-cover
test-cover: fmt
	go test $(TEST_ARGS) -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out

.PHONY: push
push: test
	git push
	git push --tags

.PHONY: update
update:
	go get -u
	go mod tidy
	go mod verify

.PHONY: clean
clean:
	rm -f coverage.out

.PHONY: setup
setup: updatelint update
