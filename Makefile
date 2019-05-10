TEST_ARGS = -failfast

fmt:
	go fmt ./...

test: fmt
	go test $(TEST_ARGS) ./...

test-regen: fmt
	mkdir -p testdata/output
	go test -regen $(TEST_ARGS) ./...

test-cover: fmt
	go test $(TEST_ARGS) -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out

push: test
	git push
	git push --tags

update:
	go get -u
	go mod tidy
	go mod verify

clean:
	rm -f coverage.out
