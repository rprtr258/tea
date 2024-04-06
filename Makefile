## DEV

.PHONY: test
test:
	# @go run gotest.tools/gotestsum@latest --format dots-v2
	@gotestsum --format dots-v2

.PHONY: fmt
fmt:
	@gofumpt -l -w .
	@gci write --skip-generated -s standard -s default -s 'prefix(github.com/rprtr258/tea)' --custom-order .

## RUN EXAMPLES

.PHONY: run-examples
run-examples:
	@go run cmd/main.go

.PHONY: run-tutorials
run-tutorials:
	@go run cmd/main.go tutorials

.PHONY: run-styles
run-styles:
	@go run cmd/main.go styles

.PHONY: run-markdown
run-markdown:
	@go run cmd/main.go markdown
