.PHONY: build test run docker-build docker-buildx docker-run clean

BINARY=docsbot

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o bin/$(BINARY) ./cmd/docsbot

run: build
	./bin/$(BINARY)

test:
	go test ./... -v

# обычный docker build (архитектура хоста)
docker-build:
	docker build -t vasenin26/docsbot:latest .

# multi-arch build (требует docker buildx и push настроенного реестра)
# пример: make docker-buildx PLATFORMS=linux/amd64,linux/arm64 TAG=vasenin26/docsbot:latest
docker-buildx:
	docker buildx build --platform $(PLATFORMS) --tag $(TAG) --push .

docker-run:
	docker run --rm -p 9090:9090 vasenin26/docsbot:latest

clean:
	rm -rf bin
