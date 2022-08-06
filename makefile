all:
	@echo building...
	@go build main.go
	@mv main ~/source/giveth

test:
	@echo Building...
	@make all
	@./test.sh

images:
	@cat R/counts.R | R --no-save

.PHONY: all test clean
