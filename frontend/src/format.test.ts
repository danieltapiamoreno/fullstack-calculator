import { describe, expect, it } from 'vitest'
import { formatResult } from './format'

describe('formatResult', () => {
  it.each([
    [5, '5'],
    [-1.5, '-1.5'],
    [0.30000000000000004, '0.3'],
    [1.23456789012345, '1.23456789012'],
    [123456789012345, '123456789012000'],
    [1e21, '1e+21'],
    [1e-7, '1e-7'],
    [0, '0'],
  ])('displays %d as %s', (value, expected) => {
    expect(formatResult(value)).toBe(expected)
  })
})
