import { useRef, useState } from 'react'
import type { FormEvent } from 'react'
import { calculate } from './api'
import type { CalculationOutcome } from './api'
import { formatResult } from './format'
import { OPERATIONS } from './operations'
import type { Field, Operation } from './operations'
import { parseOperand } from './validation'

type Status = { kind: 'idle' } | { kind: 'pending' } | CalculationOutcome

const IDLE: Status = { kind: 'idle' }

const FIELD_LABELS: Record<Field, string> = { a: 'First number', b: 'Second number' }

const REJECTION_MESSAGES: Record<string, string> = {
  division_by_zero: 'Cannot divide by zero.',
  result_out_of_range: 'The result is too large to represent.',
}

function rejectionMessage(code: string) {
  return REJECTION_MESSAGES[code] ?? "The calculator couldn't complete that request."
}

export function Calculator() {
  const [operation, setOperation] = useState(OPERATIONS[0])
  const [values, setValues] = useState<Record<Field, string>>({ a: '', b: '' })
  const [fieldErrors, setFieldErrors] = useState<Partial<Record<Field, string>>>({})
  const [status, setStatus] = useState<Status>(IDLE)
  // Cancels the request in flight, so a response to outdated input is never shown.
  const inFlightRequest = useRef<AbortController | null>(null)

  function resetStatus() {
    inFlightRequest.current?.abort()
    setStatus(IDLE)
  }

  function handleOperationChange(next: Operation) {
    setOperation(next)
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

    const controller = new AbortController()
    inFlightRequest.current = controller
    setFieldErrors({})
    setStatus({ kind: 'pending' })

    const outcome = await calculate(operation.id, { a: a.value, b: b.value }, controller.signal)
    if (controller.signal.aborted) return

    const errorField = outcome.kind === 'rejected' ? operation.errorFields?.[outcome.code] : undefined
    if (outcome.kind === 'rejected' && errorField !== undefined) {
      setFieldErrors({ [errorField]: rejectionMessage(outcome.code) })
      setStatus(IDLE)
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
          {OPERATIONS.map((option) => (
            <label key={option.id} className="operation">
              <input
                type="radio"
                name="operation"
                value={option.id}
                checked={option.id === operation.id}
                onChange={() => handleOperationChange(option)}
              />
              <span>{option.label}</span>
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
        {status.kind === 'success' && (
          <span className="result-value">{formatResult(status.result)}</span>
        )}
        {status.kind === 'rejected' && (
          <span className="result-error">{rejectionMessage(status.code)}</span>
        )}
      </output>
    </main>
  )
}
