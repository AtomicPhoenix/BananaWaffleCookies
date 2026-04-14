import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Dashboard from '@/pages/dashboard.vue'

// mock fetch
global.fetch = vi.fn()


describe('Dashboard.vue', () => {

  beforeEach(() => {
    vi.clearAllMocks()
  })

  // ---------------------------- //
  // TEST 1: RENDER + DEFAULT UI  //
  // ---------------------------- //

  it('renders dashboard title and default welcome message', async () => {
    fetch.mockResolvedValue({ ok: true, json: async () => [] })

    const wrapper = mount(Dashboard)

    await flushPromises()

    expect(wrapper.text()).toContain('My Job Dashboard')
    expect(wrapper.text()).toContain('Welcome back')
  })

  // ---------------------------- //
  // TEST 2: SEARCH FUNCTION      //
  // ---------------------------- //

  it('performs search and updates searchResults', async () => {
    const mockResults = [
      { id: 1, title: 'Software Engineer', company_name: 'Google' }
    ]

    fetch
      .mockResolvedValueOnce({ ok: true, json: async () => [] }) // fetchUserJobs
      .mockResolvedValueOnce({ ok: true, json: async () => ({}) }) // fetchProfile
      .mockResolvedValueOnce({ ok: true, json: async () => mockResults }) // search

    const wrapper = mount(Dashboard)

    await flushPromises()

    // simulate typing
    await wrapper.find('#job-search').setValue('engineer')

    // submit form
    await wrapper.find('form').trigger('submit')

    await flushPromises()

    expect(wrapper.vm.searchResults.length).toBe(1)
    expect(wrapper.vm.searchResults[0].title).toBe('Software Engineer')
  })

  // ---------------------------- //
  // TEST 3: FILTER OUT ARCHIVED  //
  // ---------------------------- //

  it('filters out archived jobs when showArchived is false', async () => {
    const jobs = [
      { id: 1, status: 'applied', is_archived: false },
      { id: 2, status: 'applied', is_archived: true }
    ]

    fetch
      .mockResolvedValueOnce({ ok: true, json: async () => jobs })
      .mockResolvedValueOnce({ ok: true, json: async () => ({}) })

    const wrapper = mount(Dashboard)

    await flushPromises()

    // default: showArchived = false
    const displayed = wrapper.vm.displayedJobs

    expect(displayed.length).toBe(1)
    expect(displayed[0].id).toBe(1)
  })

  // ---------------------------- //
  // TEST 4: SORTING LOGIC        //
  // ---------------------------- //

  it('sorts jobs by title ascending', async () => {
    const jobs = [
      { id: 1, title: 'Zebra Job', status: 'applied', is_archived: false },
      { id: 2, title: 'Alpha Job', status: 'applied', is_archived: false }
    ]

    fetch
      .mockResolvedValueOnce({ ok: true, json: async () => jobs })
      .mockResolvedValueOnce({ ok: true, json: async () => ({}) })

    const wrapper = mount(Dashboard)

    await flushPromises()

    // set sorting
    wrapper.vm.sortBy = 'title-asc'
    await wrapper.vm.$nextTick()

    const displayed = wrapper.vm.displayedJobs

    expect(displayed[0].title).toBe('Alpha Job')
    expect(displayed[1].title).toBe('Zebra Job')
  })

})
