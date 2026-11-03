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
