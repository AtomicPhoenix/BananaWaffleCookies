<template>
  <div class="dashboard">
    <div class="dashboard-header">
      <h1 class="dashboard-title">My Job Dashboard</h1>
      <h2 class="dashboard-subtitle">{{ welcomeMessage }}</h2>
      <!-- OVERVIEW & ANALYTICS -->
      <div class="overview">
        <h1 class="overview-title">Pipeline Analytics</h1>
        
        <!-- Status Counts -->
        <div class="analytics-section">
          <h3 class="section-label">Job Status</h3>
          <div class="status-grid">
            <div class="stat-item">
              <span class="stat-label">Interested</span>
              <span class="stat-value">{{ stats.interested }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">Applied</span>
              <span class="stat-value">{{ stats.applied }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">Interview</span>
              <span class="stat-value">{{ stats.interview }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">Offer</span>
              <span class="stat-value">{{ stats.offer }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">Rejected</span>
              <span class="stat-value">{{ stats.rejected }}</span>
            </div>
          </div>
        </div>

        <!-- Stage Conversion -->
        <div class="analytics-section">
          <h3 class="section-label">Conversion Rates (% of Total Jobs)</h3>
          <div class="conversion-grid">
            <div class="conversion-item">
              <span class="conversion-label">Applied</span>
              <span class="conversion-value">{{ pipelineMetrics.applicationRate }}%</span>
            </div>
            <div class="conversion-item">
              <span class="conversion-label">Interviewed</span>
              <span class="conversion-value">{{ pipelineMetrics.interviewRate }}%</span>
            </div>
            <div class="conversion-item">
              <span class="conversion-label">Success (Offer)</span>
              <span class="conversion-value success">{{ pipelineMetrics.successRate }}%</span>
            </div>
            <div class="conversion-item">
              <span class="conversion-label">Rejected</span>
              <span class="conversion-value rejection">{{ pipelineMetrics.rejectionRate }}%</span>
            </div>
          </div>
        </div>

        <!-- Velocity -->
        <div class="analytics-section">
          <h3 class="section-label">Velocity (Last 30 Days)</h3>
          <div class="velocity-grid">
            <div class="velocity-item">
              <span class="velocity-label">Avg Jobs Added/Week</span>
              <span class="velocity-value">{{ pipelineMetrics.jobsAddedPerWeek }}</span>
            </div>
            <div class="velocity-item">
              <span class="velocity-label">Avg Applications/Week</span>
              <span class="velocity-value">{{ pipelineMetrics.applicationsPerWeek }}</span>
            </div>
            <div class="velocity-item">
              <span class="velocity-label">Avg Interviews/Week</span>
              <span class="velocity-value">{{ pipelineMetrics.interviewsPerWeek }}</span>
            </div>
            <div class="velocity-item">
              <span class="velocity-label">Avg Offers/Week</span>
              <span class="velocity-value">{{ pipelineMetrics.offersPerWeek }}</span>
            </div>
          </div>
        </div>

        <!-- Time in Stage -->
        <div class="analytics-section">
          <h3 class="section-label">Average Days in Stage</h3>
          <div class="time-grid">
            <div class="time-item">
              <span class="time-label">Interested</span>
              <span class="time-value">{{ pipelineMetrics.avgDaysInterested }}</span>
            </div>
            <div class="time-item">
              <span class="time-label">Applied</span>
              <span class="time-value">{{ pipelineMetrics.avgDaysApplied }}</span>
            </div>
            <div class="time-item">
              <span class="time-label">Interview</span>
              <span class="time-value">{{ pipelineMetrics.avgDaysInterview }}</span>
            </div>
            <div class="time-item">
              <span class="time-label">Offer</span>
              <span class="time-value">{{ pipelineMetrics.avgDaysOffer }}</span>
            </div>
          </div>
        </div>
      </div>
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
          <p class="filter-section-label">Other</p>
            <div class="filter-options">
              <label class="filter-option">
                  <input type="checkbox" v-model="showArchived">
                  Show Archived Jobs
                </label>
            </div>
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
              <BDropdownItem :to="{ name: 'job-detail', params: { job_id: job.id } }">View</BDropdownItem>
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
const showArchived = ref(false)

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

const pipelineMetrics = computed(() => {
  const activeJobs = userJobs.value.filter(j => !j.is_archived)
  const totalJobs = activeJobs.length
  const interested = stats.value.interested
  const applied = stats.value.applied
  const interview = stats.value.interview
  const offer = stats.value.offer
  const rejected = stats.value.rejected

  // Conversion rates as % of total jobs
  const applicationRate = totalJobs > 0 ? Math.round((applied / totalJobs) * 100) : 0
  const interviewRate = totalJobs > 0 ? Math.round((interview / totalJobs) * 100) : 0
  const successRate = totalJobs > 0 ? Math.round((offer / totalJobs) * 100) : 0
  const rejectionRate = totalJobs > 0 ? Math.round((rejected / totalJobs) * 100) : 0

  // Velocity metrics - calculate based on last 30 days
  const thirtyDaysAgo = new Date(Date.now() - 30 * 24 * 60 * 60 * 1000)
  const jobsIn30Days = activeJobs.filter(job => {
    const jobDate = new Date(job.created_at)
    return jobDate >= thirtyDaysAgo
  })
  
  const jobsByStatusIn30Days = {
    interested: jobsIn30Days.filter(j => j.status === 'interested'),
    applied: jobsIn30Days.filter(j => j.status === 'applied'),
    interview: jobsIn30Days.filter(j => j.status === 'interview'),
    offer: jobsIn30Days.filter(j => j.status === 'offer')
  }

  const weeksInPeriod = jobsIn30Days.length > 0 ? 4.29 : 1 // Approximate weeks in 30 days
  
  const jobsAddedPerWeek = jobsIn30Days.length > 0 ? (jobsIn30Days.length / weeksInPeriod).toFixed(1) : 0
  const applicationsPerWeek = jobsIn30Days.length > 0 ? (jobsByStatusIn30Days.applied.length / weeksInPeriod).toFixed(1) : 0
  const interviewsPerWeek = jobsIn30Days.length > 0 ? (jobsByStatusIn30Days.interview.length / weeksInPeriod).toFixed(1) : 0
  const offersPerWeek = jobsIn30Days.length > 0 ? (jobsByStatusIn30Days.offer.length / weeksInPeriod).toFixed(1) : 0

  // Calculate average time spent in each stage
  // For current jobs in a stage, use time since update
  // For past jobs that moved through, estimate from creation to now
  const calculateAvgDaysInStage = (statusLabel) => {
    const now = new Date().getTime()
    const hourDurations = []

    activeJobs.forEach(job => {
      const createdTime = new Date(job.created_at).getTime()
      const updatedTime = new Date(job.updated_at || job.created_at).getTime()

      if (job.status === statusLabel) {
        // Jobs currently in this stage: calculate from last update
        const hoursInStage = Math.floor((now - updatedTime) / (1000 * 60 * 60))
        hourDurations.push(hoursInStage)
      } else if (statusLabel === 'interested' && job.status !== 'interested') {
        // All jobs moved through "interested" stage - estimate from creation to now
        const hoursInStage = Math.floor((now - createdTime) / (1000 * 60 * 60))
        hourDurations.push(hoursInStage)
      } else if (statusLabel === 'applied' && (job.status === 'interview' || job.status === 'offer' || job.status === 'rejected')) {
        // Jobs that have moved past applied: estimate from creation
        const hoursInStage = Math.floor((now - createdTime) / (1000 * 60 * 60)) 
        hourDurations.push(Math.max(0, hoursInStage - 24)) // Rough estimate
      } else if (statusLabel === 'interview' && (job.status === 'offer' || job.status === 'rejected')) {
        // Jobs that moved past interview
        const hoursInStage = Math.floor((now - createdTime) / (1000 * 60 * 60)) 
        hourDurations.push(Math.max(0, hoursInStage - 72)) // Rough estimate
      }
    })

    if (hourDurations.length === 0) return '—'
    
    const avgHours = hourDurations.reduce((a, b) => a + b, 0) / hourDurations.length
    
    if (avgHours < 1) return '< 1h'
    if (avgHours < 24) return `${avgHours.toFixed(1)}h`
    
    const days = avgHours / 24
    if (days < 1) return `${avgHours.toFixed(1)}h`
    
    const wholeDays = Math.floor(days)
    const remainingHours = Math.round((days - wholeDays) * 24)
    return remainingHours === 0 ? `${wholeDays}d` : `${wholeDays}d ${remainingHours}h`
  }

  return {
    // Conversion rates
    applicationRate,
    interviewRate,
    successRate,
    rejectionRate,
    // Velocity
    jobsAddedPerWeek,
    applicationsPerWeek,
    interviewsPerWeek,
    offersPerWeek,
    // Time in stage
    avgDaysInterested: calculateAvgDaysInStage('interested'),
    avgDaysApplied: calculateAvgDaysInStage('applied'),
    avgDaysInterview: calculateAvgDaysInStage('interview'),
    avgDaysOffer: calculateAvgDaysInStage('offer')
  }
})

const statusOptions = ['interested', 'applied', 'interview', 'offer', 'rejected']
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

    if (counts[status] !== undefined && !job.is_archived) {
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
  const confirmed = window.confirm(
    `Delete "${job.title}" at ${job.company_name}? This action cannot be undone.`
  )

  if (!confirmed) {
    return
  }

  const typedConfirmation = window.prompt(
    'To permanently delete this job, type "delete" below.'
  )

  if ((typedConfirmation || '').trim().toLowerCase() !== 'delete') {
    window.alert('Delete cancelled. You must type "delete" exactly to confirm.')
    return
  }

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
  location.reload();
}

const archiveJob = async (job) => {
  try {
    const res = await fetch(`/api/jobs/${job.id}/archive`, {
      method: 'POST',
      credentials: 'include',
    })

    if (!res.ok) {
      return
    }

  } catch (err) {
    window.alert('Unable to archive job, please try again later.')
  }
  location.reload();
}

const unArchiveJob = async (job) => {
  try {
    const res = await fetch(`/api/jobs/${job.id}/unarchive`, {
      method: 'POST',
      credentials: 'include',
    })//blah

    if (!res.ok) {
      return
    }
  } catch (err) {
    window.alert('Unable to archive job, please try again later.')
  }
  location.reload();
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

const excludeArchivedUnlessEnabled = (job) => {
  if (showArchived.value) return true
  return !job.is_archived
}

const filteredJobs = (jobs) =>
  jobs
    .filter(excludeArchivedUnlessEnabled)
    .filter(matchesStatusFilter)
    .filter(matchesSalaryFilter)

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
