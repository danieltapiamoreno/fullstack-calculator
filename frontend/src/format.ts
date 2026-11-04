const DISPLAY_PRECISION = 12

/** Rounds a result for display, hiding floating-point noise such as 0.30000000000000004. */
export function formatResult(value: number): string {
  return String(Number(value.toPrecision(DISPLAY_PRECISION)))
}
