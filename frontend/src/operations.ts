export type Field = 'a' | 'b'

export type Operation = {
  /** The API route name. */
  id: string
  label: string
  symbol: string
  /** Shown on the keypad when it differs from the symbol in the expression. */
  keyLabel?: string
  /** The input each API error code belongs to. Other codes show as form-level messages. */
  errorFields?: Partial<Record<string, Field>>
  /** Takes only `a`: the symbol comes first, and `b` is not shown, validated, or sent. */
  oneInput?: boolean
}

export const OPERATIONS: readonly Operation[] = [
  { id: 'add', label: 'Add', symbol: '+' },
  { id: 'subtract', label: 'Subtract', symbol: '−' },
  { id: 'multiply', label: 'Multiply', symbol: '×' },
  { id: 'divide', label: 'Divide', symbol: '÷', errorFields: { division_by_zero: 'b' } },
  { id: 'power', label: 'Power', symbol: '^', errorFields: { division_by_zero: 'b' } },
  {
    id: 'sqrt',
    label: 'Square root',
    symbol: '√',
    errorFields: { undefined_result: 'a' },
    oneInput: true,
  },
  { id: 'percentage', label: 'Percentage', symbol: '% of', keyLabel: '%' },
]
