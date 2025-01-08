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

publish:
	GOPROXY=proxy.golang.org go list -m github.com/danielMensah/go-ddb@v0.0.2