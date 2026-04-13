import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'

import Profile from '@/pages/profile.vue'

// Mock global fetch
global.fetch = vi.fn()

describe('Profile.vue', () => {

  beforeEach(() => {
    vi.clearAllMocks()
  })

  // ---------------------------- //
  // RENDER TEST                  //
  // ---------------------------- //

  it('renders profile page correctly', () => {
    const wrapper = mount(Profile)

    expect(wrapper.text()).toContain('Profile')
    expect(wrapper.text()).toContain('Profile Completion')
  })

  // ---------------------------- //
  // INPUT BINDING TEST           //
  // ---------------------------- //

  it('updates form fields when user types', async () => {
    const wrapper = mount(Profile)

    const inputs = wrapper.findAll('input')

    await inputs[0].setValue('Testname')
    await inputs[1].setValue('nameTest')

    expect(wrapper.vm.form.first_name).toBe('Testname')
    expect(wrapper.vm.form.last_name).toBe('nameTest')
  })

  // ---------------------------- //
  // COMPLETION PERCENTAGE TEST   //
  // ---------------------------- //

  it('calculates completion percentage correctly', async () => {
    const wrapper = mount(Profile)

    wrapper.vm.form.first_name = 'John'
    wrapper.vm.form.last_name = 'Doe'

    await wrapper.vm.$nextTick()

    const totalFields = Object.keys(wrapper.vm.form).length
    const expected = Math.round((2 / totalFields) * 100)

    expect(wrapper.vm.completionPercentage).toBe(expected)
  })

  // ---------------------------- //
  // GET PROFILE (onMounted)      //
  // ---------------------------- //

  it('loads profile data from API on mount', async () => {
    const mockData = {
      first_name: 'John',
      last_name: 'Doe',
      phone: '1234567890',
      city: 'Newark',
      state: 'NJ',
      country: 'USA',
      linkedin_url: 'linkedin.com/test',
      portfolio_url: 'portfolio.com',
      summary: 'Test summary'
    }

    fetch.mockResolvedValueOnce({
      ok: true,
      json: async () => mockData
    })

    const wrapper = mount(Profile)

    await flushPromises()

    expect(wrapper.vm.form.first_name).toBe('John')
    expect(wrapper.vm.form.city).toBe('Newark')
    expect(fetch).toHaveBeenCalledWith('/api/profile', { method: 'GET' })
  })

  // ---------------------------- //
  // SAVE PROFILE (PUT)           //
  // ---------------------------- //

  it('sends updated profile data to API on save', async () => {
    fetch.mockResolvedValueOnce({ ok: true })

    const wrapper = mount(Profile)

    wrapper.vm.form.first_name = 'John'
    wrapper.vm.form.last_name = 'Doe'

    await wrapper.vm.$nextTick()

    await wrapper.find('button').trigger('click')

    expect(fetch).toHaveBeenCalledWith('/api/profile', {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        first_name: 'John',
        last_name: 'Doe',
        phone: '',
        city: '',
        state: '',
        country: '',
        linkedin_url: '',
        portfolio_url: '',
        summary: ''
      })
    })
  })

  // ---------------------------- //
  // EMPTY FORM COMPLETION        //
  // ---------------------------- //

  it('returns 0% completion when form is empty', () => {
    const wrapper = mount(Profile)

    expect(wrapper.vm.completionPercentage).toBe(0)
  })

  // ---------------------------- //
  // FULL FORM COMPLETION         //
  // ---------------------------- //

  it('returns 100% completion when all fields are filled', async () => {
    const wrapper = mount(Profile)

    Object.keys(wrapper.vm.form).forEach(key => {
      wrapper.vm.form[key] = 'test'
    })

    await wrapper.vm.$nextTick()

    expect(wrapper.vm.completionPercentage).toBe(100)
  })

  // ---------------------------- //
  // ERROR HANDLING (GET)         //
  // ---------------------------- //

  it('handles API error gracefully on GET', async () => {
    fetch.mockRejectedValueOnce(new Error('API failure'))

    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {})

    mount(Profile)

    await flushPromises()

    expect(consoleSpy).toHaveBeenCalled()
    consoleSpy.mockRestore()
  })

  // ---------------------------- //
  // ERROR HANDLING (PUT)         //
  // ---------------------------- //

  it('handles API error gracefully on save', async () => {
    fetch.mockRejectedValueOnce(new Error('API failure'))

    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {})

    const wrapper = mount(Profile)

    await wrapper.find('button').trigger('click')
    await flushPromises()

    expect(consoleSpy).toHaveBeenCalled()
    consoleSpy.mockRestore()
  })
})
