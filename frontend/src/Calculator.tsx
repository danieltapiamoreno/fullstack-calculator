import { useRef, useState } from 'react'
import type { FormEvent } from 'react'
import { calculate } from './api'
import { formatResult } from './format'
import { OPERATIONS } from './operations'
import { parseOperand } from './validation'

type Field = 'a' | 'b'

type Status =
  | { kind: 'idle' }
  | { kind: 'pending' }
  | { kind: 'result'; value: number }
  | { kind: 'rejected'; code: string }
  | { kind: 'unavailable' }

const IDLE: Status = { kind: 'idle' }

const FIELD_LABELS: Record<Field, string> = { a: 'First number', b: 'Second number' }

const REJECTION_MESSAGES: Record<string, string> = {
  result_out_of_range: 'The result is too large to represent.',
}

export function Calculator() {
  const [operationId, setOperationId] = useState(OPERATIONS[0].id)
  const [values, setValues] = useState<Record<Field, string>>({ a: '', b: '' })
  const [fieldErrors, setFieldErrors] = useState<Partial<Record<Field, string>>>({})
  const [status, setStatus] = useState<Status>(IDLE)
  // Identifies the latest request, so a response to outdated input can be ignored.
  const latestRequest = useRef(0)

  const operation = OPERATIONS.find(({ id }) => id === operationId)!

  function resetStatus() {
    latestRequest.current += 1
    setStatus(IDLE)
  }

  function handleOperationChange(id: string) {
    setOperationId(id)
    setFieldErrors({})
    resetStatus()
  }

  function handleValueChange(field: Field, value: string) {
    setValues((current) => ({ ...current, [field]: value }))
    setFieldErrors((current) => ({ ...current, [field]: undefined }))
    resetStatus()
  }

  async function handleSubmit(event: FormEvent) {
    event.preventDefault()
    if (status.kind === 'pending') return

    const a = parseOperand(values.a)
    const b = parseOperand(values.b)
    if (!a.ok || !b.ok) {
      setFieldErrors({ a: a.ok ? undefined : a.message, b: b.ok ? undefined : b.message })
      return
    }

    latestRequest.current += 1
    const request = latestRequest.current
    setFieldErrors({})
    setStatus({ kind: 'pending' })

    const outcome = await calculate(operation.id, { a: a.value, b: b.value })
    if (request !== latestRequest.current) return

    if (outcome.kind === 'success') {
      setStatus({ kind: 'result', value: outcome.result })
    } else {
      setStatus(outcome)
    }
  }

  function renderInput(field: Field) {
    const error = fieldErrors[field]
    const errorId = `${field}-error`
    return (
      <div className="field">
        <label htmlFor={field}>{FIELD_LABELS[field]}</label>
        <input
          id={field}
          type="text"
          autoComplete="off"
          autoCapitalize="off"
          autoCorrect="off"
          spellCheck={false}
          value={values[field]}
          aria-invalid={error !== undefined}
          aria-describedby={error === undefined ? undefined : errorId}
          onChange={(event) => handleValueChange(field, event.target.value)}
        />
        {error !== undefined && (
          <p id={errorId} className="field-error">
            {error}
          </p>
        )}
      </div>
    )
  }

  return (
    <main className="calculator">
      <h1>Calculator</h1>

      {status.kind === 'unavailable' && (
        <p role="alert" className="banner">
          Couldn't reach the server. Check your connection and try again.
        </p>
      )}

      <form onSubmit={handleSubmit}>
        <div role="radiogroup" aria-labelledby="operation-label" className="operations">
          <span id="operation-label" className="visually-hidden">
            Operation
          </span>
          {OPERATIONS.map(({ id, label }) => (
            <label key={id} className="operation">
              <input
                type="radio"
                name="operation"
                value={id}
                checked={id === operationId}
                onChange={() => handleOperationChange(id)}
              />
              <span>{label}</span>
            </label>
          ))}
        </div>

        <div className="expression">
          {renderInput('a')}
          <span className="symbol" aria-hidden="true">
            {operation.symbol}
          </span>
          {renderInput('b')}
          <button type="submit" aria-label="Calculate" disabled={status.kind === 'pending'}>
            =
          </button>
        </div>
      </form>

      <output className="result">
        {status.kind === 'result' && (
          <span className="result-value">{formatResult(status.value)}</span>
        )}
        {status.kind === 'rejected' && (
          <span className="result-error">
            {REJECTION_MESSAGES[status.code] ?? "The calculator couldn't complete that request."}
          </span>
        )}
      </output>
    </main>
  )
}
