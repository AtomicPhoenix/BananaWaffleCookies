// performance.test.js

import { describe, it, expect, vi, beforeEach } from "vitest";
import { mount, flushPromises } from "@vue/test-utils";
import JobWorkspace from "@/pages/job-workspace.vue";

vi.mock("vue-router", () => ({
	useRoute: () => ({
		params: {
			job_id: "1",
		},
	}),
}));

const mockJobResponse = {
	id: 1,
	title: "Software Engineer",
	company_name: "OpenAI",
	location_text: "Remote",
	posting_url: "https://example.com",
	salary: 120000,
	deadline_date: "2026-05-01T00:00:00.000Z",
	status: "applied",
	description: "Test description",
	company_notes: "Initial notes",
	created_at: "2026-04-01T00:00:00.000Z",
	outcome: {
		status: "",
		notes: "",
	},
};

function successfulFetch(data = []) {
	return Promise.resolve({
		ok: true,
		json: () => Promise.resolve(data),
	});
}

describe("JobWorkspace Performance Tests", () => {
	beforeEach(() => {
		vi.clearAllMocks();

		global.fetch = vi.fn((url) => {
			if (url.includes("/api/jobs/1") && !url.includes("activities") && !url.includes("interviews") && !url.includes("followups")) {
				return successfulFetch(mockJobResponse);
			}

			if (url.includes("activities")) {
				return successfulFetch([]);
			}

			if (url.includes("interviews")) {
				return successfulFetch([]);
			}

			if (url.includes("followups")) {
				return successfulFetch([]);
			}

			return successfulFetch({});
		});
	});

	/*
	TEST 1:
	Initial render should complete within acceptable threshold
	*/
	it("renders initial page load under performance threshold", async () => {
		const start = performance.now();

		const wrapper = mount(JobWorkspace);

		await flushPromises();

		const end = performance.now();
		const renderTime = end - start;

		expect(wrapper.exists()).toBe(true);
		expect(renderTime).toBeLessThan(500); // threshold can be tuned
	});

	/*
	TEST 2:
	hydrateAll should trigger exactly 4 API requests on mount
	*/
	it("hydrates page with expected number of fetch requests efficiently", async () => {
		mount(JobWorkspace);

		await flushPromises();

		expect(global.fetch).toHaveBeenCalledTimes(4);
	});

	/*
	TEST 3:
	Tab switching should be fast and not trigger unnecessary API calls
	*/
	it("switches tabs quickly without extra network requests", async () => {
		const wrapper = mount(JobWorkspace);

		await flushPromises();

		const initialCalls = global.fetch.mock.calls.length;

		const start = performance.now();

		const buttons = wrapper.findAll(".workspace-tab");
		await buttons[1].trigger("click"); // Timeline tab

		const end = performance.now();
		const switchTime = end - start;

		expect(switchTime).toBeLessThan(100);
		expect(global.fetch).toHaveBeenCalledTimes(initialCalls);
	});

	/*
	TEST 4:
	Large timeline dataset should render within reasonable threshold
	*/
	it("renders large activity dataset efficiently", async () => {
		const largeActivities = Array.from({ length: 200 }, (_, i) => ({
			id: i,
			activity_type: "interview_scheduled",
			activity_at: "2026-04-01T00:00:00.000Z",
			description: `Activity ${i}`,
		}));

		global.fetch = vi.fn((url) => {
			if (url.includes("activities")) {
				return successfulFetch(largeActivities);
			}

			if (url.includes("/api/jobs/1") && !url.includes("interviews") && !url.includes("followups")) {
				return successfulFetch(mockJobResponse);
			}

			return successfulFetch([]);
		});

		const start = performance.now();

		const wrapper = mount(JobWorkspace);

		await flushPromises();

		const buttons = wrapper.findAll(".workspace-tab");
		await buttons[1].trigger("click"); // Timeline

		await flushPromises();

		const end = performance.now();
		const renderTime = end - start;

		expect(wrapper.text()).toContain("Timeline");
		expect(renderTime).toBeLessThan(1000);
	});

	/*
	TEST 5:
	Save company notes should complete request quickly
	*/
	it("saves company notes within acceptable request timing", async () => {
		const wrapper = mount(JobWorkspace);

		await flushPromises();

		await wrapper.setData?.({}); // harmless if unsupported

		const start = performance.now();

		const saveButton = wrapper.findAll("button")
			.find(btn => btn.text().includes("Save Company Notes"));

		if (saveButton) {
			await saveButton.trigger("click");
			await flushPromises();
		}

		const end = performance.now();
		const requestTime = end - start;

		expect(requestTime).toBeLessThan(500);
	});

	/*
	TEST 6:
	Repeated re-renders should remain stable
	*/
	it("maintains stable performance during repeated tab navigation", async () => {
		const wrapper = mount(JobWorkspace);

		await flushPromises();

		const buttons = wrapper.findAll(".workspace-tab");

		const start = performance.now();

		for (let i = 0; i < 10; i++) {
			await buttons[1].trigger("click");
			await buttons[0].trigger("click");
		}

		const end = performance.now();
		const totalTime = end - start;

		expect(totalTime).toBeLessThan(1000);
	});
});
