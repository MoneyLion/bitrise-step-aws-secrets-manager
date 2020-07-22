.SUFFIXES:

OUTFILE=main

.PHONY: all
all: build run

.PHONY: build
build:
	@go build -o $(OUTFILE)

.PHONY: run
run:
	@./$(OUTFILE)

.PHONY: clean
clean:
	@rm $(OUTFILE)

.PHONY: fmt
fmt:
	@go fmt

.PHONY: br
br:
	@bitrise run test
