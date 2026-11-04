export type Operation = {
  /** The API route name. */
  id: string
  label: string
  symbol: string
}

export const OPERATIONS: readonly Operation[] = [
  { id: 'add', label: 'Add', symbol: '+' },
  { id: 'subtract', label: 'Subtract', symbol: '−' },
  { id: 'multiply', label: 'Multiply', symbol: '×' },
]
