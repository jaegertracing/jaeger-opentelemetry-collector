ALL_SRC := $(shell find . -name '*.go' \
                                -not -path '*/third_party/*' \
                                -type f | sort)

# ALL_PKGS is used with 'go cover'
ALL_PKGS := $(shell $(GOCMD) list $(sort $(dir $(ALL_SRC))) 2>/dev/null)

UNIT_TEST_PACKAGES := $(shell go list ./...)
GO_FLAGS ?= GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GO111MODULE=on
FMT_LOG=fmt.log
LINT_LOG=lint.log
GOPATH ?= "$(HOME)/go"
COLLECTOR_NAME ?= jaeger-opentelemetry-collector
DOCKER_NAMESPACE ?= quay.io/jaegertracing
DOCKER_TAG ?= latest

.DEFAULT_GOAL := test

.PHONY: test
test:
	go test $(ALL_PKGS) -cover -coverprofile=cover.out

.PHONY: ci
ci: check test

.PHONY: format
format:
	@${GOPATH}/bin/goimports -local "github.com/jaegertracing/jaeger-opentelemetry-collector" -l -w $(shell git ls-files "*\.go" | grep -v vendor) > ${FMT_LOG}
	@[ ! -s "$(FMT_LOG)" ] || (echo "Formatting:" | cat - $(FMT_LOG))

.PHONY: lint
lint:
	@${GOPATH}/bin/golint ./... > ${LINT_LOG}
	@[ ! -s "$(LINT_LOG)" ] || (echo "Lint issues found in:" | cat - $(LINT_LOG) && false)

.PHONY: check
check: format lint

.PHONY: install-tools
install-tools:
	${GO_FLAGS} go install \
		golang.org/x/lint/golint \
		golang.org/x/tools/cmd/goimports
