import { render, screen, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'
import { Calculator } from './Calculator'

function jsonResponse(status: number, body: unknown) {
  return new Response(JSON.stringify(body), { status })
}

function stubFetch(...responses: Response[]) {
  const fetchStub = vi.fn<typeof fetch>()
  responses.forEach((response) => fetchStub.mockResolvedValueOnce(response))
  vi.stubGlobal('fetch', fetchStub)
  return fetchStub
}

function stubFailedFetch() {
  vi.stubGlobal('fetch', vi.fn<typeof fetch>().mockRejectedValue(new TypeError('Failed to fetch')))
}

function setup() {
  render(<Calculator />)
  const user = userEvent.setup()
  const first = screen.getByRole('textbox', { name: 'First number' })
  const second = screen.getByRole('textbox', { name: 'Second number' })
  const submit = screen.getByRole('button', { name: 'Calculate' })
  return {
    user,
    first,
    second,
    submit,
    result: screen.getByRole('status'),
    async calculate(a: string, b: string) {
      await user.type(first, a)
      await user.type(second, b)
      await user.click(submit)
    },
  }
}

describe('Calculator', () => {
  it('sends the operands to the API and shows the rounded result', async () => {
    const fetchStub = stubFetch(jsonResponse(200, { result: 0.30000000000000004 }))
    const { calculate, result } = setup()

    await calculate('0.1', '0.2')

    expect(await within(result).findByText('0.3')).toBeInTheDocument()
    expect(fetchStub).toHaveBeenCalledOnce()
    const [url, init] = fetchStub.mock.calls[0]
    expect(url).toBe('/api/v1/add')
    expect(init?.method).toBe('POST')
    expect(JSON.parse(init?.body as string)).toEqual({ a: 0.1, b: 0.2 })
  })

  it('submits when the user presses Enter', async () => {
    stubFetch(jsonResponse(200, { result: 5 }))
    const { user, first, second, result } = setup()

    await user.type(first, '2')
    await user.type(second, '3{Enter}')

    expect(await within(result).findByText('5')).toBeInTheDocument()
  })

  it('accepts negative numbers and scientific notation', async () => {
    const fetchStub = stubFetch(jsonResponse(200, { result: 99999 }))
    const { calculate } = setup()

    await calculate('-1', '1e5')

    expect(JSON.parse(fetchStub.mock.calls[0][1]?.body as string)).toEqual({ a: -1, b: 100000 })
  })

  it('offers the operations as a radio group with add selected', () => {
    setup()

    const group = screen.getByRole('radiogroup', { name: 'Operation' })
    expect(within(group).getByRole('radio', { name: 'Add' })).toBeChecked()
  })

  it.each([
    { label: 'Subtract', symbol: '−', route: '/api/v1/subtract', answer: -1, shown: '-1' },
    { label: 'Multiply', symbol: '×', route: '/api/v1/multiply', answer: 6, shown: '6' },
  ])(
    'calculates with $label and shows its symbol in the form',
    async ({ label, symbol, route, answer, shown }) => {
      const fetchStub = stubFetch(jsonResponse(200, { result: answer }))
      const { calculate, user, result } = setup()

      await user.click(screen.getByRole('radio', { name: label }))
      expect(screen.getByText(symbol)).toBeInTheDocument()
      expect(screen.queryByText('+')).not.toBeInTheDocument()
      await calculate('2', '3')

      expect(await within(result).findByText(shown)).toBeInTheDocument()
      expect(fetchStub.mock.calls[0][0]).toBe(route)
      expect(JSON.parse(fetchStub.mock.calls[0][1]?.body as string)).toEqual({ a: 2, b: 3 })
    },
  )

  describe('when the operation changes', () => {
    async function switchOperation(user: ReturnType<typeof userEvent.setup>) {
      await user.click(screen.getByRole('radio', { name: 'Subtract' }))
    }

    it('keeps the entered values', async () => {
      const { user, first, second } = setup()

      await user.type(first, '2')
      await user.type(second, '3')
      await switchOperation(user)

      expect(first).toHaveValue('2')
      expect(second).toHaveValue('3')
    })

    it('clears the result', async () => {
      stubFetch(jsonResponse(200, { result: 5 }))
      const { calculate, user, result } = setup()

      await calculate('2', '3')
      await within(result).findByText('5')
      await switchOperation(user)

      expect(result).toBeEmptyDOMElement()
    })

    it('clears a form-level error', async () => {
      stubFetch(jsonResponse(422, { error: { code: 'result_out_of_range', message: 'too large' } }))
      const { calculate, user, result } = setup()

      await calculate('1e308', '1e308')
      await within(result).findByText('The result is too large to represent.')
      await switchOperation(user)

      expect(result).toBeEmptyDOMElement()
    })

    it('clears the field errors', async () => {
      stubFetch()
      const { user, first, second, submit } = setup()

      await user.type(second, 'abc')
      await user.click(submit)
      expect(first).toBeInvalid()
      expect(second).toBeInvalid()
      await switchOperation(user)

      expect(first).toBeValid()
      expect(first).not.toHaveAccessibleDescription()
      expect(second).toBeValid()
      expect(second).not.toHaveAccessibleDescription()
      expect(second).toHaveValue('abc')
    })

    it('clears the network failure banner', async () => {
      stubFailedFetch()
      const { calculate, user } = setup()

      await calculate('2', '3')
      await screen.findByRole('alert')
      await switchOperation(user)

      expect(screen.queryByRole('alert')).not.toBeInTheDocument()
    })
  })

  it('uses text inputs whose keyboard includes a minus sign on iOS', () => {
    const { first, second } = setup()

    for (const input of [first, second]) {
      expect(input).toHaveAttribute('type', 'text')
      expect(input).not.toHaveAttribute('inputmode', 'decimal')
    }
  })

  it('shows an error under each invalid input without calling the API', async () => {
    const fetchStub = stubFetch()
    const { user, first, second, submit } = setup()

    await user.type(second, 'abc')
    await user.click(submit)

    expect(first).toBeInvalid()
    expect(first).toHaveAccessibleDescription('Enter a number.')
    expect(second).toBeInvalid()
    expect(second).toHaveAccessibleDescription('Enter a valid number, like 12, -3.5, or 1e5.')
    expect(fetchStub).not.toHaveBeenCalled()
  })

  it('clears only the error of the input being changed', async () => {
    stubFetch()
    const { user, first, second, submit } = setup()

    await user.click(submit)
    await user.type(first, '1')

    expect(first).toBeValid()
    expect(first).not.toHaveAccessibleDescription()
    expect(second).toBeInvalid()
  })

  it('clears the result when an input changes', async () => {
    stubFetch(jsonResponse(200, { result: 5 }))
    const { calculate, user, first, result } = setup()

    await calculate('2', '3')
    await within(result).findByText('5')
    await user.type(first, '1')

    expect(result).toBeEmptyDOMElement()
  })

  describe('while a request is in flight', () => {
    function stubPendingFetch() {
      let respond!: (response: Response) => void
      const fetchStub = vi.fn<typeof fetch>(
        () => new Promise((resolve) => { respond = resolve }),
      )
      vi.stubGlobal('fetch', fetchStub)
      return { fetchStub, respond: (response: Response) => respond(response) }
    }

    it('disables the button and ignores another Enter press', async () => {
      const { fetchStub, respond } = stubPendingFetch()
      const { calculate, user, second, submit, result } = setup()

      await calculate('2', '3')
      expect(submit).toBeDisabled()
      await user.type(second, '{Enter}')
      expect(fetchStub).toHaveBeenCalledOnce()

      respond(jsonResponse(200, { result: 5 }))

      expect(await within(result).findByText('5')).toBeInTheDocument()
      expect(submit).toBeEnabled()
    })

    it('cancels the request and ignores its response if an input changes', async () => {
      const { fetchStub, respond } = stubPendingFetch()
      const { calculate, user, first, submit, result } = setup()

      await calculate('2', '3')
      await user.type(first, '1')
      expect(fetchStub.mock.calls[0][1]?.signal?.aborted).toBe(true)
      expect(submit).toBeEnabled()

      respond(jsonResponse(200, { result: 5 }))
      await new Promise((resolve) => setTimeout(resolve))

      expect(result).toBeEmptyDOMElement()
    })

    it('cancels the request and ignores its response if the operation changes', async () => {
      const { fetchStub, respond } = stubPendingFetch()
      const { calculate, user, first, second, submit, result } = setup()

      await calculate('2', '3')
      await user.click(screen.getByRole('radio', { name: 'Multiply' }))
      expect(fetchStub.mock.calls[0][1]?.signal?.aborted).toBe(true)
      expect(submit).toBeEnabled()

      respond(jsonResponse(200, { result: 5 }))
      await new Promise((resolve) => setTimeout(resolve))

      expect(result).toBeEmptyDOMElement()
      expect(first).toHaveValue('2')
      expect(second).toHaveValue('3')
    })

    it('shows no banner when switching operation aborts the request', async () => {
      const fetchStub = vi.fn<typeof fetch>(
        (_url, init) =>
          new Promise((_resolve, reject) => {
            init?.signal?.addEventListener('abort', () =>
              reject(new DOMException('Aborted', 'AbortError')),
            )
          }),
      )
      vi.stubGlobal('fetch', fetchStub)
      const { calculate, user } = setup()

      await calculate('2', '3')
      await user.click(screen.getByRole('radio', { name: 'Multiply' }))
      await new Promise((resolve) => setTimeout(resolve))

      expect(screen.queryByRole('alert')).not.toBeInTheDocument()
    })

    it('shows only the new answer when an old response arrives after switching and resubmitting', async () => {
      const responders: ((response: Response) => void)[] = []
      const fetchStub = vi.fn<typeof fetch>(
        () => new Promise((resolve) => responders.push(resolve)),
      )
      vi.stubGlobal('fetch', fetchStub)
      const { calculate, user, submit, result } = setup()

      await calculate('2', '3')
      await user.click(screen.getByRole('radio', { name: 'Multiply' }))
      await user.click(submit)
      expect(fetchStub.mock.calls[1][0]).toBe('/api/v1/multiply')

      responders[1](jsonResponse(200, { result: 6 }))
      await within(result).findByText('6')
      responders[0](jsonResponse(200, { result: 5 }))
      await new Promise((resolve) => setTimeout(resolve))

      expect(result).toHaveTextContent('6')
    })

    it('ignores a division_by_zero response to a cancelled request', async () => {
      const { respond } = stubPendingFetch()
      const { calculate, user, first, second, result } = setup()

      await user.click(screen.getByRole('radio', { name: 'Divide' }))
      await calculate('1', '0')
      await user.type(first, '2')

      respond(jsonResponse(422, { error: { code: 'division_by_zero', message: 'server text' } }))
      await new Promise((resolve) => setTimeout(resolve))

      expect(second).toBeValid()
      expect(result).toBeEmptyDOMElement()
    })
  })

  describe('dividing', () => {
    const divisionByZero = () =>
      jsonResponse(422, { error: { code: 'division_by_zero', message: 'server text' } })

    async function divide(response: Response, a: string, b: string) {
      const fetchStub = stubFetch(response)
      const context = setup()
      await context.user.click(screen.getByRole('radio', { name: 'Divide' }))
      await context.calculate(a, b)
      return { ...context, fetchStub }
    }

    it('shows the divide symbol when divide is chosen', async () => {
      const { user } = setup()

      await user.click(screen.getByRole('radio', { name: 'Divide' }))

      expect(screen.getByText('÷')).toBeInTheDocument()
      expect(screen.queryByText('+')).not.toBeInTheDocument()
    })

    it('sends the operands to the divide route and shows the result', async () => {
      const { fetchStub, result } = await divide(jsonResponse(200, { result: 2 }), '6', '3')

      expect(await within(result).findByText('2')).toBeInTheDocument()
      const [url, init] = fetchStub.mock.calls[0]
      expect(url).toBe('/api/v1/divide')
      expect(JSON.parse(init?.body as string)).toEqual({ a: 6, b: 3 })
    })

    it('shows division_by_zero under the second input, not the result area', async () => {
      const { first, second, submit, result } = await divide(divisionByZero(), '1', '0')

      expect(await screen.findByText('Cannot divide by zero.')).toBeInTheDocument()
      expect(second).toBeInvalid()
      expect(second).toHaveAccessibleDescription('Cannot divide by zero.')
      expect(first).toBeValid()
      expect(result).toBeEmptyDOMElement()
      expect(submit).toBeEnabled()
    })

    it('clears the division_by_zero error when the second input changes', async () => {
      const { user, second } = await divide(divisionByZero(), '1', '0')
      await screen.findByText('Cannot divide by zero.')
      expect(second).toBeInvalid()

      await user.type(second, '5')

      expect(second).toBeValid()
      expect(second).not.toHaveAccessibleDescription()
    })

    it('clears the division_by_zero error when the operation changes', async () => {
      const { user, second } = await divide(divisionByZero(), '1', '0')
      await screen.findByText('Cannot divide by zero.')
      expect(second).toBeInvalid()

      await user.click(screen.getByRole('radio', { name: 'Add' }))

      expect(second).toBeValid()
      expect(second).toHaveValue('0')
    })

    it('shows result_out_of_range as a form-level message', async () => {
      const { first, second, result } = await divide(
        jsonResponse(422, { error: { code: 'result_out_of_range', message: 'server text' } }),
        '1e308',
        '1e-308',
      )

      expect(
        await within(result).findByText('The result is too large to represent.'),
      ).toBeInTheDocument()
      expect(first).toBeValid()
      expect(second).toBeValid()
    })
  })

  describe('when the server reports an error', () => {
    async function submitValidInput(response: Response) {
      stubFetch(response)
      const context = setup()
      await context.calculate('1e308', '1e308')
      return context
    }

    it('shows result_out_of_range as a form-level message in the result area', async () => {
      const { first, second, result } = await submitValidInput(
        jsonResponse(422, { error: { code: 'result_out_of_range', message: 'too large' } }),
      )

      expect(
        await within(result).findByText('The result is too large to represent.'),
      ).toBeInTheDocument()
      expect(first).toBeValid()
      expect(second).toBeValid()
    })

    it('keeps a code at form level for an operation that does not map it to a field', async () => {
      const { first, second, result } = await submitValidInput(
        jsonResponse(422, { error: { code: 'division_by_zero', message: 'server text' } }),
      )

      expect(await within(result).findByText('Cannot divide by zero.')).toBeInTheDocument()
      expect(first).toBeValid()
      expect(second).toBeValid()
    })

    it('shows a general message for an unrecognised error code', async () => {
      const { result } = await submitValidInput(
        jsonResponse(400, { error: { code: 'something_new', message: 'unknown' } }),
      )

      expect(
        await within(result).findByText("The calculator couldn't complete that request."),
      ).toBeInTheDocument()
    })
  })

  describe('when the server cannot be reached', () => {
    const banner = "Couldn't reach the server. Check your connection and try again."

    it('shows a banner when the request fails', async () => {
      stubFailedFetch()
      const { calculate, submit, result } = setup()

      await calculate('2', '3')

      expect(await screen.findByRole('alert')).toHaveTextContent(banner)
      expect(result).toBeEmptyDOMElement()
      expect(submit).toBeEnabled()
    })

    it('shows a banner when the response is not the JSON the API returns', async () => {
      stubFetch(new Response('Bad Gateway', { status: 502 }))
      const { calculate } = setup()

      await calculate('2', '3')

      expect(await screen.findByRole('alert')).toHaveTextContent(banner)
    })

    it('clears the banner when an input changes', async () => {
      stubFailedFetch()
      const { calculate, user, first } = setup()

      await calculate('2', '3')
      await screen.findByRole('alert')
      await user.type(first, '1')

      expect(screen.queryByRole('alert')).not.toBeInTheDocument()
    })
  })
})
