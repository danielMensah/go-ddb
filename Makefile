MODULE_PATH := github.com/danielMensah/go-ddb

mock:
	mockery

fmt:
	goimports -w . && gofmt -w .

lint:
	golangci-lint run --verbose --config golangci.yaml

vuln:
	govulncheck -show verbose ./...

sweep: mock fmt lint vuln

test:
	go test -cover -race ./...;

all: sweep test

pub:
	@if [ -z "$(version)" ]; then \
		echo "Error: version is required. Usage: make pub version=vX.Y.Z"; \
		exit 1; \
	fi; \
	if ! echo "$(version)" | grep -Eq '^v[0-9]+\.[0-9]+\.[0-9]+$$'; then \
		echo "Error: version must follow the format vX.Y.Z (e.g., v1.0.0)"; \
		exit 1; \
	fi
	GOPROXY=proxy.golang.org go list -m $(MODULE_PATH)@$(version)
