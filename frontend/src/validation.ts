const NUMBER_FORMAT = /^-?(\d+\.?\d*|\.\d+)([eE][+-]?\d+)?$/

type ParsedOperand = { ok: true; value: number } | { ok: false; message: string }

/** Parses what the user typed. Maths rules, such as dividing by zero, are left to the server. */
export function parseOperand(input: string): ParsedOperand {
  const trimmed = input.trim()
  if (trimmed === '') {
    return { ok: false, message: 'Enter a number.' }
  }
  if (!NUMBER_FORMAT.test(trimmed)) {
    return { ok: false, message: 'Enter a valid number, like 12, -3.5, or 1e5.' }
  }
  const value = Number(trimmed)
  if (!Number.isFinite(value)) {
    return { ok: false, message: 'This number is too large.' }
  }
  return { ok: true, value }
}
