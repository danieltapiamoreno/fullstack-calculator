export type Operation = {
  /** The API route name. */
  id: string
  label: string
  symbol: string
}

export const OPERATIONS: readonly Operation[] = [{ id: 'add', label: 'Add', symbol: '+' }]
