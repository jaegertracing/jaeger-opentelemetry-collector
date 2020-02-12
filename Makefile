UNIT_TEST_PACKAGES := $(shell go list ./pkg/...)
GO_FLAGS ?= GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GO111MODULE=on
GOLINT=golint
FMT_LOG=fmt.log
LINT_LOG=lint.log

.DEFAULT_GOAL := test

.PHONY: test
test:
	go test $(UNIT_TEST_PACKAGES) -cover -coverprofile=cover.out

.PHONY: ci
ci: check test

.PHONY: format
format:
	@${GOPATH}/bin/goimports -local "github.com/jaegertracing/jaeger-opentelemetry-collector" -l -w $(shell git ls-files "*\.go" | grep -v vendor)

.PHONY: lint
lint:
	@${GOLINT} ./...

.PHONY: check
check:
	$(shell make format > ${FMT_LOG})
	$(shell make lint > ${LINT_LOG})
	@[ ! -s "$(FMT_LOG)" ] || (echo "Go fmt, code style or import ordering failures, run 'make format'" | cat - $(FMT_LOG) && false)
	@[ ! -s "$(LINT_LOG)" ] || (echo "Go lint failures:" | cat - $(LINT_LOG) && false)

.PHONY: install-tools
install-tools:
	${GO_FLAGS} go install \
		golang.org/x/lint/golint \
		golang.org/x/tools/cmd/goimports
