# 06: Percentage

**What to build:** A user can calculate `a` percent of `b`, and an API consumer can do the same through `POST /api/v1/percentage`. The form makes the meaning unambiguous by reading like "20 % of 50".

**Blocked by:** 01 (Add: first working calculation)

**Status:** done

- [x] `POST /api/v1/percentage` with `{"a": 20, "b": 50}` returns `{"result": 10}`
- [x] Percentage appears in the picker, and the form reads like "20 % of 50"
- [x] Negative and decimal values calculate correctly
- [x] Overflow returns 422 with code `result_out_of_range`
- [x] Percentage rejects invalid requests with 400, as `add` does
- [x] Calculator, API, and form tests cover the behaviour above
