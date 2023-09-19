bin:
	mkdir -p bin

.PHONY: clean
clean: ## remove bin directory
	@rm -r bin 2>/dev/null || :

.PHONY: up
up:
	docker-compose up --build --detach

.PHONY: down
down:
	docker-compose down -v

.PHONY: build
build:
	mkdir -p bin && go build -o bin/xepelin-bank ./cmd/server/main.go