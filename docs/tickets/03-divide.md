# 03: Divide

**What to build:** A user can divide two numbers in the form, and an API consumer can do the same through `POST /api/v1/divide`. Dividing by zero gives a clear error next to the input that caused it.

**Blocked by:** 01 (Add: first working calculation)

**Status:** ready-for-agent

- [ ] `POST /api/v1/divide` with `{"a": 6, "b": 3}` returns `{"result": 2}`
- [ ] Division appears in the picker and shows the correct symbol in the form
- [ ] A zero `b` returns 422 with code `division_by_zero`, including when `a` is also zero
- [ ] The form shows `division_by_zero` under `b`
- [ ] Overflow returns 422 with code `result_out_of_range`
- [ ] Division rejects invalid requests with 400, as `add` does
- [ ] Calculator, API, and form tests cover the behaviour above
