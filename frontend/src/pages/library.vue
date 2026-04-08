<template>
  <div class="page-container">
    <h2 class="page-title">Document Library</h2>

    <!-- Upload Section -->
    <div class="upload-card">
      <h3 class="section-title">Upload Documents</h3>

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
        <p class="doc-sub">{{ doc.type }}</p>

        <div class="doc-actions">
          <button @click="openDocument(doc)">Open</button>
          <button @click="deleteDocument(doc.id)" class="delete-btn">
            Delete
          </button>
        </div>
      </div>
    </div>

    <!-- AI Chatbot Feature Box -->
    <div class="gemini-box">
      <Chatbox />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Chatbox from './chatbox.vue'

const selectedFile = ref(null)
const uploadMessage = ref('')
const error = ref('')

const documents = ref([])

// Load documents from backend (PERSISTENCE)
onMounted(() => {
  fetchDocuments()
})

async function fetchDocuments() {
  try {
    const res = await fetch('/api/documents')
    if (res.ok) {
      const data = await res.json()
      documents.value = data
    }
  } catch (err) {
    console.error(err)
  }
}

// Handle file selection
function handleFileUpload(e) {
  selectedFile.value = e.target.files[0]
}

// Upload file to backend
async function uploadFile() {
  if (!selectedFile.value) {
    error.value = 'Please select a file'
    return
  }

  error.value = ''
  uploadMessage.value = ''

  try {
    const formData = new FormData()
    formData.append('file', selectedFile.value)

    const res = await fetch('/api/documents/upload', {
      method: 'POST',
      body: formData
    })

    if (res.ok) {
      const data = await res.json()

      // backend response
      documents.value.push({
        id: data.id,
        title: data.name,
        type: data.type || 'File',
        url: data.url
      })

      uploadMessage.value = 'File uploaded successfully!'
      selectedFile.value = null
    } else {
      error.value = 'Upload failed'
    }
  } catch (err) {
    console.error(err)
    error.value = 'Server error'
  }
}

// Open document
function openDocument(doc) {
  if (doc.url) {
    window.open(doc.url, '_blank')
  } else {
    alert('No file URL available')
  }
}

// Delete document (backend + UI)
async function deleteDocument(id) {
  if (!confirm('Are you sure you want to delete this document?')) return

  try {
    const res = await fetch(`/api/documents/${id}`, {
      method: 'DELETE'
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