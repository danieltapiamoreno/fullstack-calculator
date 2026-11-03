.PHONY: test lint build

test: frontend/node_modules
	cd backend && go test ./...
	cd frontend && npm test

lint: frontend/node_modules
	@unformatted="$$(cd backend && gofmt -l .)"; \
	if [ -n "$$unformatted" ]; then echo "gofmt needed:"; echo "$$unformatted"; exit 1; fi
	cd backend && go vet ./...
	cd frontend && npm run lint
	cd frontend && npm run typecheck

build: frontend/node_modules
	cd backend && go build -o bin/server ./cmd/server
	cd frontend && npm run build

frontend/node_modules: frontend/package-lock.json
	cd frontend && npm ci
	@touch $@
