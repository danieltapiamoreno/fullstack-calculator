export type CalculationOutcome =
  | { kind: 'success'; result: number }
  | { kind: 'rejected'; code: string }
  | { kind: 'unavailable' }

/** Calls the API. It never throws: every failure is returned as an outcome. */
export async function calculate(
  operation: string,
  operands: Record<string, number>,
): Promise<CalculationOutcome> {
  try {
    const response = await fetch(`/api/v1/${operation}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(operands),
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
