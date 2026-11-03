# 01: Add: first working calculation

**What to build:** A user can open the app locally, enter two numbers, choose add, submit with `=` or Enter, and see the rounded result. An API consumer can do the same through `POST /api/v1/add`. This first slice also establishes everything later operations share: tool pins and project setup, the operation table, the JSON contract, request validation, the form, and the Make commands.

See the spec for the full API contract and frontend behaviour: `docs/fullstack-calculator-spec.md`.

**Blocked by:** None (can start immediately)

**Status:** ready-for-agent

Setup and tooling:

- [ ] Tool pins are updated to Go 1.27 and Node 26, and Vite and Vitest are confirmed to work on Node 26
- [ ] The backend builds and tests without a built frontend
- [ ] The frontend dev server forwards `/api` requests to the backend on port 8080
- [ ] `make test`, `make lint`, and `make build` work, with linting covering `gofmt`, `go vet`, ESLint, and `tsc`

API:

- [ ] `POST /api/v1/add` with `{"a": 2, "b": 3}` returns 200 and `{"result": 5}`
- [ ] Adding a further operation requires only a calculator function and an operation table entry
- [ ] Malformed JSON, a missing field, a wrong type, an unknown field, or `1e999` returns 400 with code `invalid_request`
- [ ] A missing `b` returns an error and never becomes `0`
- [ ] A request body over 1 KB returns 413 with code `request_too_large`
- [ ] `1e308 + 1e308` returns 422 with code `result_out_of_range`, and no non-finite value can reach the JSON encoder
- [ ] A result of `-0` is returned as `0`
- [ ] Unknown routes and wrong methods return Go's default plain-text 404 and 405 responses

Frontend:

- [ ] The operation picker is a row of buttons with radio-group accessibility semantics, large enough to tap reliably
- [ ] The inputs and operation symbol read like the expression, and values persist across operation changes
- [ ] Empty and malformed inputs show an error under the relevant input without calling the API
- [ ] Negative numbers, decimals, and scientific notation are accepted, and the input allows typing a minus sign on iOS
- [ ] The `=` button is disabled while a request is in flight
- [ ] The result appears below the form, rounded to 12 significant digits, so `0.1 + 0.2` displays as `0.3`
- [ ] Changing an input or the operation clears the result, and changing an input clears its error
- [ ] `result_out_of_range` appears as a form-level message in the result area
- [ ] A network failure appears in a general banner
- [ ] Narrow screens use a single-column layout with no sideways scrolling

Tests:

- [ ] Calculator, API, form, validation, and formatting tests cover the behaviour above
