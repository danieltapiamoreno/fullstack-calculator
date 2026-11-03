---
title: Full-stack calculator
status: ready-for-agent
labels:
  - ready-for-agent
---

# Full-stack calculator

## Problem Statement

A user needs a simple, dependable way to perform calculations in a browser on a desktop or phone, including add, subtract, multiply, divide, exponentiation, square root, and percentage. When an input is invalid or a calculation has no answer, the user needs a clear explanation rather than a blank screen, `Infinity`, or `NaN`.

API consumers need the same operations over REST, with predictable JSON responses and error codes they can rely on.

Reviewers need to be able to run the project with little setup, read clean and idiomatic code in both layers, confirm its behaviour through unit tests, and understand why it was designed this way.

## Solution

A monorepo containing a Go REST API and a React + TypeScript frontend built with Vite, deployed as one Docker image that serves both from a single origin and port.

The user chooses an operation, enters one or two numbers in a form that reads like the expression, and submits it with `=` or Enter. The backend performs every calculation. The frontend shows the result or places the error next to the relevant input.

A root README documents setup, API usage, and design rationale.

## User Stories

### Calculating

1. As a calculator user, I want to add two numbers, so that I can get their sum.
2. As a calculator user, I want to subtract one number from another, so that I can get their difference.
3. As a calculator user, I want to multiply two numbers, so that I can get their product.
4. As a calculator user, I want to divide one number by another, so that I can get their quotient.
5. As a calculator user, I want to raise a number to a power, so that I can calculate exponents.
6. As a calculator user, I want to take the square root of a number, so that I can get its root without entering a second value.
7. As a calculator user, I want to calculate `a` percent of `b`, so that I can answer questions like "what is 20% of 50?".
8. As a calculator user, I want to enter negative numbers, so that I can calculate with values below zero.
9. As a calculator user, I want to enter decimal numbers, so that I can calculate with fractional values.
10. As a calculator user, I want to enter scientific notation such as `1e5`, so that I can work with very large or very small values.

### Using the form

11. As a calculator user, I want to see every operation as a button in one row, so that I can compare and switch between them quickly.
12. As a calculator user, I want the form to read like the expression I'm calculating, so that it's obvious which value goes where.
13. As a calculator user, I want the percentage operation to be labelled in the form "20 % of 50", so that its meaning is unambiguous.
14. As a calculator user, I want the second input hidden for square root, so that I'm not asked for a value the operation doesn't use.
15. As a calculator user, I want my entered values kept when I switch operations, so that I can try another operation without retyping them.
16. As a calculator user, I want to submit with the `=` button, so that I can choose when the calculation runs.
17. As a keyboard user, I want to submit with Enter, so that I don't need to reach for the mouse.
18. As a keyboard or screen reader user, I want the operation picker to behave like a radio group, so that I can tell which operation is selected and change it with standard controls.
19. As a calculator user, I want the `=` button disabled while a request is in flight, so that I can't accidentally submit twice.

### Reading results

20. As a calculator user, I want the result shown prominently below the form, so that I can find it right away.
21. As a calculator user, I want results rounded to 12 significant digits for display, so that I see `0.3` rather than `0.30000000000000004`.
22. As a calculator user, I want the result cleared when I change an input or operation, so that I never mistake an old answer for the answer to new values.

### Handling errors

23. As a calculator user, I want to be told when a required input is empty, so that I know what to fill in.
24. As a calculator user, I want to be told when an input isn't a valid number, so that I can correct it before submitting.
25. As a calculator user, I want to see a clear error under the relevant input when I divide by zero, so that I know which value to change.
26. As a calculator user, I want a clear error when I take the square root of a negative number, so that I understand why there is no result.
27. As a calculator user, I want a clear error when a power has no real result, so that I'm not shown `NaN`.
28. As a calculator user, I want a clear error when a result is too large to represent, so that I'm not shown `Infinity`.
29. As a calculator user, I want a general banner when the server can't be reached, so that I can tell a connection failure from a problem with my input.
30. As a calculator user, I want a previous error cleared when I change the related input, so that I'm not looking at an outdated message.

### Using a phone

31. As a mobile user, I want a single-column layout on a narrow screen, so that I can use the calculator without zooming or scrolling sideways.
32. As an iOS user, I want a keyboard that lets me type a minus sign, so that I can enter negative numbers.
33. As a mobile user, I want operation buttons large enough to tap reliably, so that I don't select the wrong one.

### Consuming the API

34. As an API consumer, I want one route per operation, so that each operation's required fields are clear from its route.
35. As an API consumer, I want successful calculations returned as `{"result": <number>}`, so that every operation can be parsed the same way.
36. As an API consumer, I want the raw `float64` result, so that display rounding remains my decision.
37. As an API consumer, I want all errors from calculator routes in one JSON shape, so that I can handle them with one code path.
38. As an API consumer, I want stable error codes, so that I can branch on them without parsing message text.
39. As an API consumer, I want human-readable error messages, so that I can diagnose failures while integrating.
40. As an API consumer, I want a 400 for malformed JSON, so that I know the request needs correcting.
41. As an API consumer, I want a 400 when a required field is missing, so that a missing value never silently becomes `0`.
42. As an API consumer, I want a 400 for fields of the wrong type, so that invalid data is never calculated.
43. As an API consumer, I want a 400 for unknown fields, including `b` on square root, so that mistakes in my request don't go unnoticed.
44. As an API consumer, I want a 400 for numbers too large for `float64`, so that they aren't silently changed.
45. As an API consumer, I want a 413 for oversized request bodies, so that the limit is explicit.
46. As an API consumer, I want a 422 when valid input has no representable mathematical answer, so that I can tell it apart from a malformed request.
47. As an API consumer, I want `0^-1` reported as `division_by_zero`, so that the error describes the actual cause.
48. As an API consumer, I want `0^0` to return `1`, so that the behaviour matches Go's `math.Pow` and is documented.
49. As an API consumer, I want `-0` returned as `0`, so that I don't need to handle negative zero.
50. As an operator, I want a health check route, so that I can tell whether the service is running.

### Running and reviewing the project

51. As a reviewer, I want to start the whole app with one Docker command, so that I can try it without installing Go or Node.
52. As a developer, I want the frontend dev server to forward API requests to the backend, so that local development doesn't need CORS configuration.
53. As a developer, I want the backend to build and test without a built frontend, so that backend work stays independent.
54. As a developer, I want pinned tool versions, so that my local environment matches the Docker image.
55. As a developer, I want Make commands for testing, linting, building, and running the image, so that I don't need mise or need to remember each command.
56. As an operator, I want the port and static file directory configurable through environment variables, so that I can adapt the deployment without rebuilding.
57. As an operator, I want the server to shut down gracefully, so that in-flight requests can finish.
58. As an operator, I want the final image to be minimal and run as non-root, so that the deployment has a small attack surface.
59. As a reviewer, I want a README covering setup, API usage, and design rationale, so that I can understand the project without reading all the code.
60. As a reviewer, I want unit tests for key behaviour in both layers, so that I can trust the project works and is safe to change.
61. As a reviewer, I want a clear history of logical commits, so that I can follow how the project was built.

## Implementation Decisions

### Technology

- The backend uses Go 1.27 and only the standard library, including `net/http` method-and-path routing.
- The frontend uses Node 26, npm, React, TypeScript, and Vite, with plain CSS and `useState`. It has no UI or state management libraries.
- Go has no LTS programme, so 1.27 was chosen as the newest supported release. Node 26 is the current LTS line.
- Tool versions are pinned consistently across mise, the Go module, and the Docker image. The existing mise pins (Go 1.23 and Node 22) must be updated.

### Backend modules

- **Calculator module:** a deep, pure module containing all maths and typed domain errors. It has no knowledge of HTTP or JSON.
- **API module:** owns routing, request decoding and validation, and the mapping from domain errors to status codes and JSON. A single handler is driven by a table of operations, rather than one handwritten handler per route.
- **Server entry point:** wiring only. It handles configuration, server timeouts, and graceful shutdown.
- The Go module path is `fullstack-calculator/backend`, because the repo has no remote yet.

### API contract

All calculator routes use `POST` under `/api/v1`.

| Route | Required fields | Result |
|---|---|---|
| `add` | `a`, `b` | `a + b` |
| `subtract` | `a`, `b` | `a - b` |
| `multiply` | `a`, `b` | `a * b` |
| `divide` | `a`, `b` | `a / b` |
| `power` | `a`, `b` | `a` raised to the power of `b` |
| `percentage` | `a`, `b` | `a` percent of `b` |
| `sqrt` | `a` | Square root of `a` |

- Success responses have the shape `{"result": <number>}`.
- Error responses have the shape `{"error": {"code": <string>, "message": <string>}}`.
- `GET /healthz` returns 200 with `{"status": "ok"}`.
- Unknown routes and wrong methods keep Go's default plain-text 404 and 405 responses. This is the one exception to the error shape, and the README documents it.

### Error handling

| Status | Code | Condition |
|---|---|---|
| 400 | `invalid_request` | Malformed JSON, missing fields, wrong types, unknown fields, or numbers too large for `float64` |
| 413 | `request_too_large` | Request body over 1 KB |
| 422 | `division_by_zero` | Dividing by zero, or raising zero to a negative power |
| 422 | `undefined_result` | The result isn't a real number, such as the square root of a negative number or `(-8)^0.5` |
| 422 | `result_out_of_range` | The result overflows `float64` |

- A 400 means the request itself is wrong. A 422 means the request is well formed, but the maths has no representable answer.
- Request fields decode in a way that distinguishes a missing value from zero, and unknown fields are rejected.
- JSON can't represent `Inf` or `NaN`, so a non-finite result must never reach the encoder.

### Numeric behaviour

- All calculation uses `float64`, and the API returns the raw value. The precision trade-off is documented in the README.
- `0^0` returns `1`.
- A result of `-0` is returned as `0`.
- Results are rounded only for display, in the frontend, to 12 significant digits.

### Frontend behaviour

- The UI is a single form, not a button keypad with chained input. Each submission makes one API call and shows one result.
- The operation picker is a row of buttons with radio-group accessibility semantics.
- The inputs and operation symbol are laid out like the expression, and the second input is hidden for square root.
- Values persist across operation changes. Changing an input or the operation clears the current result.
- Calculation runs only on `=` or Enter, with no calculate-as-you-type. The submit button is disabled while a request is in flight.
- Client-side validation checks only that inputs are present and are well-formed numbers, including a leading minus sign, decimals, and scientific notation.
- Maths rules live only on the server. The frontend maps stable error codes to fields, such as showing `division_by_zero` under `b` and a negative square root's `undefined_result` under `a`.
- Errors that don't belong to one input appear as a form-level message in the result area. These are `result_out_of_range` from any operation and `undefined_result` from `power`.
- Network failures appear in a general banner.
- Inputs use `type="text"` and avoid `inputmode="decimal"`, because some iOS keyboards have no minus key.
- Narrow screens use a single-column layout.
- The frontend has three deep modules: an API client that turns responses into typed outcomes, input validation, and result formatting. The form component composes them.

### Deployment

- One multi-stage root Dockerfile builds the frontend and backend into one distroless image that runs as non-root.
- The Go server serves both the API and the built frontend from one origin and port, so production needs no CORS configuration.
- Frontend files are read from a directory set by `STATIC_DIR`, rather than compiled in with `go:embed`, so backend builds and tests don't depend on a built frontend. The server runs API-only when it's unset.
- `PORT` defaults to 8080.
- The app has no client-side routes, so a plain file server is enough and no SPA fallback is needed.
- In local development, the Vite dev server forwards `/api` requests to the backend on port 8080.

### Tooling and documentation

- A root Makefile provides `test`, `lint`, `build`, and `docker-run`.
- Linting uses `gofmt`, `go vet`, ESLint, and `tsc`.
- The root README covers setup, API usage, and design rationale. It already exists and must stay accurate as the code is built.
- Work is delivered as logical commits on `main`, and nothing is pushed. This is an exception to working on a branch because the repo has no commits yet.

## Testing Decisions

### What makes a good test

- A test exercises behaviour through a module's public interface and doesn't depend on its internals.
- A test should survive refactoring and fail only when behaviour changes.
- Tests cover the edge cases in this spec, not just normal inputs.
- External boundaries are stubbed only where they cross a process boundary, which on the frontend means `fetch`.

### Test seams

| Seam | Tool | What it covers |
|---|---|---|
| API module over HTTP | Go `testing` and `httptest` | Routes, status codes, success and error bodies, missing and unknown fields, wrong types, oversized bodies, and non-finite results |
| Calculator module | Table-driven Go tests | Every operation, its edge cases, and its typed errors |
| Form component | Vitest and React Testing Library, with `fetch` stubbed | Choosing operations, hiding the second input, submitting, in-flight state, showing results, clearing stale results, mapping errors to fields, and the network failure banner |
| Input validation and result formatting | Vitest | Accepted and rejected input formats, and rounding to 12 significant digits |

The API and form seams are the highest available, and they carry most of the behaviour. The calculator, validation, and formatting seams exist so numeric edge cases can be covered directly and cheaply.

### Prior art

The repo has no existing tests. These tests set the conventions, following standard Go table-driven testing and React Testing Library's user-focused queries.

### Verification

There is no end-to-end test suite. After implementation, the full stack is verified manually by building the Docker image, exercising the running container with `curl`, and reporting the results.

## Out of Scope

- Calculation history
- Authentication and authorisation
- Persistence or a database
- Metrics and observability beyond the health check
- CI pipelines
- Browser-driven end-to-end tests
- A button keypad, chained expressions, or operator precedence
- Exact decimal arithmetic
- Rate limiting
- CORS configuration
- Custom JSON responses for unknown routes and wrong methods
- Pushing to a remote or setting one up

## Further Notes

- Grilling settled the design, and the root README records the same decisions for readers of the finished project.
- These details were filled in during drafting rather than decided individually in the grilling session, so they're easy to change: the `power` route name, the `invalid_request` and `request_too_large` codes, the `{"status": "ok"}` health response, and the exact wording of error messages.
- Check during setup that Vite and Vitest work on Node 26.
- Commits end with the `Co-Authored-By` attribution line unless the user asks to leave it out.
