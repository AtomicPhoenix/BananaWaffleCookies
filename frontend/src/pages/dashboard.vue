<template>
  <div class="dashboard">
    <div class="dashboard-header">
      <h1 class="dashboard-title">My Job Dashboard</h1>
      <h2 class="dashboard-subtitle">{{ welcomeMessage }}</h2>
    </div>
    <!-- SEARCH -->
    <div class="search-box">
      <form @submit.prevent="handleSearch">
        <div class="search-controls">
          <div class="search-field">
            <label class="search-label" for="job-search">Search Jobs</label>
            <input
              id="job-search"
              v-model="searchQuery"
              class="search-bar"
              type="search"
              placeholder="Enter Job Information Here"
            />
            <button class="search-unicode" type="submit">⌕</button>
          </div>

          <div class="sort-box">
            <label class="sort-label" for="job-sort">Sort By</label>
            <select id="job-sort" v-model="sortBy" class="sort-select">
              <option value="updated-desc">Newest Updated</option>
              <option value="updated-asc">Oldest Updated</option>
              <option value="deadline-asc">Nearest Deadline</option>
              <option value="deadline-desc">Farthest Deadline</option>
              <option value="title-asc">Job Title (A-Z)</option>
              <option value="company-asc">Company Name (A-Z)</option>
            </select>
          </div>
        </div>
      </form>
    </div>

    <div class="job-info">
      <div class="dashboard-sidebar">
        <div class="status-filter-box">
          <h3 class="filters-title">Filters</h3>

          <p class="filter-section-label">Status</p>
          <div class="filter-options">
            <label v-for="status in statusOptions" :key="status" class="filter-option">
              <input
                v-model="selectedStatuses"
                type="checkbox"
                :value="status"
              >
              {{ status }}
            </label>
          </div>

          <p class="filter-section-label">Salary Range</p>
          <div class="filter-options">
            <label v-for="range in salaryRangeOptions" :key="range.value" class="filter-option">
              <input
                v-model="selectedSalaryRanges"
                type="checkbox"
                :value="range.value"
              >
              {{ range.label }}
            </label>
          </div>
        </div>

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
      </div>

      <!-- JOB LIST -->
      <div class="job-list">

        <div class="create-job">
          <router-link class="create-job-button" to="/jobs/create">
            Add a New Job Listing
          </router-link>
        </div>
        
        <h2 v-if="isSearchMode">Search Results: {{ searchQuery }}</h2>

        <div
          v-for="job in displayedJobs"
          :key="job.id"
          class="job-listing"
        >
          <div class="left top">
            {{ job.title }} | {{ job.company_name }} | {{ job.location_text }}
          </div>
          <div class="left mid jdesc">
            Deadline: {{ formatDate(job.deadline_date) }}
          </div>
          <div class="left bot jdesc">
            Last Modified: {{ formatDate(job.updated_at) }}
          </div>
          <div class="right mid">
            <select
              class="listing-status-button listing-status-select"
              :id="statusToCssId(job.status)"
              :value="job.status"
              @change="updateJobStatus(job, $event.target.value)"
            >
              <option v-for="status in statusOptions" :key="status" :value="status">
                {{ status }}
              </option>
            </select>
          </div>
          <div class="job-actions">
            <BDropdown auto-close="outside" class="dropdown" no-caret toggle-class="job-menu-toggle">
              <template #button-content>
                <span aria-hidden="true">☰</span>
                <span class="visually-hidden">Job actions</span>
              </template>
              <BDropdownItem :to="`/jobs/${job.id}`">View</BDropdownItem>
              <BDropdownItem :to="`/jobs/${job.id}/edit`">Modify</BDropdownItem>
              <BDropdownItem @click="archiveJob(job)">Archive</BDropdownItem>
              <BDropdownItem @click="unArchiveJob(job)">Restore</BDropdownItem>
              <BDropdownItem @click="deleteJob(job)">Delete</BDropdownItem>
            </BDropdown>
          </div>
        </div>

      </div>
    </div>
  </div>
</template>

<style scoped src="@/assets/css/dashboard.css"></style>

<script setup>
import { computed, ref, onMounted } from 'vue'
import { BDropdown, BDropdownItem } from 'bootstrap-vue-next'

/* ---------------- STATE ---------------- */
const searchQuery = ref('')
const searchResults = ref([])
const userJobs = ref([])
const sortBy = ref('updated-desc')
const selectedStatuses = ref([])
const selectedSalaryRanges = ref([])
const firstName = ref('')

const welcomeMessage = computed(() => {
  const name = String(firstName.value || '').trim()
  return name ? `Welcome back, ${name}` : 'Welcome back'
})

const stats = ref({
  interested: 0,
  applied: 0,
  interview: 0,
  offer: 0,
  accepted: 0,
  rejected: 0
})

const statusOptions = ['interested', 'applied', 'interview', 'offer', 'rejected', 'archived']
const salaryRangeOptions = [
  { value: 'under-50000', label: 'Under $50,000' },
  { value: '50000-74999', label: '$50,000 - $74,999' },
  { value: '75000-99999', label: '$75,000 - $99,999' },
  { value: '100000-149999', label: '$100,000 - $149,999' },
  { value: '150000-plus', label: '$150,000+' },
  { value: 'unknown', label: 'Not provided' }
]

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
    const res = await fetch('/api/jobs', {
      method: 'GET',
      credentials: 'include' // important if using sessions/cookies
    })

    if (!res.ok) {
      throw new Error(`Failed to fetch jobs with code ${res.status}`)
    }

    const data = await res.json()

    userJobs.value = data

    computeStats(data)
  } catch (err) {
    console.error('Failed to fetch user jobs:', err)
  }
}

const fetchUserProfile = async () => {
  try {
    const res = await fetch('/api/profile', {
      method: 'GET',
      credentials: 'include'
    })

    if (!res.ok) {
      return
    }

    const profileData = await res.json()
    firstName.value = profileData.first_name || ''
  } catch (err) {
    console.error('Failed to fetch profile:', err)
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

const deleteJob = async (job) => {
  try {
    const res = await fetch(`/api/jobs/${job.id}`, {
      method: 'DELETE',
      credentials: 'include',
    })

    if (!res.ok) {
      return
    }
  } catch (err) {
    window.alert('Unable to delete job, please try again later.')
  }
}

const archiveJob = async (job) => {
  try {
    const res = await fetch(`/api/jobs/${job.id}/archive`, {
      method: 'PUT',
      credentials: 'include',
    })

    if (!res.ok) {
      return
    }
  } catch (err) {
    window.alert('Unable to archive job, please try again later.')
  }
}

const unArchiveJob = async (job) => {
  try {
    const res = await fetch(`/api/jobs/${job.id}/unarchive`, {
      method: 'PUT',
      credentials: 'include',
    })//blah

    if (!res.ok) {
      return
    }
  } catch (err) {
    window.alert('Unable to archive job, please try again later.')
  }
}

const toTimeValue = (value) => {
  const timestamp = new Date(value).getTime()
  return Number.isNaN(timestamp) ? 0 : timestamp
}

const matchesStatusFilter = (job) => {
  if (!selectedStatuses.value.length) {
    return true
  }

  return selectedStatuses.value.includes(statusToCssId(job.status))
}

const salaryValue = (job) => {
  const parsed = Number(job.salary)
  return Number.isFinite(parsed) ? parsed : null
}

const matchesSalaryFilter = (job) => {
  if (!selectedSalaryRanges.value.length) {
    return true
  }

  const salary = salaryValue(job)

  return selectedSalaryRanges.value.some((range) => {
    if (range === 'unknown') {
      return salary === null
    }

    if (salary === null) {
      return false
    }

    switch (range) {
      case 'under-50000':
        return salary < 50000
      case '50000-74999':
        return salary >= 50000 && salary <= 74999
      case '75000-99999':
        return salary >= 75000 && salary <= 99999
      case '100000-149999':
        return salary >= 100000 && salary <= 149999
      case '150000-plus':
        return salary >= 150000
      default:
        return false
    }
  })
}

const filteredJobs = (jobs) => jobs.filter(matchesStatusFilter).filter(matchesSalaryFilter)

const sortedJobs = (jobs) => {
  return [...jobs].sort((left, right) => {
    switch (sortBy.value) {
      case 'updated-asc':
        return toTimeValue(left.updated_at) - toTimeValue(right.updated_at)
      case 'updated-desc':
        return toTimeValue(right.updated_at) - toTimeValue(left.updated_at)
      case 'deadline-asc':
        return toTimeValue(left.deadline_date) - toTimeValue(right.deadline_date)
      case 'deadline-desc':
        return toTimeValue(right.deadline_date) - toTimeValue(left.deadline_date)
      case 'title-asc':
        return String(left.title || '').localeCompare(String(right.title || ''))
      case 'company-asc':
        return String(left.company_name || '').localeCompare(String(right.company_name || ''))
      default:
        return 0
    }
  })
}

const sortedSearchResults = computed(() => sortedJobs(filteredJobs(searchResults.value)))
const sortedUserJobs = computed(() => sortedJobs(filteredJobs(userJobs.value)))
const isSearchMode = computed(() => searchQuery.value.trim().length > 0)
const displayedJobs = computed(() => {
  return isSearchMode.value ? sortedSearchResults.value : sortedUserJobs.value
})

/* ---------------- LIFECYCLE ---------------- */

onMounted(() => {
  fetchUserJobs()
  fetchUserProfile()
})
</script>
