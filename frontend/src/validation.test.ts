import { describe, expect, it } from 'vitest'
import { parseOperand } from './validation'

describe('parseOperand', () => {
  it.each([
    ['12', 12],
    ['-3.5', -3.5],
    ['.5', 0.5],
    ['5.', 5],
    ['0', 0],
    ['1e5', 100000],
    ['1.5E-3', 0.0015],
    ['  7  ', 7],
  ])('accepts %j', (input, expected) => {
    expect(parseOperand(input)).toEqual({ ok: true, value: expected })
  })

  it.each(['', '   '])('asks for a number when the input is %j', (input) => {
    expect(parseOperand(input)).toEqual({ ok: false, message: 'Enter a number.' })
  })

  it.each(['abc', '1,5', '1.2.3', '--1', '-', '.', 'e5', '1e', '0x10', 'Infinity', 'NaN'])(
    'rejects %j as malformed',
    (input) => {
      expect(parseOperand(input)).toEqual({
        ok: false,
        message: 'Enter a valid number, like 12, -3.5, or 1e5.',
      })
    },
  )

  it.each(['1e999', '-1e999'])('rejects %j as too large to represent', (input) => {
    expect(parseOperand(input)).toEqual({ ok: false, message: 'This number is too large.' })
  })
})
