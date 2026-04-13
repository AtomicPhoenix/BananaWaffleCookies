import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import DocumentLibrary from '@/pages/library.vue'
import Chatbox from '@/pages/chatbox.vue'

// Mock fetch
global.fetch = vi.fn()
describe('DocumentLibrary.vue', () => {

  beforeEach(() => {
    vi.clearAllMocks()
  })

  // ---------------------------- //
  // RENDER TEST                  //
  // ---------------------------- //

  it('renders document library page', () => {
    fetch.mockResolvedValue({ ok: true, json: async () => [] })

    const wrapper = mount(DocumentLibrary)

    expect(wrapper.text()).toContain('Document Library')
    expect(wrapper.text()).toContain('Upload Documents')
  })

  // ---------------------------- //
  // FETCH DOCUMENTS (onMounted)  //
  // ---------------------------- //

  it('loads documents from API on mount', async () => {
    const mockDocs = [
      { id: 1, title: 'Resume', type: 'PDF' }
    ]

    fetch.mockResolvedValueOnce({
      ok: true,
      json: async () => mockDocs
    })

    const wrapper = mount(DocumentLibrary)

    await flushPromises()

    expect(wrapper.vm.documents.length).toBe(1)
    expect(wrapper.vm.documents[0].title).toBe('Resume')
  })

  // ---------------------------- //
  // UPLOAD WITHOUT FILE          //
  // ---------------------------- //

  it('shows error if upload attempted without file', async () => {
    fetch.mockResolvedValue({ ok: true, json: async () => [] })

    const wrapper = mount(DocumentLibrary)

    await wrapper.find('button').trigger('click')

    expect(wrapper.vm.error).toBe('Please select a file')
  })

  // ---------------------------- //
  // SUCCESSFUL UPLOAD            //
  // ---------------------------- //

  it('uploads file and updates document list', async () => {
    fetch
      .mockResolvedValueOnce({ ok: true, json: async () => [] }) // initial fetch
      .mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          id: 1,
          name: 'Resume.pdf',
          type: 'PDF',
          url: 'test-url'
        })
      })

    const wrapper = mount(DocumentLibrary)

    const file = new File(['test'], 'Resume.pdf', { type: 'application/pdf' })

    wrapper.vm.selectedFile = file

    await wrapper.vm.uploadFile()
    await flushPromises()

    expect(wrapper.vm.documents.length).toBe(1)
    expect(wrapper.vm.uploadMessage).toBe('File uploaded successfully!')
  })

  // ---------------------------- //
  // FAILED UPLOAD (API ERROR)    //
  // ---------------------------- //

  it('handles failed upload response', async () => {
    fetch
      .mockResolvedValueOnce({ ok: true, json: async () => [] }) // initial fetch
      .mockResolvedValueOnce({ ok: false })

    const wrapper = mount(DocumentLibrary)

    wrapper.vm.selectedFile = new File(['test'], 'file.pdf')

    await wrapper.vm.uploadFile()
    await flushPromises()

    expect(wrapper.vm.error).toBe('Upload failed')
  })

  // ---------------------------- //
  // SERVER ERROR (UPLOAD)        //
  // ---------------------------- //

  it('handles server error during upload', async () => {
    fetch
      .mockResolvedValueOnce({ ok: true, json: async () => [] })
      .mockRejectedValueOnce(new Error('Server error'))

    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {})

    const wrapper = mount(DocumentLibrary)

    wrapper.vm.selectedFile = new File(['test'], 'file.pdf')

    await wrapper.vm.uploadFile()
    await flushPromises()

    expect(wrapper.vm.error).toBe('Server error')
    expect(consoleSpy).toHaveBeenCalled()

    consoleSpy.mockRestore()
  })

  // ---------------------------- //
  // OPEN DOCUMENT                //
  // ---------------------------- //

  it('opens document in new tab if URL exists', () => {
    fetch.mockResolvedValue({ ok: true, json: async () => [] })

    const wrapper = mount(DocumentLibrary)

    const openSpy = vi.spyOn(window, 'open').mockImplementation(() => {})

    wrapper.vm.openDocument({ url: 'http://test.com' })

    expect(openSpy).toHaveBeenCalledWith('http://test.com', '_blank')

    openSpy.mockRestore()
  })

  // ---------------------------- //
  // DELETE DOCUMENT              //
  // ---------------------------- //

  it('deletes document after confirmation', async () => {
    fetch
      .mockResolvedValueOnce({ ok: true, json: async () => [] }) // initial fetch
      .mockResolvedValueOnce({ ok: true }) // delete

    const confirmSpy = vi.spyOn(window, 'confirm').mockReturnValue(true)

    const wrapper = mount(DocumentLibrary)

    wrapper.vm.documents = [
      { id: 1, title: 'Test Doc' }
    ]

    await wrapper.vm.deleteDocument(1)
    await flushPromises()

    expect(wrapper.vm.documents.length).toBe(0)

    confirmSpy.mockRestore()
  })

  // ---------------------------- //
  // DELETE CANCELLED             //
  // ---------------------------- //

  it('does not delete document if user cancels', async () => {
    fetch.mockResolvedValue({ ok: true, json: async () => [] })

    const confirmSpy = vi.spyOn(window, 'confirm').mockReturnValue(false)

    const wrapper = mount(DocumentLibrary)

    wrapper.vm.documents = [
      { id: 1, title: 'Test Doc' }
    ]

    await wrapper.vm.deleteDocument(1)

    expect(wrapper.vm.documents.length).toBe(1)

    confirmSpy.mockRestore()
  })

  // ---------------------------- //
  // CHAT INTEGRATION             //
  // ---------------------------- //

  it('calls chatbox setActiveDocument when chat is opened', async () => {
    fetch.mockResolvedValue({ ok: true, json: async () => [] })

    const wrapper = mount(DocumentLibrary, {
      global: {
        stubs: {
          Chatbox: {
            template: '<div></div>',
            methods: {
              setActiveDocument: vi.fn()
            }
          }
        }
      }
    })

    const mockDoc = { id: 1, title: 'Test' }

    // manually mock ref
    wrapper.vm.chatboxRef = {
      setActiveDocument: vi.fn()
    }

    wrapper.vm.openChat(mockDoc)

    expect(wrapper.vm.chatboxRef.setActiveDocument).toHaveBeenCalledWith(mockDoc)
  })

})
