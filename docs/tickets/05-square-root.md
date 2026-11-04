# 05: Square root

**What to build:** A user can take the square root of a number using one input, and an API consumer can do the same through `POST /api/v1/sqrt`. This is the only operation that takes one operand, so the form and API must both handle that case.

**Blocked by:** 01 (Add: first working calculation)

**Status:** done

- [x] `POST /api/v1/sqrt` with `{"a": 9}` returns `{"result": 3}`
- [x] Square root appears in the picker and hides the second input
- [x] A value entered for `b` is kept and reappears when switching back to a two-input operation
- [x] A hidden `b` is neither validated nor sent to the API
- [x] A negative `a` returns 422 with code `undefined_result`, and the form shows it under `a`
- [x] A request containing `b` returns 400 with code `invalid_request`
- [x] A missing `a` and other invalid requests return 400, as `add` does
- [x] Calculator, API, and form tests cover the behaviour above
