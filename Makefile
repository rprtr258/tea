.PHONY: test
test:
	@go run gotest.tools/gotestsum@latest --format dots-v2

.PHONY: fmt
fmt:
	@gofumpt -l -w .
