// library.test.js

import { describe, it, expect, vi, beforeEach } from "vitest";
import { mount, flushPromises } from "@vue/test-utils";
import LibraryPage from "@/pages/library.vue";

const mockDocuments = [
	{
		id: 1,
		title: "Resume.pdf",
		type: "resume",
		status: "active",
		tags: ["frontend"],
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

function successfulFetch(data = {}) {
	return Promise.resolve({
		ok: true,
		json: () => Promise.resolve(data),
		text: () => Promise.resolve(""),
		blob: () => Promise.resolve(new Blob()),
	});
}

function failedFetch(message = "Upload failed") {
	return Promise.resolve({
		ok: false,
		text: () => Promise.resolve(message),
	});
}

describe("Library Upload + Error Handling Tests", () => {
	beforeEach(() => {
		vi.clearAllMocks();

		global.fetch = vi.fn((url) => {
			if (url === "/api/documents") {
				return successfulFetch(mockDocuments);
			}

			if (url === "/api/jobs") {
				return successfulFetch(mockJobs);
			}

			if (url.includes("/info")) {
				return successfulFetch(mockDocuments[0]);
			}

			return successfulFetch({});
		});
	});

	/*
	TEST 1
	Reject unsupported file type
	*/
  
	it("rejects non-PDF uploads", async () => {
		const wrapper = mount(LibraryPage);
		await flushPromises();

		const fileInput = wrapper.find('input[type="file"]');

		const badFile = new File(["bad"], "image.png", {
			type: "image/png",
		});

		Object.defineProperty(fileInput.element, "files", {
			value: [badFile],
		});

		await fileInput.trigger("change");

		expect(wrapper.text()).toContain("Only PDF files are allowed");
	});

	/*
	TEST 2
	Reject oversized file
	*/

	it("rejects files larger than 5MB", async () => {
		const wrapper = mount(LibraryPage);
		await flushPromises();

		const fileInput = wrapper.find('input[type="file"]');

		const largeFile = new File(
			[new Array(6 * 1024 * 1024).fill("a").join("")],
			"large.pdf",
			{
				type: "application/pdf",
			}
		);

		Object.defineProperty(fileInput.element, "files", {
			value: [largeFile],
		});

		await fileInput.trigger("change");

		expect(wrapper.text()).toContain("File must be under 5MB");
	});

	/*
	TEST 3
	Attempt upload with no file selected
	*/

	it("prevents upload when no file is selected", async () => {
		const wrapper = mount(LibraryPage);
		await flushPromises();

		const uploadButton = wrapper
			.findAll("button")
			.find((btn) => btn.text() === "Upload File");

		await uploadButton.trigger("click");
		await flushPromises();

		expect(wrapper.text()).toContain("Please select a file");
	});

	/*
	TEST 4
	Successful valid PDF selection
	*/

	it("accepts valid PDF upload selection", async () => {
		const wrapper = mount(LibraryPage);
		await flushPromises();

		const fileInput = wrapper.find('input[type="file"]');

		const validFile = new File(["pdf"], "resume.pdf", {
			type: "application/pdf",
		});

		Object.defineProperty(fileInput.element, "files", {
			value: [validFile],
		});

		await fileInput.trigger("change");

		expect(wrapper.text()).toContain("Selected: resume.pdf");
	});

	/*
	TEST 5
	Successful upload triggers API + success message
	*/

	it("uploads valid document successfully", async () => {
		global.fetch = vi.fn((url) => {
			if (url === "/api/documents") {
				return successfulFetch(mockDocuments);
			}

			if (url === "/api/jobs") {
				return successfulFetch(mockJobs);
			}

			return successfulFetch({});
		});

		const wrapper = mount(LibraryPage);
		await flushPromises();

		const fileInput = wrapper.find('input[type="file"]');

		const validFile = new File(["pdf"], "resume.pdf", {
			type: "application/pdf",
		});

		Object.defineProperty(fileInput.element, "files", {
			value: [validFile],
		});

		await fileInput.trigger("change");

		const uploadButton = wrapper
			.findAll("button")
			.find((btn) => btn.text() === "Upload File");

		await uploadButton.trigger("click");
		await flushPromises();

		expect(wrapper.text()).toContain("Upload successful!");
	});

	/*
	TEST 6
	API upload failure catch
	*/

	it("handles upload API failure gracefully", async () => {
		global.fetch = vi.fn((url) => {
			if (url === "/api/documents") {
				return failedFetch("Server upload failed");
			}

			if (url === "/api/jobs") {
				return successfulFetch(mockJobs);
			}

			return successfulFetch({});
		});

		const wrapper = mount(LibraryPage);
		await flushPromises();

		const fileInput = wrapper.find('input[type="file"]');

		const validFile = new File(["pdf"], "resume.pdf", {
			type: "application/pdf",
		});

		Object.defineProperty(fileInput.element, "files", {
			value: [validFile],
		});

		await fileInput.trigger("change");

		const uploadButton = wrapper
			.findAll("button")
			.find((btn) => btn.text() === "Upload File");

		await uploadButton.trigger("click");
		await flushPromises();

		expect(wrapper.text()).toContain("Server upload failed");
	});

	/*
	TEST 7
	New version upload for existing document
	*/

	it("uploads new version for existing document", async () => {
		const wrapper = mount(LibraryPage);
		await flushPromises();

		const versionFileInput = wrapper.findAll('input[type="file"]')[1];

		const validFile = new File(["pdf"], "updated.pdf", {
			type: "application/pdf",
		});

		Object.defineProperty(versionFileInput.element, "files", {
			value: [validFile],
		});

		await versionFileInput.trigger("change");

		const versionButton = wrapper
			.findAll("button")
			.find((btn) => btn.text().includes("Upload New Version"));

		await versionButton.trigger("click");
		await flushPromises();

		expect(global.fetch).toHaveBeenCalled();
		expect(wrapper.text()).toContain("Upload successful!");
	});

	/*
	TEST 8
	Failed version upload API catch
	*/

	it("handles failed version upload correctly", async () => {
		global.fetch = vi.fn((url) => {
			if (url.includes("/versions")) {
				return failedFetch("Version upload failed");
			}

			if (url === "/api/documents") {
				return successfulFetch(mockDocuments);
			}

			if (url === "/api/jobs") {
				return successfulFetch(mockJobs);
			}

			return successfulFetch({});
		});

		const wrapper = mount(LibraryPage);
		await flushPromises();

		const versionFileInput = wrapper.findAll('input[type="file"]')[1];

		const validFile = new File(["pdf"], "updated.pdf", {
			type: "application/pdf",
		});
		Object.defineProperty(versionFileInput.element, "files", {
			value: [validFile],
		});

		await versionFileInput.trigger("change");

		const versionButton = wrapper
			.findAll("button")
			.find((btn) => btn.text().includes("Upload New Version"));

		await versionButton.trigger("click");
		await flushPromises();

		expect(wrapper.text()).toContain("Version upload failed");
	});
});
