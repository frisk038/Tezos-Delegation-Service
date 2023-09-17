# Makefile

.PHONY: test up

# Run tests
test:
	go test ./...

# Start Docker Compose
up:
	docker-compose up -d

# Stop Docker Compose
down:
	docker-compose down
