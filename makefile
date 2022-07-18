all:
	@go build main.go
	@mv main ~/source/giveth

test:
	@./test.sh

.PHONY: all test clean
