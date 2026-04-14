<template>
  <div class="page-container">
    <div v-if="activeDocumentName" class="active-doc">
      Analyzing: {{ activeDocumentName }}
    </div>

    <h2 class="page-title">Document Library</h2>

    <!-- Upload Section -->
    <div class="upload-card">
      <h3 class="section-title">Upload Documents</h3>

      <!-- JOB SELECT -->
      <div class="form-group">
        <label>Select Job</label>
        <select v-model="selectedJobId">
          <option disabled value="">Select a job</option>
          <option
            v-for="job in jobs"
            :key="job.id"
            :value="job.id"
          >
            {{ job.title }} - {{ job.company_name }}
          </option>
        </select>
      </div>

      <!-- DOCUMENT TYPE -->
      <div class="form-group">
        <label>Document Type</label>
        <select v-model="documentType">
          <option value="resume">Resume</option>
          <option value="cover_letter">Cover Letter</option>
          <option value="other">Other</option>
        </select>
      </div>

      <!-- TAG INPUT -->
      <div class="form-group">
        <label>Tags (Job Context)</label>

        <div class="tag-input-row">
          <input
            v-model="tagInput"
            placeholder="e.g. frontend, internship"
          />
          <button @click="addTag">Add</button>
        </div>

        <div class="tag-list">
          <span
            v-for="tag in selectedTags"
            :key="tag"
            class="tag"
          >
            {{ tag }}
          </span>
        </div>
      </div>

      <!-- FILE INPUT -->
      <input type="file" @change="handleFileUpload" />

      <button @click="uploadFile">
        Upload File
      </button>

      <p v-if="uploadMessage" class="success">{{ uploadMessage }}</p>
      <p v-if="error" class="error">{{ error }}</p>
    </div>

    <!-- Documents Grid -->
    <div class="documents-grid">
      <div
        v-for="doc in documents"
        :key="doc.id"
        class="document-card"
      >
        <div class="doc-preview">
          <p>{{ doc.title }}</p>
        </div>

        <h4>{{ doc.title }}</h4>

        <p class="doc-sub">
          {{ doc.type }} • job ID: {{ doc.job_id || 'N/A' }}
        </p>

        <!-- TAGS -->
        <div class="doc-tags" v-if="doc.tags && doc.tags.length">
          <span
            v-for="tag in doc.tags"
            :key="tag"
            class="tag"
          >
            {{ tag }}
          </span>
        </div>

        <!-- ACTIONS -->
        <div class="doc-actions">
          <button @click="openDocument(doc)">Open</button>
          <button @click="openChat(doc)">Chat</button>
          <button @click="deleteDocument(doc.id)" class="delete-btn">
            Delete
          </button>
        </div>
      </div>
    </div>

    <!-- AI Chatbot Feature Box -->
    <div class="gemini-box">
      <Chatbox ref="chatboxRef" />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Chatbox from '@/pages/chatbox.vue'

const chatboxRef = ref(null)

const activeDocumentName = ref('')

function openChat(doc) {
  chatboxRef.value?.setActiveDocument?.(doc)
  chatboxRef.value?.setActiveJobID?.(job)
}

/* STATE */
const selectedFile = ref(null)
const uploadMessage = ref('')
const error = ref('')

const documents = ref([])
const jobs = ref([])

const selectedJobId = ref('')
const selectedTags = ref([])
const tagInput = ref('')
const documentType = ref('resume')

/* FETCH DATA */
onMounted(() => {
  fetchDocuments()
  fetchJobs()
})

async function fetchDocuments() {
  try {
    const res = await fetch('/api/documents')
    if (res.ok) {
      documents.value = await res.json()
    }
  } catch (err) {
    console.error(err)
  }
}

async function fetchJobs() {
  try {
    const res = await fetch('/api/jobs')
    if (res.ok) {
      jobs.value = await res.json()
    }
  } catch (err) {
    console.error(err)
  }
}

/* TAG LOGIC */
function addTag() {
  if (tagInput.value.trim()) {
    selectedTags.value.push(tagInput.value.trim())
    tagInput.value = ''
  }
}

/* FILE HANDLING */
function handleFileUpload(e) {
  selectedFile.value = e.target.files[0]
}

/* UPLOAD */
async function uploadFile() {
  if (!selectedFile.value) {
    error.value = 'Please select a file'
    return
  }

  if (!selectedJobId.value) {
    error.value = 'Please select a job'
    return
  }

  error.value = ''
  uploadMessage.value = ''

  try {
    const formData = new FormData()
    formData.append('file', selectedFile.value)
    formData.append('job_id', selectedJobId.value)
    formData.append('tags', JSON.stringify(selectedTags.value))
    formData.append('type', documentType.value)

    const res = await fetch('/api/documents', {
      method: 'POST',
      body: formData,
      credentials: 'include'
    })

    if (res.ok) {
      const data = await res.json()

      documents.value.push({
        id: data.id,
        title: data.title || data.name,
        type: data.document_type || documentType.value || 'File',
        url: data.url || `/documents/${data.id}`,
        job_id: selectedJobId.value,
        tags: [...selectedTags.value]
      })

      uploadMessage.value = 'File uploaded successfully!'
      selectedFile.value = null
      selectedTags.value = []
    } else {
      error.value = 'Upload failed'
    }
  } catch (err) {
    console.error(err)
    error.value = 'Server error'
  }
}

/* ACTIONS */
function openDocument(doc) {
  if (doc.url) {
    window.open(doc.url, '_blank')
  } else {
    alert('No file URL available')
  }
}

async function deleteDocument(id) {
  if (!confirm('Are you sure you want to delete this document?')) return

  try {
    const res = await fetch(`/api/documents/${id}`, {
      method: 'DELETE',
      credentials: 'include'
    })

    if (res.ok) {
      documents.value = documents.value.filter(doc => doc.id !== id)
    } else {
      alert('Delete failed')
    }
  } catch (err) {
    console.error(err)
  }
}
</script>

<style scoped src="@/assets/css/library.css"></style>
