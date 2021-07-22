.SUFFIXES:

.PHONY: fmt
fmt:
	@go fmt

.PHONY: br
br:
	@bitrise run test
