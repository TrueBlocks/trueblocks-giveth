all:
	@go build main.go
	@mv main ~/source/giveth

test:
	@./test.sh

images:
	@cat R/counts.R | R --no-save

.PHONY: all test clean
