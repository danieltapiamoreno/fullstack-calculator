# 04: Power

**What to build:** A user can raise `a` to the power of `b` in the form, and an API consumer can do the same through `POST /api/v1/power`. Results that have no representable answer give a specific, clear error rather than `NaN` or `Infinity`.

**Blocked by:** 01 (Add: first working calculation)

**Status:** done

- [x] `POST /api/v1/power` with `{"a": 2, "b": 3}` returns `{"result": 8}`
- [x] Power appears in the picker and shows the correct symbol in the form
- [x] `0^0` returns `1`
- [x] `0^-1` returns 422 with code `division_by_zero`, and the form shows it under `b`
- [x] `(-8)^0.5` returns 422 with code `undefined_result`, and the form shows it as a form-level message in the result area
- [x] `10^10000` returns 422 with code `result_out_of_range`, and the form shows it as a form-level message in the result area
- [x] Power rejects invalid requests with 400, as `add` does
- [x] Calculator, API, and form tests cover the behaviour above
