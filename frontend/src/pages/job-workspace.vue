<template>
	<!-- Job Workspace Page -->
	<div class="job-workspace">
		<div class="workspace-header">
			<h1 class="workspace-title">{{ form.title || 'Job Workspace' }}</h1>
			<div class="workspace-actions">
				<BDropdown auto-close="outside" class="dropdown" no-caret toggle-class="job-menu-toggle">
					<template #button-content>
						<span aria-hidden="true">☰</span>
						<span class="visually-hidden">Job actions</span>
					</template>
					<BDropdownItem @click="archiveJob">Archive</BDropdownItem>
					<BDropdownItem @click="unArchiveJob">Restore</BDropdownItem>
					<BDropdownItem @click="deleteJob">Delete</BDropdownItem>
				</BDropdown>
			</div>
		</div>

		<div class="workspace-tabs"> 
			<button
				type="button"
				class="workspace-tab"
				:class="{ active: activeTab === 'details' }"
				@click="activeTab = 'details'"
			>
				Job Details
			</button>
			<button
				type="button"
				class="workspace-tab"
				:class="{ active: activeTab === 'timeline' }"
				@click="activeTab = 'timeline'"
			>
				Timeline
			</button>
			<button
				type="button"
				class="workspace-tab"
				:class="{ active: activeTab === 'documents' }"
				@click="activeTab = 'documents'"
			>
				Documents
			</button>
			<button
				type="button"
				class="workspace-tab"
				:class="{ active: activeTab === 'interviews' }"
				@click="activeTab = 'interviews'"
			>
				Interviews
			</button>
		</div>

		<p v-if="loading" class="feedback">Loading...</p>
		<p v-if="error" class="feedback error">{{ error }}</p>

		<div v-if="activeTab === 'details'" class="workspace-panel">
			<form class="job-form" @submit.prevent="saveJobDetails">
				<div class="form-group">
					<label>Job Title</label>
					<input v-model="form.title" type="text" required class="status-bar" />
				</div>

				<div class="form-group">
					<label>Company</label>
					<input v-model="form.company_name" type="text" required class="status-bar" />
				</div>

				<div class="form-group">
					<label>Salary</label>
					<input v-model="form.salary" type="number" class="status-bar" />
				</div>

				<div class="form-group">
					<label>Location</label>
					<input v-model="form.location_text" type="text" class="status-bar" />
				</div>

				<div class="form-group">
					<label>Posting URL</label>
					<input v-model="form.posting_url" type="url" class="status-bar" />
				</div>

				<div class="form-group">
					<label>Created</label>
					<input :value="formatDateTime(createdAt)" type="text" class="status-bar" readonly />
				</div>

				<div class="form-group">
					<label>Deadline</label>
					<input v-model="form.deadline_date" type="date" required class="status-bar" />
				</div>

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

				<div class="form-group">
					<label>Notes</label>
					<textarea v-model="form.description" class="status-bar"></textarea>
				</div>

				<button type="submit" class="action-button" :disabled="saving">
					{{ saving ? 'Saving...' : 'Save Job Details' }}
				</button>
			</form>
		</div>

		<div v-else-if="activeTab === 'timeline'" class="workspace-panel">
			<div class="section">
				<h3>Timeline</h3>
				<div v-if="!sortedActivities.length" class="empty-state">No timeline entries yet.</div>
				<div v-for="event in sortedActivities" :key="event.id" class="item-card">
					<strong>{{ formatActivityType(event.activity_type) }}</strong>
					<p class="sub-text">{{ formatDateTime(event.activity_at) }}</p>
					<p>{{ event.description }}</p>
				</div>
			</div>
		</div>

		<div v-else-if="activeTab === 'documents'" class="workspace-panel">
			<div class="section">
				<h3>Documents</h3>
				<div class="inline-form">
					<select v-model="selectedDocumentId">
						<option disabled value="">Select a document</option>
						<option v-for="doc in availableDocuments" :key="doc.id" :value="doc.id">
							{{ doc.title }}
						</option>
					</select>
					<div class="row-actions">
						<button type="button" class="action-button" @click="addSelectedDocument">
							Add Document
						</button>
					</div>
				</div>
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
        	<button @click="saveGeneratedResume">
          		Save to Job
        	</button>

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
			
			<div v-if="!availableDocuments.length" class="empty-state">No documents available.</div>
		</div>

		<div v-else class="workspace-panel">
			<div class="section">
				<h3>Interviews</h3>
				<div class="inline-form">
					<input v-model="newInterview.round" placeholder="Round Type" />
					<input v-model="newInterview.datetime" type="datetime-local" />
					<input v-model="newInterview.notes" placeholder="Notes" />
					<button type="button" class="action-button" @click="addInterview">Add Interview</button>
				</div>

				<div v-if="!interviews.length" class="empty-state">No interviews yet.</div>
				<div v-for="i in interviews" :key="i.id" class="item-card row-between">
					<div>
						<strong>{{ i.round }}</strong>
						<p class="sub-text">{{ formatDateTime(i.datetime) }}</p>
						<p>{{ i.notes || 'No notes' }}</p>
					</div>
					<button type="button" class="danger-button" @click="deleteInterview(i.id)">Delete</button>
				</div>
			</div>

			<div class="section">
				<h3>Follow-Ups / Reminders</h3>
				<div class="inline-form">
					<input v-model="newFollow.title" placeholder="Title" />
					<input v-model="newFollow.due_at" type="date" />
					<input v-model="newFollow.notes" placeholder="Notes (optional)" />
					<button type="button" class="action-button" @click="addFollowUp">Add Reminder</button>
				</div>

				<div v-if="!followUps.length" class="empty-state">No follow-ups yet.</div>
				<div v-for="f in followUps" :key="f.id" class="item-card row-between">
					<div>
						<strong>{{ f.title }}</strong>
						<p class="sub-text">{{ formatDateTime(f.due_at) }}</p>
						<p>{{ f.notes || 'No notes' }}</p>
						<p class="sub-text">Completed: {{ f.is_completed ? 'Yes' : 'No' }}</p>
					</div>
					<div class="row-actions">
						<button type="button" class="action-button" @click="toggleDone(f)">
							{{ f.is_completed ? 'Undo' : 'Done' }}
						</button>
						<button type="button" class="danger-button" @click="deleteFollowUp(f.id)">Delete</button>
					</div>
				</div>
			</div>

			<div class="section">
				<h3>Company Notes</h3>
				<textarea v-model="company_notes" class="status-bar"></textarea><br><br>

				<div class="row-actions">
					<button
						type="button"
						class="action-button"
						@click="saveCompanyNotes"
						:disabled="savingCompanyNotes"
					>
						{{ savingCompanyNotes ? 'Saving...' : 'Save Company Notes' }}
					</button>

					<button
						v-if="company_notes && company_notes.trim().length > 0"
						type="button"
						class="action-button-enhance"
						@click="enhanceCompanyNotes"
						:disabled="enhancingAI"
					>
						{{ enhancingAI ? 'Enhancing...' : 'Enhance with AI' }}
					</button>
				</div>
			</div>

			<div class="section">
				<h3>Outcome</h3>
				<div class="inline-form">
					<select v-model="outcome.status">
						<option value="">Select Status</option>
						<option>Offer</option>
						<option>Rejected</option>
						<option>Ghosted</option>
					</select>
					<input v-model="outcome.notes" placeholder="Notes" />
					<button type="button" class="action-button" @click="saveOutcome">Save Outcome</button>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup>
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { BDropdown, BDropdownItem } from 'bootstrap-vue-next'

const route = useRoute()
const props = defineProps({
	jobId: {
		type: [String, Number],
		default: null
	}
})

const resolvedJobId = computed(() => props.jobId ?? route.params.job_id)

const loading = ref(false)
const saving = ref(false)
const savingCompanyNotes = ref(false)
const error = ref('')
const activeTab = ref('details')
const company_notes = ref('')
const enhancingAI = ref(false)
const selectedDocumentId = ref('')
const availableDocuments = ref([])

const resumeResponse = ref('')
const isGeneratingResume = ref(false)
const coverLetterResponse = ref('')
const isGeneratingCoverLetter = ref(false)


const form = reactive({
	id: null,
	company_name: '',
	title: '',
	location_text: '',
	posting_url: '',
	salary: '',
	deadline_date: '',
	status: '',
	description: ''
})
	const createdAt = ref('')

const activities = ref([])
const interviews = ref([])
const followUps = ref([])

const newInterview = reactive({ round: '', datetime: '', notes: '' })
const newFollow = reactive({ title: '', due_at: '', notes: '' })
const outcome = reactive({ status: '', notes: '' })

const sortedActivities = computed(() =>
	[...activities.value].sort(
		(a, b) => new Date(b.activity_at) - new Date(a.activity_at)
	)
)

function toDateInput(value) {
	if (!value) return ''
	const date = new Date(value)
	if (Number.isNaN(date.getTime())) return ''
	return date.toISOString().split('T')[0]
}

function formatDate(value) {
	if (!value) return 'N/A'
	const date = new Date(value)
	if (Number.isNaN(date.getTime())) return 'N/A'
	return date.toLocaleDateString()
}

function formatDateTime(value) {
	if (!value) return 'N/A'
	const date = new Date(value)
	if (Number.isNaN(date.getTime())) return 'N/A'
	return date.toLocaleString()
}

function formatActivityType(type) {
	if (!type) return 'Activity'
	return String(type)
		.replaceAll('_', ' ')
		.replace(/\b\w/g, c => c.toUpperCase())
}

function addSelectedDocument() {
	return null
}

async function fetchJob() {
	if (!resolvedJobId.value) return

	loading.value = true
	error.value = ''

	try {
		const res = await fetch(`/api/jobs/${resolvedJobId.value}`, {
			method: 'GET',
			credentials: 'include'
		})

		if (!res.ok) {
			throw new Error('Failed to fetch job')
		}

		const data = await res.json()

		form.id = data.id
		form.company_name = data.company_name || ''
		form.title = data.title || ''
		form.location_text = data.location_text || ''
		form.posting_url = data.posting_url || ''
		form.salary = data.salary ?? ''
		form.deadline_date = toDateInput(data.deadline_date)
		form.status = data.status || ''
		form.description = data.description || ''
		company_notes.value = data.company_notes || ''
		createdAt.value = data.created_at || ''

		Object.assign(outcome, data.outcome || { status: '', notes: '' })
	} catch (err) {
		error.value = 'Failed to fetch job information.'
		console.error('Failed to fetch job:', err)
	} finally {
		loading.value = false
	}
}

async function fetchActivities() {
	if (!resolvedJobId.value) return

	try {
		const res = await fetch(`/api/jobs/${resolvedJobId.value}/activities`, {
			method: 'GET',
			credentials: 'include'
		})

		if (!res.ok) {
			activities.value = []
			return
		}

		const data = await res.json()
		activities.value = data || []
	} catch {
		activities.value = []
	}
}

async function fetchInterviews() {
	if (!resolvedJobId.value) return

	try {
		const res = await fetch(`/api/jobs/${resolvedJobId.value}/interviews`, {
			method: 'GET',
			credentials: 'include'
		})

		if (!res.ok) {
			interviews.value = []
			return
		}

		const data = await res.json()
		interviews.value = (data || []).map(i => ({
			id: i.id,
			round: i.round_type,
			datetime: i.scheduled_at,
			notes: i.notes
		}))
	} catch {
		interviews.value = []
	}
}

async function fetchFollowUps() {
	if (!resolvedJobId.value) return

	try {
		const res = await fetch(`/api/jobs/${resolvedJobId.value}/followups`, {
			method: 'GET',
			credentials: 'include'
		})

		if (!res.ok) {
			followUps.value = []
			return
		}

		followUps.value = await res.json()
	} catch {
		followUps.value = []
	}
}

async function hydrateAll() {
	await fetchJob()
	await Promise.all([fetchActivities(), fetchInterviews(), fetchFollowUps()])
}

async function saveJobDetails() {
	if (!resolvedJobId.value) return

	try {
		saving.value = true
		error.value = ''

		if (!form.title || !form.company_name || !form.deadline_date || !form.status) {
			error.value = 'Please fill in all required fields'
			return
		}

		const payload = {
			...form,
			id: form.id || Number(resolvedJobId.value),
			deadline_date: form.deadline_date ? new Date(form.deadline_date).toISOString() : null
		}

		const res = await fetch('/api/jobs', {
			method: 'PUT',
			headers: {
				'Content-Type': 'application/json'
			},
			credentials: 'include',
			body: JSON.stringify(payload)
		})

		if (!res.ok) {
			throw new Error('Failed to update job')
		}

		await fetchJob()
	} catch (err) {
		error.value = 'Error saving job details.'
		console.error(err)
	} finally {
		saving.value = false
	}
}

async function saveCompanyNotes() {
	if (!resolvedJobId.value) return

	try {
		savingCompanyNotes.value = true
		error.value = ''

		const res = await fetch(`/api/jobs/${resolvedJobId.value}/company-notes`, {
			method: 'PUT',
			headers: {
				'Content-Type': 'application/json'
			},
			credentials: 'include',
			body: JSON.stringify({
				company_notes: company_notes.value
			})
		})

		if (!res.ok) throw new Error()

		await fetchJob()
	} catch (err) {
		error.value = 'Unable to save company notes.'
		console.error(err)
	} finally {
		saving.value = false
	}
}

async function enhanceCompanyNotes() {
	if (!resolvedJobId.value || !company_notes.value.trim()) return

	try {
		enhancingAI.value = true
		error.value = ''

		const res = await fetch(`/api/jobs/${resolvedJobId.value}/enhance-notes`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			credentials: 'include',
			body: JSON.stringify({
				type: 'enhance_company_notes', //idk
				content: company_notes.value
			})
		})

		if (!res.ok) {
			throw new Error('AI enhancement failed, try again later')
		}

		const data = await res.json()

		// expecting something like: { enhanced_text: "..." }
		if (data?.enhanced_text) {
			company_notes.value = data.enhanced_text
		} else {
			throw new Error('Invalid AI response format')
		}

	} catch (err) {
		error.value = 'Unable to enhance notes right now.'
		console.error(err)
	} finally {
		enhancingAI.value = false
	}
}

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

async function addInterview() {
	if (!resolvedJobId.value) return

	if (!newInterview.round || !newInterview.datetime) {
		error.value = 'Round and date are required for interviews.'
		return
	}

	try {
		error.value = ''

		const res = await fetch(`/api/jobs/${resolvedJobId.value}/interviews`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			credentials: 'include',
			body: JSON.stringify({
				round_type: newInterview.round,
				scheduled_at: new Date(newInterview.datetime).toISOString(),
				notes: newInterview.notes
			})
		})

		if (!res.ok) {
			throw new Error('Failed to add interview')
		}

		newInterview.round = ''
		newInterview.datetime = ''
		newInterview.notes = ''
		await Promise.all([fetchInterviews(), fetchActivities()])
	} catch (err) {
		error.value = 'Unable to add interview right now.'
		console.error(err)
	}
}

async function deleteInterview(id) {
	if (!resolvedJobId.value) return

	try {
		const res = await fetch(`/api/jobs/${resolvedJobId.value}/interviews/${id}`, {
			method: 'DELETE',
			credentials: 'include'
		})

		if (!res.ok) return

		await Promise.all([fetchInterviews(), fetchActivities()])
	} catch (err) {
		console.error(err)
	}
}

async function addFollowUp() {
	if (!resolvedJobId.value) return

	if (!newFollow.title || !newFollow.due_at) {
		error.value = 'Title and due date are required for reminders.'
		return
	}

	try {
		error.value = ''

		const res = await fetch(`/api/jobs/${resolvedJobId.value}/followups`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			credentials: 'include',
			body: JSON.stringify({
				title: newFollow.title,
				due_at: new Date(newFollow.due_at).toISOString(),
				notes: newFollow.notes
			})
		})

		if (!res.ok) {
			throw new Error('Failed to add follow-up')
		}

		newFollow.title = ''
		newFollow.due_at = ''
		newFollow.notes = ''
		await Promise.all([fetchFollowUps(), fetchActivities()])
	} catch (err) {
		error.value = 'Unable to add follow-up right now.'
		console.error(err)
	}
}

async function deleteFollowUp(id) {
	if (!resolvedJobId.value) return

	try {
		await fetch(`/api/jobs/${resolvedJobId.value}/followups/${id}`, {
			method: 'DELETE',
			credentials: 'include'
		})

		await Promise.all([fetchFollowUps(), fetchActivities()])
	} catch (err) {
		console.error(err)
	}
}

async function toggleDone(followUp) {
	if (!resolvedJobId.value) return

	try {
		const res = await fetch(`/api/jobs/${resolvedJobId.value}/followups/${followUp.id}`, {
			method: 'PUT',
			headers: { 'Content-Type': 'application/json' },
			credentials: 'include',
			body: JSON.stringify({
				...followUp,
				is_completed: !followUp.is_completed
			})
		})

		if (!res.ok) {
			throw new Error('Failed to update follow-up')
		}

		await fetchFollowUps()
	} catch (err) {
		console.error(err)
	}
}

async function saveOutcome() {
	if (!resolvedJobId.value) return

	try {
		if (!outcome.status) {
			error.value = 'Outcome status required.'
			return
		}

		const res = await fetch(`/api/jobs/${resolvedJobId.value}/outcome`, {
			method: 'PUT',
			headers: { 'Content-Type': 'application/json' },
			credentials: 'include',
			body: JSON.stringify(outcome)
		})

		if (!res.ok) {
			throw new Error('Failed to save outcome')
		}

		await Promise.all([fetchJob(), fetchActivities()])
	} catch (err) {
		error.value = 'Unable to save outcome right now.'
		console.error(err)
	}
}

const deleteJob = async () => {
	const confirmed = window.confirm(
		`Delete "${form.title}" at ${form.company_name}? This action cannot be undone.`
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
		const res = await fetch(`/api/jobs/${resolvedJobId.value}`, {
			method: 'DELETE',
			credentials: 'include',
		})

		if (!res.ok) {
			return
		}
	} catch (err) {
		window.alert('Unable to delete job, please try again later.')
	}
	location.reload()
}

const archiveJob = async () => {
	try {
		const res = await fetch(`/api/jobs/${resolvedJobId.value}/archive`, {
			method: 'POST',
			credentials: 'include',
		})

		if (!res.ok) {
			return
		}
	} catch (err) {
		window.alert('Unable to archive job, please try again later.')
	}
	location.reload()
}

const unArchiveJob = async () => {
	try {
		const res = await fetch(`/api/jobs/${resolvedJobId.value}/unarchive`, {
			method: 'POST',
			credentials: 'include',
		})

		if (!res.ok) {
			return
		}
	} catch (err) {
		window.alert('Unable to restore job, please try again later.')
	}
	location.reload()
}

watch(
	resolvedJobId,
	() => {
		hydrateAll()
	}
)

onMounted(() => {
	hydrateAll()
})
</script>

<style scoped src="@/assets/css/job-workspace.css"></style>
