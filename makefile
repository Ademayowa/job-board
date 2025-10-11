# Makefile
.PHONY: test test-coverage

test:
	go test ./tests/... -v

test-coverage:
	go test ./tests/... -v -coverprofile=coverage.out
	go tool cover -func=coverage.out