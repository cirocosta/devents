infra:
	cd ./infra && \
		docker-compose up -d

image:
	docker build -t cirocosta/devents .

install:
	go install -v

fmt:
	gofmt -s -w ./main.go
	find ./lib -name "*.go" -exec gofmt -s -w {} \;

toc:
	doctoc ./README.md

.PHONY: install build fmt image toc infra
