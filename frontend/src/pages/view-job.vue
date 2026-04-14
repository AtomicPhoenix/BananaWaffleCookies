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
      </div>

      <!-- FEEDBACK -->
      <p v-if="error" class="error">{{ error }}</p>
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
    user_id.value = data.id
  } catch (err) {
    console.error('Failed to fetch user id:', err) 
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
</script>

<style scoped src="@/assets/css/job-page.css"></style>
