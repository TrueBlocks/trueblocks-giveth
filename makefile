all:
	@go build main.go
	@mv main ~/source/giveth

test:
	@./test

.PHONY: all test clean