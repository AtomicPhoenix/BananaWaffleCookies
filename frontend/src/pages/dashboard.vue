<template>
  <div class="dashboard">
    <div class="dashboard-header">
      <h1 class="dashboard-title">My Job Dashboard</h1>
      <h2 class="dashboard-subtitle">Welcome back, [Name]!</h2>
    </div>
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
      <!--Filter By-->
        <div class="filter-box">
          Filter By: 
          <select>
            <option value="all">All</option>
            <option value="status">Status</option>
            <option value="date">Date</option>
            <option value="salary">Salary</option>
          </select>
        </div>
    </div>

    <div class="job-info">

      <!-- OVERVIEW -->
      <div class="overview">
        <h1 class="overview-title">Overview</h1>
        <p class="overview-items">
          Interested: {{ stats.interested }}<br>
          Applied: {{ stats.applied }}<br>
          Interview: {{ stats.interview }}<br>
          Offer: {{ stats.offer }}<br>
          Accepted: {{ stats.accepted }}<br>
          Rejected: {{ stats.rejected }}<br>
        </p>
      </div>

      <!-- JOB LIST -->
      <div class="job-list">

        <div class="create-job">
          <router-link class="create-job-button" to="/create-job">
            Add a New Job Listing
          </router-link>
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
              {{ result.title }} | {{ result.company_name }} | {{ result.location_text }}
            </div>
            <div class="left mid">
              {{ result.deadline_date }}
              <linebreak>|</linebreak>
              {{ result.updated_at }}
            </div>
            <div class="right mid">
              <select
                class="listing-status-button listing-status-select"
                :id="statusToCssId(result.status)"
                :value="result.status"
                @change="updateJobStatus(result, $event.target.value)"
              >
                <option v-for="status in statusOptions" :key="status" :value="status">
                  {{ status }}
                </option>
              </select>
            </div>
            <div class="job-actions">
              <!-- Dropdown Menu for modify, archive, delete, etc. job -->
              <BDropdown auto-close="outside" class="dropdown" no-caret toggle-class="job-menu-toggle">
                <template #button-content>
                  <span aria-hidden="true">☰</span>
                  <span class="visually-hidden">Job actions</span>
                </template>
                <BDropdownItem :to="`/jobs/${result.id}`">View</BDropdownItem>
                <BDropdownItem :to="`/jobs/${result.id}/edit`">Modify</BDropdownItem>
                <BDropdownItem>Archive</BDropdownItem>
                <BDropdownItem>Delete</BDropdownItem>
              </BDropdown>
            </div>
          </div>
        </div>
        <!-- USER JOBS -->
        <div
          v-for="job in userJobs"
          :key="job.id"
          class="job-listing"
        >
          <div class="left top">
            {{ job.title }} | {{ job.company_name }} | {{ job.location_text }}
          </div>

          <div class="left mid jdesc">
            Last Modified: {{ formatDate(job.updated_at) }}
          </div>

          <div class="left bot jdesc">
            Deadline: {{ formatDate(job.deadline_date) }}
          </div>
          <div class="listing-status-button right mid" :id="statusToCssId(job.status)">
            {{ job.status }}
          </div>
        </div>

      </div>
    </div>
  </div>
</template>

<style scoped src="@/assets/css/dashboard.css"></style>

<script setup>
import { ref, onMounted } from 'vue'
import { BDropdown, BDropdownItem } from 'bootstrap-vue-next'

/* ---------------- STATE ---------------- */
const searchQuery = ref('')
const searchResults = ref([])
const userJobs = ref([])

const stats = ref({
  interested: 0,
  applied: 0,
  interview: 0,
  offer: 0,
  accepted: 0,
  rejected: 0
})

const statusOptions = ['interested', 'applied', 'interview', 'offer', 'accepted', 'rejected', 'archived']

/* ---------------- API CALLS ---------------- */

// Search jobs (external or internal API)
const handleSearch = async () => {
  try {
    const query = encodeURIComponent(searchQuery.value.trim())
    const url = query ? `/api/jobs?search=${query}` : `/api/jobs`
    const res = await fetch(url, { method: 'GET' })
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
    offer: 0,
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

const statusToCssId = (status) => {
  return String(status || '')
    .toLowerCase()
    .trim()
    .replace(/\s+/g, '-')
}

const updateJobStatus = async (job, newStatus) => {
  if (!newStatus || newStatus === job.status) {
    return
  }

  const previousStatus = job.status
  job.status = newStatus

  const matchingUserJob = userJobs.value.find((item) => item.id === job.id)
  if (matchingUserJob) {
    matchingUserJob.status = newStatus
  }
  computeStats(userJobs.value)

  try {
    const res = await fetch('/api/jobs', {
      method: 'PUT',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ ...job, status: newStatus })
    })

    if (!res.ok) {
      throw new Error(`Status update failed with code ${res.status}`)
    }
  } catch (err) {
    job.status = previousStatus
    if (matchingUserJob) {
      matchingUserJob.status = previousStatus
    }
    computeStats(userJobs.value)
    window.alert('Unable to update status right now. Please try again.')
    console.error('Failed to update status:', err)
  }
}

/* ---------------- LIFECYCLE ---------------- */

onMounted(() => {
  fetchUserJobs()
})
</script>
