default: build

build:
	go build -o terraform-provider-azuresim .

install: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/ndonathan/azuresim/0.1.0/$$(go env GOOS)_$$(go env GOARCH)
	cp terraform-provider-azuresim ~/.terraform.d/plugins/registry.terraform.io/ndonathan/azuresim/0.1.0/$$(go env GOOS)_$$(go env GOARCH)/terraform-provider-azuresim_v0.1.0

test:
	go test ./... -timeout 120s

testv:
	go test ./... -v -timeout 120s

testacc:
	TF_ACC=1 go test ./... -v -timeout 120s

testcover:
	go test ./... -coverprofile=coverage.out -timeout 120s
	go tool cover -func=coverage.out | tail -1

fmt:
	go fmt ./...

vet:
	go vet ./...

.PHONY: build install test testv testacc testcover fmt vet
