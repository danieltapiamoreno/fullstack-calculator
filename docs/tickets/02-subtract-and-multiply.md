# 02: Subtract and multiply

**What to build:** A user can subtract `b` from `a` or multiply two numbers in the form, and an API consumer can do the same through `POST /api/v1/subtract` and `POST /api/v1/multiply`.

**Blocked by:** 01 (Add: first working calculation)

**Status:** ready-for-agent

- [ ] `POST /api/v1/subtract` with `{"a": 2, "b": 3}` returns `{"result": -1}`
- [ ] `POST /api/v1/multiply` with `{"a": 2, "b": 3}` returns `{"result": 6}`
- [ ] Both operations appear in the picker and show the correct symbol in the form
- [ ] Switching operation keeps the entered values and clears the result, field errors, and any request in flight (ticket 01 built this, but it can't be exercised with one operation)
- [ ] Overflow in either operation returns 422 with code `result_out_of_range`
- [ ] `0 * -1` returns `0`, not `-0`
- [ ] Both operations reject invalid requests with 400, as `add` does
- [ ] Calculator, API, and form tests cover the behaviour above
