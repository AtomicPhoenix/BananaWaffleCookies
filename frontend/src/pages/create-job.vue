<template>
  <div class="job-page">
    <h1>Create a Job Application</h1>

    <form @submit.prevent="handleSubmit" class="job-form">

      <!-- TITLE -->
      <div class="form-group">
        <label>Job Title</label>
        <input v-model="form.title" type="text" required class="status-bar" />
      </div>

      <!-- COMPANY -->
      <div class="form-group">
        <label>Company</label>
        <input v-model="form.company" type="text" required class="status-bar" />
      </div>

      <!-- LOCATION -->
      <div class="form-group">
        <label>Location</label>
        <input v-model="form.location" type="text" class="status-bar" />
      </div>

      <!-- SALARY -->
      <div class="form-group">
        <label>Salary</label>
        <input v-model="form.salary" type="number" class="status-bar" />
      </div>

      <!-- DATE APPLIED -->
      <div class="form-group">
        <label>Date Applied</label>
        <input v-model="form.date_applied" type="date" class="status-bar" />
      </div>

      <!-- DEADLINE -->
      <div class="form-group">
        <label>Deadline</label>
        <input v-model="form.deadline" type="date" required class="status-bar" />
      </div>

      <!-- STATUS -->
      <div class="form-group">
        <label>Status</label>
        <select v-model="form.status" required class="status-bar">
          <option disabled value="">Select status</option>
          <option>Interested</option>
          <option>Applied</option>
          <option>Interview</option>
          <option>Accepted</option>
          <option>Rejected</option>
        </select>
      </div>

      <!-- NOTES -->
      <div class="form-group">
        <label>Notes</label>
        <textarea v-model="form.description" class="status-bar"></textarea>
      </div>

      <!-- SUBMIT -->
      <button type="submit" class="submit-job-button">Submit</button>

      <!-- FEEDBACK -->
      <p v-if="error" class="error">{{ error }}</p>
    </form>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

/* ---------------- STATE ---------------- */
const form = reactive({
  title: '',
  company: '',
  location: '',
  salary: '',
  date_applied: '',
  deadline: '',
  status: '',
  description: ''
})

const error = ref('')

/* ---------------- SUBMIT ---------------- */
async function handleSubmit() {
  try {
    error.value = ''

    // Basic validation
    if (!form.title || !form.company || !form.deadline || !form.status) {
      error.value = 'Please fill in all required fields'
      return
    }

    const res = await fetch('/jobs', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      credentials: 'include',
      body: JSON.stringify(form)
    })

    if (!res.ok) {
      throw new Error('Failed to create job')
    }

    // Reset form
    Object.keys(form).forEach(key => form[key] = '')

    // Redirect (FIXED)
    router.push('/dashboard')

  } catch (err) {
    console.error(err)
    error.value = 'Error submitting form'
  }
}
</script>

<style scoped src="@/assets/css/job-page.css"></style>
