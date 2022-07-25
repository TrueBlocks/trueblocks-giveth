all:
	@go build main.go
	@mv bin
	@mv main ~/bin/giveth

test:
	@echo Building...
	@make all
	@./test.sh

images:
	@cat R/counts.R | R --no-save

.PHONY: all test clean
