.PHONY: test
test:
	go test -v ./tests/... -cover

.PHONY: test-coverage
test-coverage:
	go test -v ./tests/... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html