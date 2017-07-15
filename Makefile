image:
	docker build -t cirocosta/devents .

install:
	cd ./devents && go install -v

build:
	cd ./devents && go build -v

fmt:
	cd ./devents && gofmt -s -w .
	cd ./lib && gofmt -s -w .

toc:
	doctoc ./README.md

.PHONY: install build fmt image toc
