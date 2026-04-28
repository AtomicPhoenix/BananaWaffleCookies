// performance.test.js

import { describe, it, expect, vi, beforeEach } from "vitest";
import { mount, flushPromises } from "@vue/test-utils";

import JobWorkspace from "@/pages/job-workspace.vue";
import LibraryPage from "@/pages/library.vue";

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

const mockDocuments = [
	{
		id: 1,
		title: "Resume 2026",
		type: "resume",
		status: "active",
		tags: ["frontend", "internship"],
		created_at: "2026-04-01T00:00:00.000Z",
		updated_at: "2026-04-02T00:00:00.000Z",
		versions: [],
	},
];

const mockJobs = [
	{
		id: 1,
		title: "Software Engineer",
		company_name: "OpenAI",
	},
];

function successfulFetch(data = []) {
	return Promise.resolve({
		ok: true,
		json: () => Promise.resolve(data),
		text: () => Promise.resolve(""),
		blob: () => Promise.resolve(new Blob()),
	});
}

describe("Performance Tests", () => {
	beforeEach(() => {
		vi.clearAllMocks();

		global.fetch = vi.fn((url) => {
			/* JOB WORKSPACE ROUTES */

			if (
				url.includes("/api/jobs/1") &&
				!url.includes("activities") &&
				!url.includes("interviews") &&
				!url.includes("followups") &&
				!url.includes("company-notes")
			) {
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

			if (url.includes("company-notes")) {
				return successfulFetch({});
			}

			/* LIBRARY ROUTES */

			if (url === "/api/documents") {
				return successfulFetch(mockDocuments);
			}

			if (url === "/api/jobs") {
				return successfulFetch(mockJobs);
			}

			return successfulFetch({});
		});
	});

	/*
	================================================
	JOB WORKSPACE PERFORMANCE TESTS
	================================================
	*/

	it("renders initial page load under performance threshold", async () => {
		const start = performance.now();

		const wrapper = mount(JobWorkspace);
		await flushPromises();

		const end = performance.now();
		const renderTime = end - start;

		expect(wrapper.exists()).toBe(true);
		expect(renderTime).toBeLessThan(500);
	});

	it("hydrates page with expected number of fetch requests efficiently", async () => {
		mount(JobWorkspace);
		await flushPromises();

		expect(global.fetch).toHaveBeenCalledTimes(4);
	});

	it("switches tabs quickly without extra network requests", async () => {
		const wrapper = mount(JobWorkspace);
		await flushPromises();

		const initialCalls = global.fetch.mock.calls.length;

		const start = performance.now();

		const buttons = wrapper.findAll(".workspace-tab");
		await buttons[1].trigger("click");

		const end = performance.now();
		const switchTime = end - start;

		expect(switchTime).toBeLessThan(100);
		expect(global.fetch).toHaveBeenCalledTimes(initialCalls);
	});

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

			if (
				url.includes("/api/jobs/1") &&
				!url.includes("interviews") &&
				!url.includes("followups")
			) {
				return successfulFetch(mockJobResponse);
			}

			return successfulFetch([]);
		});

		const start = performance.now();

		const wrapper = mount(JobWorkspace);
		await flushPromises();

		const buttons = wrapper.findAll(".workspace-tab");
		await buttons[1].trigger("click");
		await flushPromises();

		const end = performance.now();
		const renderTime = end - start;

		expect(wrapper.text()).toContain("Timeline");
		expect(renderTime).toBeLessThan(1000);
	});

	it("saves company notes within acceptable request timing", async () => {
		const wrapper = mount(JobWorkspace);
		await flushPromises();

		const start = performance.now();

		const saveButton = wrapper
			.findAll("button")
			.find((btn) => btn.text().includes("Save Company Notes"));

		if (saveButton) {
			await saveButton.trigger("click");
			await flushPromises();
		}

		const end = performance.now();
		const requestTime = end - start;

		expect(requestTime).toBeLessThan(500);
	});

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

	/*
	NEW TEST:
	Concurrent API resolution timing for Promise.all()
	*/

	it("handles concurrent API resolution timing efficiently during hydrateAll", async () => {
		global.fetch = vi.fn((url) => {
			return new Promise((resolve) => {
				setTimeout(() => {
					if (
						url.includes("/api/jobs/1") &&
						!url.includes("activities") &&
						!url.includes("interviews") &&
						!url.includes("followups")
					) {
						resolve(successfulFetch(mockJobResponse));
						return;
					}

					resolve(successfulFetch([]));
				}, 50);
			});
		});

		const start = performance.now();

		const wrapper = mount(JobWorkspace);
		await flushPromises();

		const end = performance.now();
		const totalTime = end - start;

		expect(wrapper.exists()).toBe(true);

		/*
		If sequential this would be much higher.
		This helps verify Promise.all concurrency behavior.
		*/
		expect(totalTime).toBeLessThan(400);
	});

	/*
	================================================
	LIBRARY PAGE PERFORMANCE TESTS
	================================================
	*/

	it("renders document library initial load under threshold", async () => {
		const start = performance.now();

		const wrapper = mount(LibraryPage);
		await flushPromises();

		const end = performance.now();
		const renderTime = end - start;

		expect(wrapper.exists()).toBe(true);
		expect(renderTime).toBeLessThan(500);
	});

	it("filters documents quickly without triggering additional API calls", async () => {
		const wrapper = mount(LibraryPage);
		await flushPromises();

		const initialCalls = global.fetch.mock.calls.length;

		const start = performance.now();

		const filterInput = wrapper.find('input[placeholder="Filter by tag"]');
		await filterInput.setValue("frontend");
		await flushPromises();

		const end = performance.now();
		const filterTime = end - start;

		expect(filterTime).toBeLessThan(100);
		expect(global.fetch).toHaveBeenCalledTimes(initialCalls);
	});

	it("renders large document dataset efficiently", async () => {
		const largeDocuments = Array.from({ length: 300 }, (_, i) => ({
			id: i,
			title: `Document ${i}`,
			type: "resume",
			status: "active",
			tags: ["tag1"],
			created_at: "2026-04-01T00:00:00.000Z",
			updated_at: "2026-04-02T00:00:00.000Z",
			versions: [],
		}));

		global.fetch = vi.fn((url) => {
			if (url === "/api/documents") {
				return successfulFetch(largeDocuments);
			}

			if (url === "/api/jobs") {
				return successfulFetch(mockJobs);
			}

			return successfulFetch([]);
		});

		const start = performance.now();

		const wrapper = mount(LibraryPage);
		await flushPromises();

		const end = performance.now();
		const renderTime = end - start;

		expect(wrapper.text()).toContain("Document Library");
		expect(renderTime).toBeLessThan(1000);
	});
});
