export type CalculationOutcome =
  | { kind: 'success'; result: number }
  | { kind: 'rejected'; code: string }
  | { kind: 'unavailable' }

type Operands = { a: number; b?: number }

/** Calls the API. It never throws: every failure is returned as an outcome. */
export async function calculate(
  operation: string,
  operands: Operands,
  signal: AbortSignal,
): Promise<CalculationOutcome> {
  try {
    const response = await fetch(`/api/v1/${operation}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(operands),
      signal,
    })
    const body = await response.json()
    if (typeof body?.result === 'number') {
      return { kind: 'success', result: body.result }
    }
    if (typeof body?.error?.code === 'string') {
      return { kind: 'rejected', code: body.error.code }
    }
    return { kind: 'unavailable' }
  } catch {
    return { kind: 'unavailable' }
  }
}
