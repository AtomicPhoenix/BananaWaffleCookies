import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Settings from '@/pages/settings.vue'

// Mock global fetch from backend information
global.fetch = vi.fn()

describe('Settings Page', () => {

  beforeEach(() => {
    vi.clearAllMocks()
  })

  // ========================= //
  // Load Existing Settings    //
  // ========================= //

  it('loads saved email on mount', async () => {
    fetch.mockResolvedValueOnce({
      ok: true,
      json: async () => ({ Email: 'test@example.com' })
    })

    const wrapper = mount(Settings)
    await flushPromises()
    const emailInput = wrapper.find('input[type="email"]')

    // test 1
    expect(emailInput.element.value).toBe('test@example.com')
  })

  // ========================= //
  // Add Data + Save           //
  // ========================= //

  it('allows user to update email and persists after save + reload', async () => {
    // Initial load
    fetch.mockResolvedValueOnce({
      ok: true,
      json: async () => ({ Email: 'old@example.com' })
    })

    const wrapper = mount(Settings)
    await flushPromises()

    const emailInput = wrapper.find('input[type="email"]')
    await emailInput.setValue('new@example.com')

    // Mock save
    fetch.mockResolvedValueOnce({ ok: true })
    await wrapper.find('button').trigger('click')

    //test 2
    expect(fetch).toHaveBeenCalledWith('/api/settings', expect.objectContaining({
      method: 'PUT'
    }))

    // Simulate refresh (new mount)
    fetch.mockResolvedValueOnce({
      ok: true,
      json: async () => ({ Email: 'new@example.com' })
    })

    const wrapper2 = mount(Settings)
    await flushPromises()

    const newEmailInput = wrapper2.find('input[type="email"]')
    //test 3
    expect(newEmailInput.element.value).toBe('new@example.com')
  })

  // ========================= //
  // Remove Data + Save        //
  // ========================= //

  it('allows user to remove email and persists after reload', async () => {
    fetch.mockResolvedValueOnce({
      ok: true,
      json: async () => ({ Email: 'remove@example.com' })
    })

    const wrapper = mount(Settings)
    await flushPromises()

    const emailInput = wrapper.find('input[type="email"]')

    await emailInput.setValue('')

    fetch.mockResolvedValueOnce({ ok: true })

    await wrapper.find('button').trigger('click')

    // Simulate reload
    fetch.mockResolvedValueOnce({
      ok: true,
      json: async () => ({ Email: '' })
    })

    const wrapper2 = mount(Settings)
    await flushPromises()

    //test4
    expect(wrapper2.find('input[type="email"]').element.value).toBe('')
  })

  // =========================
  // Password Reset
  // =========================

  it('successfully submits password reset when passwords match', async () => {
    // 1. Initial GET
    fetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({ Email: 'test@example.com' })
    })

    const wrapper = mount(Settings)
    await flushPromises()

    const passwordInputs = wrapper.findAll('input[type="password"]')

    await passwordInputs[0].setValue('newpassword')
    await passwordInputs[1].setValue('newpassword')

    // 2. PUT request
    fetch.mockResolvedValueOnce({ ok: true })

    await wrapper.find('button').trigger('click')

    // wait for change
    await flushPromises()
    //test 6
    expect(wrapper.text()).toContain('Settings saved successfully!')
  })


  it('shows error when passwords do not match', async () => {
    fetch.mockResolvedValueOnce({
      ok: true,
      json: async () => ({ Email: 'test@example.com' })
    })

    const wrapper = mount(Settings)
    await flushPromises()

    const passwordInputs = wrapper.findAll('input[type="password"]')

    await passwordInputs[0].setValue('password1')
    await passwordInputs[1].setValue('password2')

    await wrapper.find('button').trigger('click')

    // test 7
    expect(wrapper.text()).toContain('Passwords do not match')
  })

  // =========================
  // Notifications Checkbox
  // =========================

  it('toggles notification checkbox correctly', async () => {
    fetch.mockResolvedValueOnce({
      ok: true,
      json: async () => ({ Email: 'test@example.com' })
    })

    const wrapper = mount(Settings)
    await flushPromises()
    const checkbox = wrapper.find('input[type="checkbox"]')

    // test 8
    expect(checkbox.element.checked).toBe(false)

    await checkbox.setValue(true)
    // test 9
    expect(checkbox.element.checked).toBe(true)

    await checkbox.setValue(false)
    // test 10
    expect(checkbox.element.checked).toBe(false)
  })

})
