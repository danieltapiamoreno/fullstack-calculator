# 06: Percentage

**What to build:** A user can calculate `a` percent of `b`, and an API consumer can do the same through `POST /api/v1/percentage`. The form makes the meaning unambiguous by reading like "20 % of 50".

**Blocked by:** 01 (Add: first working calculation)

**Status:** ready-for-agent

- [ ] `POST /api/v1/percentage` with `{"a": 20, "b": 50}` returns `{"result": 10}`
- [ ] Percentage appears in the picker, and the form reads like "20 % of 50"
- [ ] Negative and decimal values calculate correctly
- [ ] Overflow returns 422 with code `result_out_of_range`
- [ ] Percentage rejects invalid requests with 400, as `add` does
- [ ] Calculator, API, and form tests cover the behaviour above
