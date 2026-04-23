import { describe, it, expect, vi, beforeEach } from "vitest"
import { mount, flushPromises } from "@vue/test-utils"
import JobWorkspace from "@/pages/job-workspace.vue"

global.fetch = vi.fn()

const mockJob = {
	id: 1,
	title: "Software Engineer",
	company_name: "Google",
	location_text: "NYC",
	posting_url: "",
	salary: 100000,
	deadline_date: "2026-05-01T00:00:00Z",
	status: "applied",
	description: "Test desc",
	company_notes: "Initial notes",
	created_at: new Date().toISOString(),
	outcome: { status: "", notes: "" }
}

describe("JobWorkspace.vue", () => {

	beforeEach(() => {
		vi.clearAllMocks()
	})

	// 1. Rendering
	it("renders job workspace with fetched data", async () => {
		fetch
			.mockResolvedValueOnce({ ok: true, json: async () => mockJob })
			.mockResolvedValueOnce({ ok: true, json: async () => [] })
			.mockResolvedValueOnce({ ok: true, json: async () => [] })
			.mockResolvedValueOnce({ ok: true, json: async () => [] })

		const wrapper = mount(JobWorkspace, {
			props: { jobId: 1 }
		})

		await flushPromises()

		expect(wrapper.find(".workspace-title").text())
			.toContain("Software Engineer")

		expect(wrapper.find("form").exists()).toBe(true)
	})

	// 2. Modify Job Details + Submit
	it("updates job details and calls PUT /api/jobs", async () => {
		fetch
			.mockResolvedValueOnce({ ok: true, json: async () => mockJob }) // fetchJob
			.mockResolvedValueOnce({ ok: true, json: async () => [] })
			.mockResolvedValueOnce({ ok: true, json: async () => [] })
			.mockResolvedValueOnce({ ok: true, json: async () => [] })
			.mockResolvedValueOnce({ ok: true }) // PUT

		const wrapper = mount(JobWorkspace, { props: { jobId: 1 } })
		await flushPromises()

		const titleInput = wrapper.find("input[type='text']")
		await titleInput.setValue("New Title")

		await wrapper.find("form").trigger("submit.prevent")
		await flushPromises()

		expect(fetch).toHaveBeenCalledWith(
			"/api/jobs",
			expect.objectContaining({ method: "PUT" })
		)
	})

	// 3. Timeline Sorting Logic
	it("renders timeline sorted by newest first", async () => {
		const activities = [
			{ id: 1, activity_type: "applied", activity_at: "2024-01-01", description: "Old" },
			{ id: 2, activity_type: "interview", activity_at: "2025-01-01", description: "New" }
		]

		fetch
			.mockResolvedValueOnce({ ok: true, json: async () => mockJob })
			.mockResolvedValueOnce({ ok: true, json: async () => activities })
			.mockResolvedValueOnce({ ok: true, json: async () => [] })
			.mockResolvedValueOnce({ ok: true, json: async () => [] })

		const wrapper = mount(JobWorkspace, { props: { jobId: 1 } })
		await flushPromises()

		await wrapper.findAll(".workspace-tab")[1].trigger("click")

		const items = wrapper.findAll(".item-card")
		expect(items[0].text()).toContain("New")
		expect(items[1].text()).toContain("Old")
	})

	// 4. Add Interview
	it("adds interview and calls POST endpoint", async () => {
		fetch
			.mockResolvedValueOnce({ ok: true, json: async () => mockJob })
			.mockResolvedValue({ ok: true, json: async () => [] })

		const wrapper = mount(JobWorkspace, { props: { jobId: 1 } })
		await flushPromises()

		await wrapper.findAll(".workspace-tab")[2].trigger("click")

		const inputs = wrapper.findAll("input")
		await inputs[0].setValue("Technical")
		await inputs[1].setValue("2026-04-22T10:00")

		const addBtn = wrapper.findAll("button")
			.find(b => b.text().includes("Add Interview"))

		await addBtn.trigger("click")
		await flushPromises()

		expect(fetch).toHaveBeenCalledWith(
			"/api/jobs/1/interviews",
			expect.objectContaining({ method: "POST" })
		)
	})

	// 5. Interview Validation
	it("does not add interview if required fields missing", async () => {
		fetch
			.mockResolvedValueOnce({ ok: true, json: async () => mockJob })
			.mockResolvedValueOnce({ ok: true, json: async () => [] })
			.mockResolvedValueOnce({ ok: true, json: async () => [] })
			.mockResolvedValueOnce({ ok: true, json: async () => [] })

		const wrapper = mount(JobWorkspace, { props: { jobId: 1 } })
		await flushPromises()

		await wrapper.findAll(".workspace-tab")[2].trigger("click")

		const addBtn = wrapper.findAll("button")
			.find(b => b.text().includes("Add Interview"))

		await addBtn.trigger("click")

		expect(wrapper.text()).toContain("Round and date are required")
	})

	// 6. Company Notes Save
	it("saves company notes", async () => {
		fetch
			.mockResolvedValueOnce({ ok: true, json: async () => mockJob }) // job
			.mockResolvedValueOnce({ ok: true, json: async () => [] }) // activities
			.mockResolvedValueOnce({ ok: true, json: async () => [] }) // interviews
			.mockResolvedValueOnce({ ok: true, json: async () => [] }) // followups
			.mockResolvedValueOnce({ ok: true }) // PUT notes

		const wrapper = mount(JobWorkspace, { props: { jobId: 1 } })
		await flushPromises()

		await wrapper.findAll(".workspace-tab")[2].trigger("click")

		const textarea = wrapper.find("textarea")
		await textarea.setValue("Updated notes")

		const saveBtn = wrapper.findAll("button")
			.find(b => b.text().includes("Save Company Notes"))

		await saveBtn.trigger("click")
		await flushPromises()

		expect(fetch).toHaveBeenCalledWith(
			"/api/jobs/1/company-notes",
			expect.objectContaining({ method: "PUT" })
		)
	})

	// 7. AI Enhance Button
	it("shows and triggers Enhance with AI button", async () => {
		fetch
			.mockResolvedValueOnce({ ok: true, json: async () => mockJob }) // job
			.mockResolvedValueOnce({ ok: true, json: async () => [] }) // activities
			.mockResolvedValueOnce({ ok: true, json: async () => [] }) // interviews
			.mockResolvedValueOnce({ ok: true, json: async () => [] }) // followups
			.mockResolvedValueOnce({
				ok: true,
				json: async () => ({ enhanced_text: "Enhanced version" })
			}) // AI

		const wrapper = mount(JobWorkspace, { props: { jobId: 1 } })
		await flushPromises()

		await wrapper.findAll(".workspace-tab")[2].trigger("click")

		const enhanceBtn = wrapper.find(".action-button-enhance")
		expect(enhanceBtn.exists()).toBe(true)

		await enhanceBtn.trigger("click")
		await flushPromises()

		expect(fetch).toHaveBeenCalledWith(
			"/api/jobs/1/enhance-notes",
			expect.objectContaining({ method: "POST" })
		)

		expect(wrapper.find("textarea").element.value)
			.toBe("Enhanced version")
	})

})
