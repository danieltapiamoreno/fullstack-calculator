import { useRef, useState } from 'react'
import type { FormEvent } from 'react'
import { calculate } from './api'
import type { CalculationOutcome } from './api'
import { formatResult } from './format'
import { appendEntry, deleteLast, toggleSign } from './keypad'
import { OPERATIONS } from './operations'
import type { Field, Operation } from './operations'
import { parseOperand } from './validation'

type Status = { kind: 'idle' } | { kind: 'pending' } | CalculationOutcome

const IDLE: Status = { kind: 'idle' }

const ENTRY_KEYS = ['7', '8', '9', '4', '5', '6', '1', '2', '3', '0', '.']

/** Past this many characters, a result needs a smaller size to fit the display. */
const LONG_RESULT_LENGTH = 9

const FIELD_LABELS: Record<Field, string> = { a: 'First number', b: 'Second number' }

const REJECTION_MESSAGES: Record<string, string> = {
  division_by_zero: 'Cannot divide by zero.',
  undefined_result: "The result isn't a real number.",
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
  const [activeField, setActiveField] = useState<Field>('a')
  // Cancels the request in flight, so a response to outdated input is never shown.
  const inFlightRequest = useRef<AbortController | null>(null)

  function resetStatus() {
    inFlightRequest.current?.abort()
    setStatus(IDLE)
  }

  // Choosing an operation after a result continues from that result, as on a pocket calculator.
  function handleOperationChange(next: Operation) {
    const first = status.kind === 'success' ? formatResult(status.result) : values.a
    if (status.kind === 'success') setValues({ a: first, b: '' })
    setOperation(next)
    setActiveField(!next.oneInput && first !== '' ? 'b' : 'a')
    setFieldErrors({})
    resetStatus()
  }

  function handleValueChange(field: Field, value: string) {
    setValues((current) => ({ ...current, [field]: value }))
    setFieldErrors((current) => ({ ...current, [field]: undefined }))
    resetStatus()
  }

  const keypadField: Field = operation.oneInput ? 'a' : activeField

  // A digit after a result starts a new calculation, instead of extending the old input.
  function handleEntryKey(entry: string) {
    if (status.kind === 'success') {
      setValues({ a: appendEntry('', entry), b: '' })
      setActiveField('a')
      setFieldErrors({})
      resetStatus()
    } else {
      handleValueChange(keypadField, appendEntry(values[keypadField], entry))
    }
  }

  function handleClear() {
    setValues({ a: '', b: '' })
    setActiveField('a')
    setFieldErrors({})
    resetStatus()
  }

  async function handleSubmit(event: FormEvent) {
    event.preventDefault()
    if (status.kind === 'pending') return

    const a = parseOperand(values.a)
    const b = operation.oneInput ? null : parseOperand(values.b)
    if (!a.ok || b?.ok === false) {
      setFieldErrors({
        a: a.ok ? undefined : a.message,
        b: b?.ok === false ? b.message : undefined,
      })
      return
    }

    const controller = new AbortController()
    inFlightRequest.current = controller
    setFieldErrors({})
    setStatus({ kind: 'pending' })

    const operands = b === null ? { a: a.value } : { a: a.value, b: b.value }
    const outcome = await calculate(operation.id, operands, controller.signal)
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
      <div className={field === keypadField ? 'field field-active' : 'field'}>
        <label htmlFor={field} className="visually-hidden">
          {FIELD_LABELS[field]}
        </label>
        <input
          id={field}
          type="text"
          placeholder="0"
          autoComplete="off"
          autoCapitalize="off"
          autoCorrect="off"
          spellCheck={false}
          value={values[field]}
          aria-invalid={error !== undefined}
          aria-describedby={error === undefined ? undefined : errorId}
          onChange={(event) => handleValueChange(field, event.target.value)}
          onFocus={() => setActiveField(field)}
        />
        {error !== undefined && (
          <p id={errorId} className="field-error">
            {error}
          </p>
        )}
      </div>
    )
  }

  const symbol = (
    <span className="symbol" aria-hidden="true">
      {operation.symbol}
    </span>
  )

  const resultText = status.kind === 'success' ? formatResult(status.result) : ''

  return (
    <main className="calculator">
      <h1>Calculator</h1>

      <form onSubmit={handleSubmit}>
        <div className="display">
          {status.kind === 'unavailable' && (
            <p role="alert" className="banner">
              Couldn't reach the server. Check your connection and try again.
            </p>
          )}

          <div
            role="group"
            aria-label="Expression"
            className={operation.oneInput ? 'expression expression-one-input' : 'expression'}
          >
            {operation.oneInput && symbol}
            {renderInput('a')}
            {!operation.oneInput && symbol}
            {!operation.oneInput && renderInput('b')}
          </div>

          <output className="result">
            {status.kind === 'success' && (
              <span
                className={
                  resultText.length > LONG_RESULT_LENGTH
                    ? 'result-value result-value-long'
                    : 'result-value'
                }
              >
                {resultText}
              </span>
            )}
            {status.kind === 'rejected' && (
              <span className="result-error">{rejectionMessage(status.code)}</span>
            )}
          </output>
        </div>

        <div className="keypad">
          <div className="keys-entry">
            <button type="button" className="key key-clear" onClick={handleClear}>
              AC
            </button>
            <button
              type="button"
              className="key key-function"
              aria-label="Toggle sign"
              onClick={() => handleValueChange(keypadField, toggleSign(values[keypadField]))}
            >
              ±
            </button>
            <button
              type="button"
              className="key key-function"
              aria-label="Delete"
              onClick={() => handleValueChange(keypadField, deleteLast(values[keypadField]))}
            >
              ⌫
            </button>
            {ENTRY_KEYS.map((entry) => (
              <button
                key={entry}
                type="button"
                className={entry === '0' ? 'key key-wide' : 'key'}
                onClick={() => handleEntryKey(entry)}
              >
                {entry}
              </button>
            ))}
          </div>

          <div className="keys-operations">
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
                    // A radio that is already checked fires no change event, but pressing + after
                    // typing the first number with Add selected must still move on to the second.
                    onClick={() => {
                      if (option.id === operation.id) handleOperationChange(option)
                    }}
                  />
                  <span className="key" aria-hidden="true">
                    {option.keyLabel ?? option.symbol}
                  </span>
                  <span className="visually-hidden">{option.label}</span>
                </label>
              ))}
            </div>

            <button
              type="submit"
              className="key key-equals key-wide"
              aria-label="Calculate"
              disabled={status.kind === 'pending'}
            >
              =
            </button>
          </div>
        </div>
      </form>
    </main>
  )
}
