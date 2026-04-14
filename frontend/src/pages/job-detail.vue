<template>
  <div class="page-container">
    <h2 class="page-title">{{ job.company_name }} — {{ job.title }}</h2>

    <div class="form-card">
      <!-- ================= TIMELINE ================= -->
      <div class="section">
        <h3 class="section-title">Activity Timeline</h3>

        <div
          v-for="(event, index) in sortedTimeline"
          :key="index"
          class="item-card"
        >
          <strong>{{ event.type.toUpperCase() }}</strong>
          <p class="sub-text">{{ formatDate(event.date) }}</p>
          <p>{{ event.note }}</p>
        </div>
      </div>

      <!-- ================= Last Modified ================= -->
      <div class="section">
        <h3 class="section-title">Last Modified</h3>
        <p class="sub-text">
          {{ timestamp ? formatDate(timestamp) : 'No activity yet' }}
        </p>
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
              <button @click="startEditInterview(i)">Edit</button>
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

        <input v-model="newFollow.task" placeholder="Task" />
        <input v-model="newFollow.date" type="date" />

        <button @click="addFollowUp">Add Reminder</button>

        <p v-if="messages.follow.success" class="success">Saved!</p>
        <p v-if="messages.follow.error" class="error">
          {{ messages.follow.error }}
        </p>

        <div v-for="f in followUps" :key="f.id" class="item-card">
          <!-- VIEW -->
          <div v-if="editFollowId !== f.id" class="item-row">
            <div>
              <strong>{{ f.task }}</strong>
              <p class="sub-text">{{ formatDate(f.date) }}</p>
            </div>

            <div class="actions">
              <button @click="startEditFollow(f)">Edit</button>
              <button @click="toggleDone(f)">
                {{ f.done ? 'Undo' : 'Done' }}
              </button>
              <button @click="deleteFollowUp(f.id)">Delete</button>
            </div>
          </div>

          <!-- EDIT -->
          <div v-else class="edit-row">
            <input v-model="editFollow.task" />
            <input v-model="editFollow.date" type="date" />

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
const timestamp = ref([null]) 

const newInterview = reactive({ round: '', datetime: '', notes: '' })
const newFollow = reactive({ task: '', date: '' })

const outcome = reactive({ status: '', notes: '' })

const editInterviewId = ref(null)
const editInterview = reactive({ round: '', datetime: '', notes: '' })

const editFollowId = ref(null)
const editFollow = reactive({ task: '', date: '' })

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

function formatDate(date) {
  return new Date(date).toLocaleString()
}

const sortedTimeline = computed(() =>
  [...job.timeline].sort((a, b) => new Date(b.date) - new Date(a.date))
)
// ================= Last Modified =================

const getTimestamp = async () => {
  try {
    const res = await fetch(`/api/jobs/${job.id}/activities`, {
      method: 'GET',
      credentials: 'include'
    })

    if (!res.ok) {
      throw new Error(`Failed with status ${res.status}`)
    }

    const data = await res.json()

    if (!Array.isArray(data) || data.length === 0) {
      timestamp.value = null
      return
    }

    // find most recent activity
    const latest = data.reduce((latest, current) => {
      return new Date(current.activity_at) > new Date(latest.activity_at)
        ? current
        : latest
    })

    timestamp.value = latest.activity_at

  } catch (err) {
    console.error('Failed to fetch timestamp:', err)
    timestamp.value = null
  }
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
    interviews.value = data.interviews || []
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

    await getTimestamp()
  }
}

// ================= INTERVIEWS =================
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
      body: JSON.stringify(newInterview)
    })

    if (res.ok) {
      const saved = await res.json()
      interviews.value.push(saved)

      job.timeline.push({
        type: 'interview',
        date: saved.datetime,
        note: saved.round
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
    await fetch(`/api/jobs/${job.id}/interviews/${id}`, {
      method: 'DELETE'
    })

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
async function addFollowUp() {
  reset('follow')

  if (!newFollow.task || !newFollow.date) {
    messages.follow.error = 'Task and date required'
    return
  }

  try {
    const res = await fetch(`/api/jobs/${job.id}/followups`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(newFollow)
    })

    if (res.ok) {
      const saved = await res.json()
      followUps.value.push(saved)

      job.timeline.push({
        type: 'follow-up',
        date: saved.date,
        note: saved.task
      })

      Object.keys(newFollow).forEach(k => (newFollow[k] = ''))
      messages.follow.success = true
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
      e => !(e.type === 'follow-up' && e.note === old.task)
    )
  } catch (err) {
    console.error(err)
  }
}

function startEditFollow(f) {
  editFollowId.value = f.id
  Object.assign(editFollow, f)
}

function cancelEditFollow() {
  editFollowId.value = null
}

async function updateFollowUp(id) {
  try {
    const res = await fetch(`/api/jobs/${job.id}/followups/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(editFollow)
    })

    if (res.ok) {
      const updated = await res.json()

      const index = followUps.value.findIndex(f => f.id === id)
      followUps.value[index] = updated

      job.timeline = job.timeline.filter(
        e => !(e.type === 'follow-up' && e.note === updated.task)
      )

      job.timeline.push({
        type: 'follow-up',
        date: updated.date,
        note: updated.task
      })

      editFollowId.value = null
    }
  } catch (err) {
    console.error(err)
  }
}

function toggleDone(f) {
  f.done = !f.done
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
