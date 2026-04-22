import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import DocumentLibrary from '@/pages/library.vue'

// Mock fetch
global.fetch = vi.fn()

describe('DocumentLibrary.vue', () => {

  beforeEach(() => {
    vi.clearAllMocks()
  })

  // ---------------------------- //
  // RENDER TEST                  //
  // ---------------------------- //

  it('renders document library page', async () => {
    fetch
      .mockResolvedValueOnce({ ok: true, json: async () => [] }) // documents
      .mockResolvedValueOnce({ ok: true, json: async () => [] }) // jobs

    const wrapper = mount(DocumentLibrary)
    await flushPromises()

    expect(wrapper.text()).toContain('Document Library')
    expect(wrapper.text()).toContain('Upload Documents')
  })

  // ---------------------------- //
  // FETCH DOCUMENTS (onMounted)  //
  // ---------------------------- //

  it('loads documents from API on mount', async () => {
    const mockDocs = [{ id: 1, title: 'Resume', type: 'PDF' }]

    fetch
      .mockResolvedValueOnce({
        ok: true,
        json: async () => mockDocs
      }) // documents
      .mockResolvedValueOnce({
        ok: true,
        json: async () => []
      }) // jobs

    const wrapper = mount(DocumentLibrary)
    await flushPromises()

    expect(wrapper.vm.documents.length).toBe(1)
    expect(wrapper.vm.documents[0].title).toBe('Resume')
  })

  // ---------------------------- //
  // UPLOAD WITHOUT FILE          //
  // ---------------------------- //

  it('shows error if upload attempted without file', async () => {
    fetch
      .mockResolvedValueOnce({ ok: true, json: async () => [] }) // documents
      .mockResolvedValueOnce({ ok: true, json: async () => [] }) // jobs

    const wrapper = mount(DocumentLibrary)
    await flushPromises()

    const uploadBtn = wrapper.findAll('button')
      .find(btn => btn.text().includes('Upload File'))

    await uploadBtn.trigger('click')

    expect(wrapper.vm.error).toBe('Please select a file')
  })

  // ---------------------------- //
  // SUCCESSFUL UPLOAD            //
  // ---------------------------- //

  it('uploads file and updates document list', async () => {
    fetch
      .mockResolvedValueOnce({ ok: true, json: async () => [] }) // documents
      .mockResolvedValueOnce({ ok: true, json: async () => [] }) // jobs
      .mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          id: 1,
          name: 'Resume.pdf',
          document_type: 'PDF',
          url: 'test-url'
        })
      })

    const wrapper = mount(DocumentLibrary)
    await flushPromises()

    const file = new File(['test'], 'Resume.pdf', { type: 'application/pdf' })

    wrapper.vm.selectedFile = file
    wrapper.vm.selectedJobId = 1 // REQUIRED now

    await wrapper.vm.uploadFile()
    await flushPromises()

    expect(wrapper.vm.documents.length).toBe(1)
    expect(wrapper.vm.uploadMessage).toBe('Upload successful!')
  })

  // ---------------------------- //
  // FAILED UPLOAD (API ERROR)    //
  // ---------------------------- //

  it('handles failed upload response', async () => {
    fetch
      .mockResolvedValueOnce({ ok: true, json: async () => [] }) // documents
      .mockResolvedValueOnce({ ok: true, json: async () => [] }) // jobs
      .mockResolvedValueOnce({ ok: false })

    const wrapper = mount(DocumentLibrary)
    await flushPromises()

    wrapper.vm.selectedFile = new File(['test'], 'file.pdf')
    wrapper.vm.selectedJobId = 1

    await wrapper.vm.uploadFile()
    await flushPromises()

    expect(wrapper.vm.error).toBe('Only PDF files allowed')
  })

  // ---------------------------- //
  // DELETE DOCUMENT              //
  // ---------------------------- //

  it('deletes document after confirmation', async () => {
    fetch
      .mockResolvedValueOnce({ ok: true, json: async () => [] }) // documents
      .mockResolvedValueOnce({ ok: true, json: async () => [] }) // jobs
      .mockResolvedValueOnce({ ok: true }) // delete

    const confirmSpy = vi.spyOn(window, 'confirm').mockReturnValue(true)

    const wrapper = mount(DocumentLibrary)
    await flushPromises()

    wrapper.vm.documents = [{ id: 1, title: 'Test Doc' }]

    await wrapper.vm.deleteDocument(1)
    await flushPromises()

    expect(wrapper.vm.documents.length).toBe(0)

    confirmSpy.mockRestore()
  })

  // ---------------------------- //
  // DELETE CANCELLED             //
  // ---------------------------- //

  it('does not delete document if user cancels', async () => {
    fetch
      .mockResolvedValueOnce({ ok: true, json: async () => [] })
      .mockResolvedValueOnce({ ok: true, json: async () => [] })

    const confirmSpy = vi.spyOn(window, 'confirm').mockReturnValue(false)

    const wrapper = mount(DocumentLibrary)
    await flushPromises()

    wrapper.vm.documents = [{ id: 1, title: 'Test Doc' }]

    await wrapper.vm.deleteDocument(1)

    expect(wrapper.vm.documents.length).toBe(1)

    confirmSpy.mockRestore()
  })

})
