export type Field = 'a' | 'b'

export type Operation = {
  /** The API route name. */
  id: string
  label: string
  symbol: string
  /** The input each API error code belongs to. Other codes show as form-level messages. */
  errorFields?: Partial<Record<string, Field>>
}

export const OPERATIONS: readonly Operation[] = [
  { id: 'add', label: 'Add', symbol: '+' },
  { id: 'subtract', label: 'Subtract', symbol: '−' },
  { id: 'multiply', label: 'Multiply', symbol: '×' },
  { id: 'divide', label: 'Divide', symbol: '÷', errorFields: { division_by_zero: 'b' } },
]
