<template>
  <div class="dashboard">

    <!-- SEARCH -->
    <div class="search-box">
      <form @submit.prevent="handleSearch">
        <label class="search-label" for="job-search">Search Jobs</label>
        <input
          v-model="searchQuery"
          class="search-bar"
          type="search"
          placeholder="Enter Job Information Here"
        />
        <button class="search-unicode" type="submit">⌕</button>
      </form>
    </div>

    <div class="job-info">

      <!-- OVERVIEW (can be dynamic later) -->
      <div class="overview">
        <h1 class="overview-title">Overview</h1>
        <p class="overview-items">
          Interested: {{ stats.interested }}<br>
          Applied: {{ stats.applied }}<br>
          Interview Offered: {{ stats.interview }}<br>
          Accepted: {{ stats.accepted }}<br>
          Rejected: {{ stats.rejected }}<br>
        </p>
      </div>

      <!-- JOB LIST -->
      <div class="job-list">

        <div class="create-job">
          <router-link class="create-job-button" to="/create-job">
            Create New Job Application
          </router-link>
        </div>

        <!-- USER JOBS -->
        <div
          v-for="job in userJobs"
          :key="job.id"
          class="job-listing"
        >
          <div class="left top">
            {{ job.title }} | {{ job.company }} | {{ job.location }}
          </div>

          <div class="left mid jdesc">
            Last Modified: {{ formatDate(job.updated_at) }}
          </div>

          <div class="left bot jdesc">
            Deadline: {{ formatDate(job.deadline) }}
          </div>

          <div
            class="listing-status-button right mid"
            :id="job.status.toLowerCase()"
          >
            {{ job.status }}
          </div>
        </div>

      </div>

      <!-- OPTIONAL: SEARCH RESULTS -->
      <div v-if="searchResults.length">
        <h2>Search Results</h2>

        <div
          v-for="result in searchResults"
          :key="result.id"
          class="job-listing"
        >
          <div class="left top">
            {{ result.title }} | {{ result.company }} | {{ result.location }}
          </div>
        </div>
      </div>

    </div>
  </div>
</template>

<style scoped src="@/assets/css/dashboard.css"></style>

<script setup>
import { ref, onMounted } from 'vue'

/* ---------------- STATE ---------------- */
const searchQuery = ref('')
const searchResults = ref([])
const userJobs = ref([])

const stats = ref({
  interested: 0,
  applied: 0,
  interview: 0,
  accepted: 0,
  rejected: 0
})

/* ---------------- API CALLS ---------------- */

// Search jobs (external or internal API)
const handleSearch = async () => {
  try {
    const query = encodeURIComponent(searchQuery.value.trim())
    const res = await fetch(`/api/jobs/search?q=${query}`)
    const data = await res.json()

    searchResults.value = data
  } catch (err) {
    console.error('Search failed:', err)
  }
}

// Fetch logged-in user's jobs
const fetchUserJobs = async () => {

  try {
    const res = await fetch('/api/user/jobs', {
      //define GET right? not sure so not included
      method: 'GET',
      credentials: 'include' // important if using sessions/cookies
    })
    const data = await res.json()

    userJobs.value = data

    computeStats(data)
  } catch (err) {
    console.error('Failed to fetch user jobs:', err)
  }
}

/* ---------------- HELPERS ---------------- */

const computeStats = (jobs) => {
  const counts = {
    interested: 0,
    applied: 0,
    interview: 0,
    accepted: 0,
    rejected: 0
  }

  jobs.forEach(job => {
    const status = job.status.toLowerCase()

    if (counts[status] !== undefined) {
      counts[status]++
    }
  })

  stats.value = counts
}

const formatDate = (dateStr) => {
  if (!dateStr) return 'N/A'
  return new Date(dateStr).toLocaleDateString()
}

/* ---------------- LIFECYCLE ---------------- */

onMounted(() => {
  fetchUserJobs()
})
</script>
