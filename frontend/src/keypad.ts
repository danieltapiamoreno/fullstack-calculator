/** Adds a digit or decimal point to a number being entered on the keypad. */
export function appendEntry(value: string, entry: string): string {
  if (entry !== '.') return value + entry
  if (/[.eE]/.test(value)) return value
  return value === '' || value === '-' ? `${value}0.` : `${value}.`
}

export function toggleSign(value: string): string {
  return value.startsWith('-') ? value.slice(1) : `-${value}`
}

export function deleteLast(value: string): string {
  return value.slice(0, -1)
}
