<template>
  <div class="job-page">
      <h1>{{job.title}}</h1>

      <div class="job-information-box">
        <ul class="job-list-column">
          <li><b>Company:</b> {{job.company_name}}</li>
          <li><b>Salary:</b> {{job.salary}}</li>
          <li><b>Location:</b> {{job.location_text}}</li>
          <li><b>Posting:</b> {{job.posting_url}}</li>
          <li><b>Date Applied:</b> {{job.date_applied}}</li>
          <li><b>Deadline:</b> {{job.deadline_date}}</li>
          <li><b>Status:</b> {{job.status}}</li>
        </ul>

        <h2 class="job-notes">Notes</h2>
        <p>{{job.description}}</p>

        <div class="edit-button-view" v-if="isOwner && job.id">
          <button @click="edit" class="edit-job-button">Edit</button>
        </div>


        <div class="generate-resume-button" v-if="isOwner && job.id">
          <button @click="generateResume" class="generate-resume-button">Edit</button>
        </div>
      </div>
  
    <!-- FEEDBACK -->
    <p v-if="error" class="error">{{ error }}</p>
  </div>
  <div class="job-page">
    <h1>Edit a Job Listing</h1>

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

      <!-- Salary -->
      <div class="form-group">
        <label>Salary</label>
        <input v-model="form.salary" type="number" class="status-bar" />
      </div>

      <!-- LOCATION -->
      <div class="form-group">
        <label>Location</label>
        <input v-model="form.location_text" type="text" class="status-bar" />
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
import { reactive, onMounted, watch, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()

/* ---------------- STATE ---------------- */
const isOwner = ref(false)
const user_id = ref(-1)
const job = ref({})

const error = ref('')

// Edit
const edit = () => {
  router.push({ name: 'edit-job', params: { job_id: route.params.job_id } })
}

// Fetch job
const fetchJob = async () => {
  try {
    let path = '/api/jobs/' + route.params.job_id
    const res = await fetch(path, {
      method: 'GET',
      credentials: 'include' // important if using sessions/cookies
    })
    job.value = await res.json()
  } catch (err) {
    console.error('Failed to fetch job:', err)
  }
}

// Fetch User ID
const getUser = async () => {
  try {
    const res = await fetch('/api/profile', {
      method: 'GET',
      credentials: 'include' // important if using sessions/cookies
    })
    let data = await res.json()
    user_id.value = data.user_id
  } catch (err) {
    console.error('Failed to fetch user id:', err) 
  }
}

// Get resume draft
const generateResume = async () => {
  try {
    let path = '/api/jobs/' + route.params.job_id + '/resume'
    const res = await fetch(path, {
      method: 'GET',
      credentials: 'include' // important if using sessions/cookies
    })
    job.value = await res.json()
  } catch (err) {
    console.error('Failed to generate job resume:', err)
  }
}


watch(() => route.params.job_id, (newId) => {
  form.id = newId
})

onMounted(async () => {
  await fetchJob()
  await getUser()
  isOwner.value = user_id.value === job.value.user_id
  console.log(user_id.value)
  console.log(job.value.user_id)
})

/* ---------------- MODIFY --------------- */
/* ---------------- STATE ---------------- */
const form = reactive({
  id: route.params.job_id,
  company_name: '',
  title: '',
  location_text: '',
  posting_url: '',
  salary: '',
  deadline_date: '',
  status: '',
  description: ''
})

watch(() => route.params.job_id, (newId) => {
  form.id = newId
})


/* ---------------- GET ---------------- */
// Fetch job
const fetchJob = async () => {
  try {
    let path = '/api/jobs/' + route.params.job_id
    const res = await fetch(path, {
      method: 'GET',
      credentials: 'include' // important if using sessions/cookies
    })
    const data = await res.json()

    // Convert to YYYY-MM-DD format
    if (data.deadline_date) {
      data.deadline_date = data.deadline_date.split('T')[0]
    }

    Object.assign(form, data) 
  } catch (err) {
    console.error('Failed to fetch job:', err)
  }
}


/* ---------------- SUBMIT ---------------- */
async function handleSubmit() {
  try {
    error.value = ''

    // Basic validation
    if (!form.title || !form.company_name || !form.deadline_date || !form.status) {
      error.value = 'Please fill in all required fields'
      return
    }

    form.deadline_date = new Date(form.deadline_date).toISOString()

    const res = await fetch('/api/jobs', {
      method: 'PUT',
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

onMounted(async () => {
  await fetchJob()
})

</script>

<style scoped src="@/assets/css/job-page.css"></style>
