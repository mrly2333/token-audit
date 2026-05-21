APP := newapi-audit-proxy

.PHONY: build run fmt hash-password

build:
	go build -o bin/$(APP) ./cmd/$(APP)

run:
	go run ./cmd/$(APP) -config config.yaml

fmt:
	go fmt ./...

hash-password:
	go run ./scripts/hash_password.go "change-me"
