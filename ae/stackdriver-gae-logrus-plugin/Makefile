ODIR=./bin
VERSION=`cat VERSION`
GCLOUD_VERSION=`cat VERSION | sed 's/\./-/g'`

.PHONY: build
build:
	go build -o bin/http-service-stackdriver-logging -ldflags "-X main.version=$(VERSION)" ./examples/http-service/main.go

.PHONY: run
run:
	go run -ldflags "-X main.version=$(VERSION) -X main.projectID=$(PROJECT_ID)" ./examples/http-service/main.go

.PHONY: deploy-gae
deploy-gae:
	gcloud app deploy --version=$(GCLOUD_VERSION) ./examples/http-service/app.yaml

.PHONY: clean
clean:
	-@rm -r $(ODIR)/* 2> /dev/null || true
