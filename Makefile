PROJECT_NAME := "blacklist"
VERSION := $(shell cat ./VERSION)

.PHONY: test test-ci cover cover-ci coverhtml lint race version msan examples

test:
	@go test -short

test-ci:
	@go test -v

cover:
	@go test -cover

cover-ci:
	@go test -v -coverprofile=.testCoverage.txt

coverhtml:
	@go test -coverprofile=coverage.out
	@go tool cover -html=coverage.out
	@rm -f ./coverage.out

lint:
	@golint -set_exit_status

race:
	@go test -race -short

msan:
	@go test -msan -short

version:
	@cat VERSION

examples:
	$(MAKE) -C examples
