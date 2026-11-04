.PHONY: test coverage lint build docker-run docker-up

test: frontend/node_modules
	cd backend && go test ./...
	cd frontend && npm test

coverage: frontend/node_modules
	cd backend && go test -coverprofile=coverage.out ./...
	cd backend && go tool cover -func=coverage.out
	cd backend && go tool cover -html=coverage.out -o coverage.html
	cd frontend && npm run coverage

lint: frontend/node_modules
	@unformatted="$$(cd backend && gofmt -l .)"; \
	if [ -n "$$unformatted" ]; then echo "gofmt needed:"; echo "$$unformatted"; exit 1; fi
	cd backend && go vet ./...
	cd frontend && npm run lint
	cd frontend && npm run typecheck

build: frontend/node_modules
	cd backend && go build -o bin/server ./cmd/server
	cd frontend && npm run build

docker-run:
	docker build -t fullstack-calculator .
	docker run --rm -p 8080:8080 fullstack-calculator

docker-up:
	docker compose up --build

frontend/node_modules: frontend/package-lock.json
	cd frontend && npm ci
	@touch $@
