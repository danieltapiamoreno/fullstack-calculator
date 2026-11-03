# 01: Add: first working calculation

**What to build:** A user can open the app locally, enter two numbers, choose add, submit with `=` or Enter, and see the rounded result. An API consumer can do the same through `POST /api/v1/add`. This first slice also establishes everything later operations share: tool pins and project setup, the operation table, the JSON contract, request validation, the form, and the Make commands.

See the spec for the full API contract and frontend behaviour: `docs/fullstack-calculator-spec.md`.

**Blocked by:** None (can start immediately)

**Status:** done

Setup and tooling:

- [x] Tool pins are updated to Go 1.27 and Node 26, and Vite and Vitest are confirmed to work on Node 26
- [x] The backend builds and tests without a built frontend
- [x] The frontend dev server forwards `/api` requests to the backend on port 8080
- [x] `make test`, `make lint`, and `make build` work, with linting covering `gofmt`, `go vet`, ESLint, and `tsc`

API:

- [x] `POST /api/v1/add` with `{"a": 2, "b": 3}` returns 200 and `{"result": 5}`
- [x] Adding a further operation requires only a calculator function and an operation table entry
- [x] Malformed JSON, a missing field, a wrong type, an unknown field, or `1e999` returns 400 with code `invalid_request`
- [x] A missing `b` returns an error and never becomes `0`
- [x] A request body over 1 KB returns 413 with code `request_too_large`
- [x] `1e308 + 1e308` returns 422 with code `result_out_of_range`, and no non-finite value can reach the JSON encoder
- [x] A result of `-0` is returned as `0`
- [x] Unknown routes and wrong methods return Go's default plain-text 404 and 405 responses

Frontend:

- [x] The operation picker is a row of buttons with radio-group accessibility semantics, large enough to tap reliably
- [x] The inputs and operation symbol read like the expression, and values persist across operation changes
- [x] Empty and malformed inputs show an error under the relevant input without calling the API
- [x] Negative numbers, decimals, and scientific notation are accepted, and the input allows typing a minus sign on iOS
- [x] The `=` button is disabled while a request is in flight
- [x] The result appears below the form, rounded to 12 significant digits, so `0.1 + 0.2` displays as `0.3`
- [x] Changing an input or the operation clears the result, and changing an input clears its error
- [x] `result_out_of_range` appears as a form-level message in the result area
- [x] A network failure appears in a general banner
- [x] Narrow screens use a single-column layout with no sideways scrolling

Tests:

- [x] Calculator, API, form, validation, and formatting tests cover the behaviour above, except switching operation, which needs a second operation and is tested in ticket 02
