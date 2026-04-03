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
        <input v-model="form.company_name" type="text" required class="status-bar" />
      </div>

      <!-- LOCATION -->
      <div class="form-group">
        <label>Location</label>
        <input v-model="form.location_text" type="text" class="status-bar" />
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
        <input v-model="form.deadline_date" type="date" required class="status-bar" />
      </div>

      <!-- STATUS -->
      <div class="form-group">
        <label>Status</label>
        <select v-model="form.status" required class="status-bar">
            <option disabled value="">Select status</option>
            <option value="interested">Interested</option>
            <option value="applied">Applied</option>
            <option value="interview">Interview</option>
            <option value="offer">Offer</option>
            <option value="rejected">Rejected</option>        
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
  company_name: '',
  title: '',
  location_text: '',
  salary: 0,
  deadline_date: '',
  status: '',
  description: ''
})

const error = ref('')

/* ---------------- SUBMIT ---------------- */
async function handleSubmit() {
  try {
    error.value = ''

    // Basic validation
    if (!form.title || !form.company_name || !form.deadline_date || !form.status) {
      error.value = 'Please fill in all required fields'
      return
    }

    const res = await fetch('/api/jobs', {
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
