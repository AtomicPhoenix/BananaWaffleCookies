import { describe, it, expect, vi, beforeEach } from "vitest";
import { mount, flushPromises } from "@vue/test-utils";

import Profile from "@/pages/profile.vue";

// global fetch mock for ALL endpoints, might have to move this into a setup.js but temp solution for this sprint
global.fetch = vi.fn((url, options) => {
  // Handle GET endpoints
  if (!options || options.method === "GET") {
    if (url === "/api/profile") {
      return Promise.resolve({
        ok: true,
        json: async () => ({
          first_name: "",
          last_name: "",
          phone: "",
          preferences: {
            target_roles: "",
            location: "",
          },
        }),
      });
    }

    if (
      url.includes("/education") ||
      url.includes("/skills") ||
      url.includes("/experiences") ||
      url.includes("/projects")
    ) {
      return Promise.resolve({
        ok: true,
        json: async () => [],
      });
    }
  }

  // Handle PUT endpoints
  if (options?.method === "PUT") {
    return Promise.resolve({
      ok: true,
      json: async () => ({}),
    });
  }

  // Default fallback
  return Promise.resolve({
    ok: false,
    json: async () => ({}),
  });
});

describe("Profile.vue", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  // ---------------------------- //
  // RENDER TEST                  //
  // ---------------------------- //

  it("renders profile page correctly", () => {
    const wrapper = mount(Profile);

    expect(wrapper.text()).toContain("Profile");
    expect(wrapper.text()).toContain("Profile Completion");
  });

  // ---------------------------- //
  // INPUT BINDING TEST           //
  // ---------------------------- //

  it("updates form fields when user types", async () => {
    const wrapper = mount(Profile);

    await flushPromises();

    const inputs = wrapper.findAll("input");

    await inputs[0].setValue("Testname");
    await inputs[1].setValue("nameTest");

    expect(wrapper.vm.form.first_name).toContain("Testname");
    expect(wrapper.vm.form.last_name).toContain("nameTest");
  });
/* 
  // ---------------------------- //
  // COMPLETION PERCENTAGE TEST   //
  // ---------------------------- //

  it("calculates completion percentage correctly", async () => {
    const wrapper = mount(Profile);

    wrapper.vm.form.first_name = "John";
    wrapper.vm.form.last_name = "Doe";

    await wrapper.vm.$nextTick();

    const expected = Math.round((2 / 7) * 100);

    expect(wrapper.vm.completionPercentage).toBe(expected);
  });

  // ---------------------------- //
  // GET PROFILE (onMounted)      //
  // ---------------------------- //

  it("loads profile data from API on mount", async () => {
    const mockData = {
      first_name: "John",
      last_name: "Doe",
      phone: "1234567890",
      preferences: {
        target_roles: "Engineer",
        location: "NJ",
      },
    };

    // Override ONLY the profile endpoint
    fetch.mockImplementationOnce((url) => {
      if (url === "/api/profile") {
        return Promise.resolve({
          ok: true,
          json: async () => mockData,
        });
      }
      return global.fetch(url);
    });

    const wrapper = mount(Profile);

    await flushPromises();

    expect(wrapper.vm.form.first_name).toBe("John");
    expect(wrapper.vm.preferences.location).toBe("NJ");
    expect(fetch).toHaveBeenCalledWith("/api/profile");
  });
*/
  // ---------------------------- //
  // SAVE PROFILE (PUT)           //
  // ---------------------------- //

  it("sends updated basic profile data to API on save", async () => {
    const wrapper = mount(Profile);

    wrapper.vm.form.first_name = "John";
    wrapper.vm.form.last_name = "Doe";

    await wrapper.vm.$nextTick();

    const saveButton = wrapper
      .findAll("button")
      .find((btn) => btn.text().includes("Save Basic Info"));

    await saveButton.trigger("click");

    expect(fetch).toHaveBeenCalledWith("/api/profile/basic", {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        first_name: "John",
        last_name: "Doe",
        phone: "",
      }),
    });
  });

  // ---------------------------- //
  // EMPTY FORM COMPLETION        //
  // ---------------------------- //

  it("returns 0% completion when form is empty", () => {
    const wrapper = mount(Profile);

    expect(wrapper.vm.completionPercentage).toBe(0);
  });

  // ---------------------------- //
  // FULL FORM COMPLETION         //
  // ---------------------------- //

  it("returns 100% completion when all fields are filled", async () => {
    const wrapper = mount(Profile);

    Object.keys(wrapper.vm.form).forEach((key) => {
      wrapper.vm.form[key] = "test";
    });

    Object.keys(wrapper.vm.preferences).forEach((key) => {
      wrapper.vm.preferences[key] = "test";
    });

    await wrapper.vm.$nextTick();

    expect(wrapper.vm.completionPercentage).toBe(100);
  });

  // ---------------------------- //
  // ERROR HANDLING (GET)         //
  // ---------------------------- //

  it("handles API error gracefully on GET", async () => {
    fetch.mockImplementationOnce(() =>
      Promise.reject(new Error("API failure")),
    );

    const consoleSpy = vi.spyOn(console, "error").mockImplementation(() => {});

    mount(Profile);

    await flushPromises();

    expect(consoleSpy).toHaveBeenCalled();
    consoleSpy.mockRestore();
  });

  // ---------------------------- //
  // ERROR HANDLING (PUT)         //
  // ---------------------------- //

  it("handles API error gracefully on save", async () => {
    fetch.mockImplementation((url, options) => {
      if (options?.method === "PUT") {
        return Promise.reject(new Error("API failure"));
      }
      return Promise.resolve({
        ok: true,
        json: async () => ({}),
      });
    });

    const wrapper = mount(Profile);

    await flushPromises(); // ensure mount calls finish

    wrapper.vm.form.first_name = "John";
    wrapper.vm.form.last_name = "Doe";

    await wrapper.vm.$nextTick();

    const saveButton = wrapper
      .findAll("button")
      .find((btn) => btn.text().includes("Save Basic Info"));

    await saveButton.trigger("click");
    await flushPromises();

    expect(wrapper.vm.messages.basic.error).toBe("Server error");
  });
});
