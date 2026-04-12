<template>
    <div class="job-page">
      <h1>Create a Job Application</h1>

      <!-- Form for Job Creation -->
      <form @submit.prevent="handleSubmit" class="job-form">
    
      <!-- Job Name -->
      <div class="form-job-name">
        <label for="jobName">Job Name</label>
        <input
          id="jobName"
          v-model="form.jobName"
          type="text"
          required
          placeholder="e.g. Software Engineer Intern"
          class="status-bar"
        />
      </div>

      <!-- Deadline -->
      <div class="form-deadline">
        <label for="deadline">Deadline</label>
        <input
          id="deadline"
          v-model="form.deadline"
          type="date"
          required
          class="status-bar"
        />
      </div>

      <!-- Application Status -->
      <div class="form-status">
        <label for="status">Application Status</label>
        <select id="status" v-model="form.status" required class="status-bar">
          <option disabled value="">Select status</option>
          <option>Not Applied</option>
          <option>Applied</option>
          <option>Interviewing</option>
          <option>Offer</option>
          <option>Rejected</option>
        </select>
      </div>

      <!-- Submit -->
      <button type="submit" class="submit-job-button">Submit</button>
      </form>
    </div>
</template>
<script setup>
import { reactive } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter();
const form = reactive({
  jobName: '',
  deadline: '',
  status: ''
})

function handleSubmit() {
  //FOR TESTING TO SEE
  console.log('Form Data:', { ...form })

  // Example: reset after submit
  form.jobName = ''
  form.deadline = ''
  form.status = ''

  //move back to dashboard page
  router.push('@/pages/dashboard')
}
</script>
<style scoped src="@/assets/css/job-page.css"></style>
