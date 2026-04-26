<template>
  <div class="page-container">
    <h2 class="page-title">{{ job.company_name }} — {{ job.title }}</h2>

    <div class="form-card">
      <!-- ================= TIMELINE ================= -->
      <div class="section">
        <h3 class="section-title">Timeline</h3>
          <div
            v-for="(event, index) in sortedActivities"
            :key="event.id"
            class="item-card"
          >
            <strong>{{ formatActivityType(event.activity_type) }}</strong>
            <p class="sub-text">{{ formatDate(event.activity_at) }}</p>
            <p class="activity-desc">{{ event.description }}</p>
          </div>
      </div>
      <!-- ================= INTERVIEWS ================= -->
      <div class="section">
        <h3 class="section-title">Interviews</h3>

        <input v-model="newInterview.round" placeholder="Round Type" />
        <input v-model="newInterview.datetime" type="datetime-local" />
        <input v-model="newInterview.notes" placeholder="Notes" />

        <button @click="addInterview">Add Interview</button>

        <p v-if="messages.interview.success" class="success">Saved!</p>
        <p v-if="messages.interview.error" class="error">
          {{ messages.interview.error }}
        </p>

        <div v-for="i in interviews" :key="i.id" class="item-card">
          <!-- VIEW -->
          <div v-if="editInterviewId !== i.id" class="item-row">
            <div>
              <strong>{{ i.round }}</strong>
              <p class="sub-text">{{ formatDate(i.datetime) }}</p>
              <p>{{ i.notes }}</p>
            </div>

            <div class="actions">
              <button @click="deleteInterview(i.id)">Delete</button>
            </div>
          </div>

          <!-- EDIT -->
          <div v-else class="edit-row">
            <input v-model="editInterview.round" />
            <input v-model="editInterview.datetime" type="datetime-local" />
            <input v-model="editInterview.notes" />

            <div class="actions">
              <button @click="updateInterview(i.id)">Save</button>
              <button @click="cancelEditInterview">Cancel</button>
            </div>
          </div>
        </div>
      </div>

      <!-- ================= FOLLOW-UPS ================= -->
      <div class="section">
        <h3 class="section-title">Follow-Ups / Reminders</h3>

        <input v-model="newFollow.title" placeholder="Title" />
        <input v-model="newFollow.due_at" type="date" />
        <input v-model="newFollow.notes" placeholder="Notes (optional)" />

        <button @click="addFollowUp">Add Reminder</button>

        <p v-if="messages.follow.success" class="success">Saved!</p>
        <p v-if="messages.follow.error" class="error">
          {{ messages.follow.error }}
        </p>

        <div v-for="f in followUps" :key="f.id" class="item-card">
          <!-- VIEW -->
          <div v-if="editFollowId !== f.id" class="item-row">
            <div>
              <strong>{{ f.title }}</strong>
              <p class="sub-text">{{ formatDate(f.due_at ) }}</p>
              <p v-if="f.notes">{{ f.notes }}</p>
            </div>

            <div class="actions">
              <button @click="startEditFollow(f)">Edit</button>
              <button @click="toggleDone(f)">
                {{ f.is_completed ? 'Undo' : 'Done' }}
              </button>
              <button @click="deleteFollowUp(f.id)">Delete</button>
            </div>
          </div>

          <!-- EDIT -->
          <div v-else class="edit-row">
            <input v-model="editFollow.title" />
            <input v-model="editFollow.due_at" type="date" />

            <div class="actions">
              <button @click="updateFollowUp(f.id)">Save</button>
              <button @click="cancelEditFollow">Cancel</button>
            </div>
          </div>
        </div>
      </div>

      <!-- ================= OUTCOME ================= -->
      <div class="section">
        <h3 class="section-title">Outcome</h3>

        <select v-model="outcome.status">
          <option value="">Select Status</option>
          <option>Offer</option>
          <option>Rejected</option>
          <option>Ghosted</option>
        </select>

        <input v-model="outcome.notes" placeholder="Notes" />

        <button @click="saveOutcome">Save Outcome</button>

        <p v-if="messages.outcome.success" class="success">Saved!</p>
        <p v-if="messages.outcome.error" class="error">
          {{ messages.outcome.error }}
        </p>
      </div>

      <!-- ================= RESUME GENERATION ================= -->
      <div class="section">
        <h3 class="section-title">AI Resume Generator</h3>
      
        <button @click="generateResume" :disabled="isGeneratingResume">
          {{ isGeneratingResume ? 'Generating...' : 'Generate Resume' }}
        </button>
      
        <div v-if="resumeResponse" class="item-card">
          <h4>Generated Resume</h4>
          <pre class="sub-text" style="white-space: pre-wrap;">
      {{ resumeResponse }}
          </pre>
        </div>
      </div>

      <!-- ================= SAVE RESUME ================= -->
      <div v-if="resumeResponse" class="item-card">
        <h4>Generated Resume</h4>
      
        <pre class="sub-text" style="white-space: pre-wrap;">
      {{ resumeResponse }}
        </pre>
      
        <button @click="saveGeneratedResume">
          Save to Job
        </button>
      </div>

      <!-- ================= COVER LETTER GENERATION ================= -->
      <div class="section">
        <h3 class="section-title">AI Cover Letter Generator</h3>
      
        <button @click="generateCoverLetter" :disabled="isGeneratingCoverLetter">
          {{ isGeneratingCoverLetter ? 'Generating...' : 'Generate Cover Letter' }}
        </button>
      
        <div v-if="coverLetterResponse" class="item-card">
          <h4>Generated Cover Letter</h4>
          <pre class="sub-text" style="white-space: pre-wrap;">
      {{ coverLetterResponse }}
          </pre>
        </div>
      </div>

      <button @click="saveGeneratedCoverLetter">
        Save to Job
      </button>

    </div>
  </div>
</template>

<script setup>
import { reactive, ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

// ================= STATE =================
const job = reactive({
  id: null,
  title: '',
  company_name: '',
  timeline: []
})

const resumeResponse = ref('')
const isGeneratingResume = ref(false)
const coverLetterResponse = ref('')
const isGeneratingCoverLetter = ref(false)

const interviews = ref([])
const followUps = ref([])
const activities = ref([])


const newInterview = reactive({ round: '', datetime: '', notes: '' })
const newFollow = reactive({
  title: '',
  due_at: '',
  notes: ''
})

const outcome = reactive({ status: '', notes: '' })

const editInterviewId = ref(null)
const editInterview = reactive({ round: '', datetime: '', notes: '' })

const editFollowId = ref(null)
const editFollow = reactive({
  title: '',
  due_at: '',
  notes: '',
  is_completed: false
})

const messages = reactive({
  interview: { success: false, error: '' },
  follow: { success: false, error: '' },
  outcome: { success: false, error: '' }
})

// ================= HELPERS =================
function reset(section) {
  messages[section].success = false
  messages[section].error = ''
}

// Get resume draft
const generateResume = async () => {
  isGeneratingResume.value = true
  resumeResponse.value = ''

  try {
    const path = `/api/jobs/${route.params.job_id}/resume`

    const res = await fetch(path, {
      method: 'POST',
      credentials: 'include'
    })

    const data = await res.json()

    if (data?.success) {
      resumeResponse.value = data.response
    } else {
      resumeResponse.value = 'Failed to generate resume.'
    }
  } catch (err) {
    console.error('Failed to generate job resume:', err)
    resumeResponse.value = 'Error generating resume.'
  } finally {
    isGeneratingResume.value = false
  }
}

// Get cover letter draft
const generateCoverLetter = async () => {
  isGeneratingCoverLetter.value = true
  coverLetterResponse.value = ''

  try {
    const path = `/api/jobs/${route.params.job_id}/cover-letter`

    const res = await fetch(path, {
      method: 'POST',
      credentials: 'include'
    })

    const data = await res.json()

    if (data?.success) {
      coverLetterResponse.value = data.response
    } else {
      coverLetterResponse.value = 'Failed to generate cover letter.'
    }
  } catch (err) {
    console.error('Failed to generate cover letter:', err)
    coverLetterResponse.value = 'Error generating cover letter.'
  } finally {
    isGeneratingCoverLetter.value = false
  }
}

async function saveGeneratedResume() {
  try {
    const res = await fetch(`/api/jobs/${job.id}/documents/ai-save`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({
        type: 'resume',
        content: resumeResponse.value
      })
    })

    const data = await res.json()

    if (!data.success) {
      alert('Failed to save resume')
      return
    }

    alert('Resume saved to job!')
  } catch (err) {
    console.error(err)
  }
}

async function saveGeneratedCoverLetter() {
  try {
    const res = await fetch(`/api/jobs/${job.id}/documents/ai-save`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({
        type: 'cover_letter',
        content: coverLetterResponse.value
      })
    })

    const data = await res.json()

    if (!data.success) {
      alert('Failed to save cover letter')
      return
    }

    alert('Cover letter saved to job!')
  } catch (err) {
    console.error(err)
  }
}

function formatDate(date) {
  return new Date(date).toLocaleString()
}

const sortedActivities = computed(() =>
  [...activities.value].sort(
    (a, b) => new Date(b.activity_at) - new Date(a.activity_at)
  )
)

// ================= Modified History =================

const getActivities = async () => {
  try {
    const res = await fetch(`/api/jobs/${job.id}/activities`, {
      method: 'GET',
      credentials: 'include'
    })

    if (!res.ok) {
      throw new Error(`Failed with status ${res.status}`)
    }

    const data = await res.json()

    activities.value = data || []

  } catch (err) {
    console.error('Failed to fetch activities:', err)
    activities.value = []
  }
}

function formatActivityType(type) {
  return type
    .replaceAll('_', ' ')
    .replace(/\b\w/g, c => c.toUpperCase())
}



// ================= FETCH =================
import { watch } from 'vue'

onMounted(() => {
  getJob(route.params.job_id)
})

watch(() => route.params.job_id, (newId) => {
  getJob(newId)
})

async function getJob(id) {
  const res = await fetch(`/api/jobs/${id}`, {
    credentials: 'include'
  })
  if (res.ok) {
    const data = await res.json()

    Object.assign(job, data)
    followUps.value = data.followUps || []
    Object.assign(outcome, data.outcome || {})

    if (!job.timeline || job.timeline.length === 0) {
      job.timeline = []
    }

    if (!job.timeline.some(e => e.type === 'applied')) {
      job.timeline.push({
        type: 'applied',
        date: job.created_at || new Date().toISOString(),
        note: 'Application submitted'
      })
    }
    await getActivities()
    await getFollowUps(id)
  }
}

// ================= INTERVIEWS =================
async function getInterviews(jobId) {
  try {
    const res = await fetch(`/api/jobs/${jobId}/interviews`, {
      credentials: 'include'
    })

    if (res.ok) {
      const data = await res.json()

      // map backend → frontend shape
      interviews.value = data.map(i => ({
        id: i.id,
        round: i.round_type,
        datetime: i.scheduled_at,
        notes: i.notes
      }))

      // rebuild timeline entries
      interviews.value.forEach(i => {
        job.timeline.push({
          type: 'interview',
          date: i.datetime,
          note: i.round
        })
      })
    }
  } catch (err) {
    console.error('Failed to fetch interviews', err)
  }
}

async function addInterview() {
  reset('interview')

  if (!newInterview.round || !newInterview.datetime) {
    messages.interview.error = 'Round and date required'
    return
  }

  try {
    const res = await fetch(`/api/jobs/${job.id}/interviews`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({
        round_type: newInterview.round,
        scheduled_at: new Date(newInterview.datetime).toISOString(),
        notes: newInterview.notes
      })
    })

    if (res.ok) {
      const data = await res.json()

      const newItem = {
        id: data.id,
        round: newInterview.round,
        datetime: newInterview.datetime,
        notes: newInterview.notes
      }

      interviews.value.push(newItem)

      job.timeline.push({
        type: 'interview',
        date: newItem.datetime,
        note: newItem.round
      })

      Object.keys(newInterview).forEach(k => (newInterview[k] = ''))
      messages.interview.success = true
    } else {
      messages.interview.error = 'Save failed'
    }
  } catch {
    messages.interview.error = 'Server error'
  }

}

async function deleteInterview(id) {
  try {
    const res = await fetch(`/api/jobs/${job.id}/interviews/${id}`, {
      method: 'DELETE',
      credentials: 'include'
    })

    if (!res.ok) return

    const old = interviews.value.find(i => i.id === id)

    interviews.value = interviews.value.filter(i => i.id !== id)

    job.timeline = job.timeline.filter(
      e => !(e.type === 'interview' && e.note === old.round)
    )
  } catch (err) {
    console.error(err)
  }

}

function startEditInterview(i) {
  editInterviewId.value = i.id
  Object.assign(editInterview, i)
}

function cancelEditInterview() {
  editInterviewId.value = null
}

async function updateInterview(id) {
  try {
    const res = await fetch(`/api/jobs/${job.id}/interviews/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(editInterview)
    })

    if (res.ok) {
      const updated = await res.json()

      const index = interviews.value.findIndex(i => i.id === id)
      interviews.value[index] = updated

      job.timeline = job.timeline.filter(
        e => !(e.type === 'interview' && e.note === updated.round)
      )

      job.timeline.push({
        type: 'interview',
        date: updated.datetime,
        note: updated.round
      })

      editInterviewId.value = null
    }
  } catch (err) {
    console.error(err)
  }

}

// ================= FOLLOW UPS =================
async function getFollowUps(jobId) {
  const res = await fetch(`/api/jobs/${jobId}/followups`, {
    credentials: 'include'
  })

  if (res.ok) {
    followUps.value = await res.json()
  }
}

async function addFollowUp() {
  reset('follow')

  if (!newFollow.title || !newFollow.due_at) {
    messages.follow.error = 'Title and due date required'
    return
  }

  try {
    const res = await fetch(`/api/jobs/${job.id}/followups`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({
        title: newFollow.title,
        due_at: new Date(newFollow.due_at).toISOString(),
        notes: newFollow.notes
      })
    })

    if (res.ok) {
      const saved = await res.json()

      followUps.value.push(saved)

      job.timeline.push({
        type: 'follow-up',
        date: saved.due_at,
        note: saved.title
      })

      Object.keys(newFollow).forEach(k => (newFollow[k] = ''))
      messages.follow.success = true
    } else {
      messages.follow.error = 'Save failed'
    }
  } catch {
    messages.follow.error = 'Server error'
  }

}

async function deleteFollowUp(id) {
  try {
    await fetch(`/api/jobs/${job.id}/followups/${id}`, {
      method: 'DELETE'
    })

    const old = followUps.value.find(f => f.id === id)
    followUps.value = followUps.value.filter(f => f.id !== id)

    job.timeline = job.timeline.filter(
      e => !(e.type === 'follow-up' && e.note === old.title)
    )
  } catch (err) {
    console.error(err)
  }

}

function startEditFollow(f) {
  editFollowId.value = f.id
  Object.assign(editFollow, {
    ...f,
    due_at: f.due_at ? f.due_at.slice(0, 10) : ''
  })
}

function cancelEditFollow() {
  editFollowId.value = null
}

async function updateFollowUp(id) {
  try {
    const res = await fetch(`/api/jobs/${job.id}/followups/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({
        ...editFollow,
        due_at: new Date(editFollow.due_at).toISOString()
      })
    })

    if (res.ok) {
      const updated = await res.json()

      const index = followUps.value.findIndex(f => f.id === id)
      followUps.value[index] = updated

      job.timeline = job.timeline.filter(
        e => !(e.type === 'follow-up' && e.note === updated.title)
      )

      job.timeline.push({
        type: 'follow-up',
        date: updated.due_at,
        note: updated.title
      })

      editFollowId.value = null
    }
  } catch (err) {
    console.error(err)
  }

}

async function toggleDone(f) {
  try {
    const res = await fetch(`/api/jobs/${job.id}/followups/${f.id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({
        ...f,
        is_completed: !f.is_completed
      })
    })

    if (res.ok) {
      const updated = await res.json()
      Object.assign(f, updated)
    }
  } catch (err) {
    console.error(err)
  }
}

// ================= OUTCOME =================
async function saveOutcome() {
  reset('outcome')

  if (!outcome.status) {
    messages.outcome.error = 'Status required'
    return
  }

  try {
    const res = await fetch(`/api/jobs/${job.id}/outcome`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(outcome)
    })

    if (res.ok) {
      job.timeline = job.timeline.filter(e => e.type !== 'outcome')

      job.timeline.push({
        type: 'outcome',
        date: new Date().toISOString(),
        note: `${outcome.status} — ${outcome.notes}`
      })

      messages.outcome.success = true
    }
  } catch {
    messages.outcome.error = 'Server error'
  }
}
</script>

<style scoped src="@/assets/css/profile.css"></style>
