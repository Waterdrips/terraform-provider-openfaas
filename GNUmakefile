SWEEP?=us-east-1,us-west-2
TEST?=./...
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
PKG_NAME=openfaas

default: build

build: fmtcheck
	go install

travisbuild: default

localbuild: fmt default
	mkdir -p ~/.terraform.d/plugins/terraform.openfaas.com/openfaas/openfaas/1.0.0/linux_amd64/
	cp $(GOPATH)/bin/terraform-provider-openfaas ~/.terraform.d/plugins/terraform.openfaas.com/openfaas/openfaas/1.0.0/linux_amd64/terraform-provider-openfaas

sweep:
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	go test $(TEST) -v -sweep=$(SWEEP) $(SWEEPARGS)

test: fmtcheck
	go test $(TEST) -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

goimports:
	goimports -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

vendor-status:
	@govendor status

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

.PHONY: build sweep test testacc vet fmt fmtcheck errcheck vendor-status test-compile goimports travisbuild release release-deps

