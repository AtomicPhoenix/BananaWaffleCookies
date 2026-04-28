<template>
  <div class="page-container">
    <div v-if="activeDocumentName" class="active-doc">
      Analyzing: {{ activeDocumentName }}
    </div>

    <h2 class="page-title">Document Library</h2>

    <!-- UPLOAD SECTION -->
    <div class="upload-card">
      <h3 class="section-title">Upload Documents</h3>

      <div class="form-group">
        <label>Select Job</label>
        <select v-model="selectedJobId">
          <option disabled value="">Select a job</option>
          <option v-for="job in jobs" :key="job.id" :value="job.id">
            {{ job.title }} - {{ job.company_name }}
          </option>
        </select>
      </div>

      <div class="form-group">
        <label>Document Type</label>
        <select v-model="documentType">
          <option value="resume">Resume</option>
          <option value="cover_letter">Cover Letter</option>
          <option value="other">Other</option>
        </select>
      </div>

      <div class="form-group">
        <label>Tags</label>

        <div class="tag-input-row">
          <input v-model="tagInput" placeholder="e.g. frontend" />
          <button @click="addTag">Add</button>
        </div>

        <div class="tag-list">
          <span v-for="tag in selectedTags" :key="tag" class="tag">
            {{ tag }}
          </span>
        </div>
      </div>

      <input type="file" @change="handleFileUpload" />
      <button @click="uploadFile()">Upload File</button>

      <p v-if="uploadMessage" class="success">{{ uploadMessage }}</p>
      <p v-if="error" class="error">{{ error }}</p>
    </div>

    <!-- FILTERS -->
    <div class="filters">
      <select v-model="filterType">
        <option value="">All Types</option>
        <option value="resume">Resume</option>
        <option value="cover_letter">Cover Letter</option>
      </select>

      <select v-model="filterStatus">
        <option value="">All Status</option>
        <option value="active">Active</option>
        <option value="archived">Archived</option>
      </select>

      <input v-model="filterTag" placeholder="Filter by tag" />

      <select v-model="sortBy">
        <option value="newest">Newest</option>
        <option value="oldest">Oldest</option>
      </select>
    </div>

    <!-- DOCUMENTS GRID -->
    <div class="documents-grid">
      <div v-for="doc in filteredDocuments" :key="doc.id" class="document-card">
        <div v-if="doc.isEditing">
          <input v-model="doc.title" />
          <button @click="saveTitle(doc)">Save</button>
        </div>

        <div v-else>
          <h4>{{ doc.title }}</h4>
          <button class="text-button" @click="doc.isEditing = true">Rename</button>
        </div>

        <p class="doc-sub">
          Owner: {{ doc.owner_name || 'You' }}
        </p>

        <div class="doc-sub">
          {{ doc.type }} • Status:
          <select v-model="doc.status" @change="updateStatus(doc)">
            <option value="active">Active</option>
            <option value="archived">Archived</option>
          </select>
        </div>

        <p class="doc-sub">
          Last Updated: {{ formatDate(doc.updated_at || doc.created_at) }}
        </p>

        <div class="doc-tags" v-if="doc.tags?.length">
          <span v-for="tag in doc.tags" :key="tag" class="tag">
            {{ tag }}
          </span>
        </div>

        <div class="versions" v-if="doc.versions?.length">
          <p><strong>Versions:</strong></p>

          <div
            v-for="v in doc.versions"
            :key="v.version_number"
            class="version-row"
          >
            v{{ v.version_number }} — {{ formatDate(v.created_at) }}

            <button @click="openVersion(v)">Open</button>
            <button @click="downloadVersion(v)">Download</button>
          </div>
        </div>

        <div class="doc-actions">
          <input class="browse-input" type="file" @change="(e) => handleFileUpload(e, doc)" />

          <button @click="uploadFile(doc)">
            Upload New Version
          </button>

          <button @click="duplicateDocument(doc)">Duplicate</button>

          <button v-if="doc.status === 'active'" @click="archiveDocument(doc)">
            Archive
          </button>

          <button v-else @click="restoreDocument(doc)">
            Restore
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'

const activeDocumentName = ref('')

async function openVersion(v) {
  try {
    const res = await fetch(`/api/documents/${v.document_id}`, {
      credentials: 'include'
    })

    if (!res.ok) throw new Error('Failed to open document')

    const blob = await res.blob()
    const url = window.URL.createObjectURL(blob)

    window.open(url, '_blank')

    // cleanup
    setTimeout(() => window.URL.revokeObjectURL(url), 1000)
  } catch (err) {
    console.error(err)
    error.value = 'Failed to open document'
  }
}

async function downloadVersion(v) {
  try {
    const res = await fetch(`/api/documents/${v.document_id}`, {
      credentials: 'include'
    })

    if (!res.ok) throw new Error('Download failed')

    const blob = await res.blob()
    const url = window.URL.createObjectURL(blob)

    const a = document.createElement('a')
    a.href = url
    a.download = 'document.pdf'
    document.body.appendChild(a)
    a.click()
    a.remove()

    window.URL.revokeObjectURL(url)
  } catch (err) {
    console.error(err)
    error.value = 'Download failed'
  }
}

function formatDate(date) {
  return new Date(date).toLocaleString()
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

const filterType = ref('')
const filterStatus = ref('')
const filterTag = ref('')
const sortBy = ref('newest')

/* COMPUTED */
const filteredDocuments = computed(() => {
  let docs = Array.isArray(documents.value) ? [...documents.value] : []

  const typeFilter = filterType.value?.toLowerCase()
  const statusFilter = filterStatus.value?.toLowerCase()
  const tagFilter = filterTag.value?.toLowerCase()

  if (typeFilter) {
    docs = docs.filter(d =>
      String(d.type || '').toLowerCase() === typeFilter
    )
  }

  if (statusFilter) {
    docs = docs.filter(d =>
      String(d.status || '').toLowerCase() === statusFilter
    )
  }

  if (tagFilter) {
    docs = docs.filter(d =>
      (d.tags || []).some(tag =>
        String(tag).toLowerCase().includes(tagFilter)
      )
    )
  }

  if (sortBy.value === 'newest') {
    docs.sort(
      (a, b) =>
        new Date(b.updated_at || b.created_at) -
        new Date(a.updated_at || a.created_at)
    )
  } else if (sortBy.value === 'oldest') {
    docs.sort(
      (a, b) =>
        new Date(a.updated_at || a.created_at) -
        new Date(b.updated_at || b.created_at)
    )
  }

  return docs
})

/* FETCH DATA */
onMounted(() => {
  fetchDocuments()
  fetchJobs()
})

async function fetchDocuments() {
  try {
    const res = await fetch('/api/documents', { credentials: 'include' })
    if (res.ok) {
      documents.value = (await res.json()).map(doc => ({
        ...doc,
        status: doc.status || 'active',
        tags: doc.tags || [],
        versions: doc.versions || [],
        updated_at: doc.updated_at || doc.created_at
      }))
    }
  } catch (err) {
    console.error(err)
  }
}

async function fetchDocument(id) {
  const res = await fetch(`/api/documents/${id}/info`, {
    credentials: 'include'
  })

  if (!res.ok) throw new Error('Failed to fetch document')

  const doc = await res.json()

  const normalized = {
    ...doc,
    status: doc.status || 'active',
    tags: doc.tags || [],
    versions: doc.versions || [],
    updated_at: doc.updated_at || doc.created_at
  }

  const index = documents.value.findIndex(d => d.id === id)

  if (index !== -1) {
    documents.value[index] = normalized
  } else {
    documents.value.unshift(normalized)
  }

  return normalized
}

async function fetchJobs() {
  try {
    const res = await fetch('/api/jobs', { credentials: 'include' })
    if (res.ok) {
      jobs.value = await res.json()
    }
  } catch (err) {
    console.error(err)
  }
}

/* TAG LOGIC */
function addTag() {
  const tag = tagInput.value.trim()
  if (tag && !selectedTags.value.includes(tag)) {
    selectedTags.value.push(tag)
  }
  tagInput.value = ''
}

/* FILE HANDLING */
function handleFileUpload(event, doc = null) {
  error.value = ''
  uploadMessage.value = ''

  const file = event.target.files?.[0]
  if (!file) return

  const allowedTypes = ['application/pdf']
  const maxSize = 5 * 1024 * 1024

  if (!allowedTypes.includes(file.type)) {
    error.value = 'Only PDF files are allowed'
    event.target.value = ''
    return
  }

  if (file.size > maxSize) {
    error.value = 'File must be under 5MB'
    event.target.value = ''
    return
  }

  if (doc) {
    doc._newFile = file
  } else {
    selectedFile.value = file
  }

  uploadMessage.value = `Selected: ${file.name}`
}

async function uploadFile(existingDoc = null) {
  uploadMessage.value = ''
  error.value = ''

  const fileToUpload = existingDoc
    ? existingDoc._newFile
    : selectedFile.value

  if (!fileToUpload) {
    error.value = 'Please select a file'
    return
  }

  const formData = new FormData()
  formData.append('file', fileToUpload)

  try {
    let res

    if (existingDoc) {
      res = await fetch(`/api/documents/${existingDoc.id}/versions`, {
        method: 'POST',
        body: formData,
        credentials: 'include'
      })
    } else {
      formData.append('type', documentType.value)
      formData.append('status', 'active')
      formData.append('job_id', selectedJobId.value)
      formData.append('tags', JSON.stringify(selectedTags.value))
      formData.append('title', fileToUpload.name)

      res = await fetch('/api/documents', {
        method: 'POST',
        body: formData,
        credentials: 'include'
      })
    }

    if (!res.ok) {
      throw new Error(await res.text())
    }

    await fetchDocuments()

    if (existingDoc) {
      existingDoc._newFile = null
    } else {
      selectedFile.value = null
      selectedTags.value = []
      selectedJobId.value = ''
    }

    uploadMessage.value = 'Upload successful!'
  } catch (err) {
    console.error(err)
    error.value = err.message || 'Upload failed'
  }
}

async function updateStatus(doc) {
  const oldStatus = doc.status

  try {
    const res = await fetch(`/api/documents/${doc.id}/status`, {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ status: doc.status })
    })

    if (!res.ok) throw new Error()

    await fetchDocument(doc.id)
  } catch (err) {
    console.error(err)
    doc.status = oldStatus
    error.value = 'Failed to update status'
  }
}

async function duplicateDocument(doc) {
  try {
    const res = await fetch(`/api/documents/${doc.id}/duplicate`, {
      method: 'POST',
      credentials: 'include'
    })

    if (!res.ok) throw new Error()

    await fetchDocuments()
  } catch (err) {
    console.error(err)
    error.value = 'Duplicate failed'
  }
}

async function saveTitle(doc) {
  try {
    const res = await fetch(`/api/documents/${doc.id}`, {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ title: doc.title })
    })

    if (!res.ok) throw new Error()

    doc.isEditing = false

    await fetchDocument(doc.id)
  } catch (err) {
    console.error(err)
    error.value = 'Rename failed'
  }
}

async function archiveDocument(doc) {
  doc.status = 'archived'
  await updateStatus(doc)
}

async function restoreDocument(doc) {
  doc.status = 'active'
  await updateStatus(doc)
}
</script>

<style scoped src="@/assets/css/library.css"></style>
