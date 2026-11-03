# Full-Stack Calculator

A calculator with a Go REST API and a React + TypeScript frontend. It supports addition, subtraction, multiplication, division, exponentiation, square root, and percentage.

All calculation happens on the backend. The frontend collects input, sends it to the API, and displays the result or error.

```
.
├── backend/
│   ├── cmd/server/          # Entry point: configuration, wiring, graceful shutdown
│   └── internal/
│       ├── api/             # Routing, request decoding, error-to-status mapping
│       └── calculator/      # Pure maths and typed errors, with no HTTP code
├── frontend/                # React + TypeScript app built with Vite
├── Dockerfile               # Builds one image containing the API and frontend
├── Makefile                 # test, lint, build, and docker-run commands
└── mise.toml                # Pinned Go and Node versions
```

## Setup

### Run with Docker

You only need Docker for this option.

```sh
make docker-run
```

Or run the equivalent commands without Make:

```sh
docker build -t fullstack-calculator .
docker run --rm -p 8080:8080 fullstack-calculator
```

Then open <http://localhost:8080>.

### Run locally

You need Go 1.27 and Node 26. If you use [mise](https://mise.jdx.dev), `mise install` installs both.

Start the backend in one terminal:

```sh
cd backend
go run ./cmd/server
```

Start the frontend in another:

```sh
cd frontend
npm ci
npm run dev
```

Then open <http://localhost:5173>. Vite forwards `/api` requests to the backend on port 8080, so CORS isn't needed.

### Configuration

The server reads two environment variables:

| Variable | Default | Purpose |
|---|---|---|
| `PORT` | `8080` | Port the server listens on |
| `STATIC_DIR` | Unset | Directory of built frontend files to serve. When unset, the server runs API-only. |

### Common commands

| Command | What it does |
|---|---|
| `make test` | Runs the backend and frontend tests |
| `make lint` | Runs `gofmt`, `go vet`, ESLint, and `tsc` |
| `make build` | Builds the backend and frontend |
| `make docker-run` | Builds the Docker image and runs it on port 8080 |

## API

All calculator routes use `POST`, and both requests and responses are JSON.

| Route | Body | Result |
|---|---|---|
| `/api/v1/add` | `{"a": 2, "b": 3}` | `a + b` |
| `/api/v1/subtract` | `{"a": 2, "b": 3}` | `a - b` |
| `/api/v1/multiply` | `{"a": 2, "b": 3}` | `a * b` |
| `/api/v1/divide` | `{"a": 6, "b": 3}` | `a / b` |
| `/api/v1/power` | `{"a": 2, "b": 3}` | `a` raised to the power of `b` |
| `/api/v1/percentage` | `{"a": 20, "b": 50}` | `a` percent of `b` |
| `/api/v1/sqrt` | `{"a": 9}` | Square root of `a` |

`GET /healthz` returns `200` with `{"status":"ok"}`.

### Successful response

```sh
curl -X POST http://localhost:8080/api/v1/percentage \
  -H 'Content-Type: application/json' \
  -d '{"a": 20, "b": 50}'
```

```json
{"result": 10}
```

### Error response

```sh
curl -X POST http://localhost:8080/api/v1/divide \
  -H 'Content-Type: application/json' \
  -d '{"a": 1, "b": 0}'
```

```json
{"error": {"code": "division_by_zero", "message": "cannot divide by zero"}}
```

The `code` value is stable and meant for programs to check. The `message` value is for people to read and may change.

| Status | Code | When it's returned |
|---|---|---|
| 400 | `invalid_request` | Malformed JSON, missing fields, wrong types, unknown fields (including `b` on `sqrt`), or numbers too large for `float64`, such as `1e999` |
| 413 | `request_too_large` | The request body is larger than 1 KB |
| 422 | `division_by_zero` | Dividing by zero, or raising zero to a negative power, such as `0^-1` |
| 422 | `undefined_result` | The result isn't a real number, such as the square root of a negative number or `(-8)^0.5` |
| 422 | `result_out_of_range` | The result is too large for `float64`, such as `10^10000` |

Unknown routes and wrong methods are the exception to this format: they return Go's default plain-text `404` and `405` responses.

### Edge cases

- `0^0` returns `1`, matching Go's `math.Pow`.
- A result of `-0` is returned as `0`.

## Design rationale

### Single-operation form instead of a button keypad

The UI has number inputs, a row of operation buttons, and an `=` button. Each submission makes one API call and shows one result. This handles one-input and two-input operations cleanly and makes per-field validation natural. A keypad with chained input would need a client-side state machine and rules for operator precedence, without adding any of the required functionality.

### One route per operation

Each operation's required fields are clear from its route, and the routes map directly to the required operations. They share one handler driven by a table of operations, so adding a route takes little code.

### `float64` instead of a decimal library

A decimal library would make `0.1 + 0.2` exact, but square roots and fractional exponents still can't be exact, so the API would end up with mixed precision rules. `float64` needs no dependencies and behaves consistently.

The trade-off is standard floating-point behaviour: the API returns `0.30000000000000004` for `0.1 + 0.2`. The UI rounds results to 12 significant digits for display. This approach isn't suitable for financial calculations that need exact decimal results.

### Separate 400 and 422 errors

A 400 means the request itself is wrong and the client needs to fix it. A 422 means the request is well formed, but the maths has no representable answer. Non-finite results have to be errors because JSON can't represent `Inf` or `NaN`.

Request fields decode into pointers so that a missing `b` returns an error instead of silently becoming `0`.

### Maths rules live only on the server

The client checks that inputs are present and are well-formed numbers. Rules such as "can't divide by zero" live only in the backend, so there's a single source of truth. The UI maps stable error codes to the relevant field, which gives inline errors without parsing message text.

Inputs use `type="text"` because `type="number"` reports an empty value for partial input, and some iOS decimal keyboards have no minus key.

### Pure calculator package

`internal/calculator` knows nothing about HTTP, so the maths can be tested directly with table-driven tests. `internal/api` is tested separately through `httptest`.

### One image and one origin

Go serves both the API and the built frontend from a single port, so deployment is one `docker run` with no nginx, Compose, or production CORS configuration.

The server reads frontend files from `STATIC_DIR` rather than using `go:embed`. That way, `go build` and `go test` don't depend on a built frontend. The app has no client-side routes, so a plain file server is enough. The final image is distroless and runs as a non-root user.

### Minimal dependencies

The backend uses only the Go standard library, including `net/http` method-and-path routing. The frontend uses plain CSS and `useState`, because a single form doesn't need a UI or state management library.

### Testing scope

- **Backend:** table-driven tests for the maths, and `httptest` tests for routes, status codes, and error bodies.
- **Frontend:** Vitest and React Testing Library tests for validation, display formatting, and the full form, with `fetch` stubbed.

There's no end-to-end test suite, because the API and form tests cover the key behaviour without a browser-driven suite to maintain.

### Intentionally left out

History, authentication, persistence, metrics, and CI aren't included because the calculator's requirements don't need them.
