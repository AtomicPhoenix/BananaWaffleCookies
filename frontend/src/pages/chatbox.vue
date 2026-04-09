<template>
  <div v-if="isOpen" class="chatbox-container">
    
    <!-- HEADER -->
    <div class="chatbox-header">
      <h3 class="chatbox-title">AI Resume Assistant</h3>
      <button class="close-btn" @click="closeChat">X</button>
    </div>

    <!-- MESSAGES -->
    <div class="chatbox-messages">
      <div
        v-for="(msg, index) in messages"
        :key="index"
        :class="['message', msg.role]"
      >
        {{ msg.content }}
      </div>

      <!-- LOADING BUBBLE -->
      <div v-if="isLoading" class="message ai loading">
        Thinking...
      </div>

      <!-- No Document Selected -->
      <div v-if="!activeDocumentId" class="message ai system">
        Please select a document from the library to begin.
      </div>
    </div>


    <!-- INPUT -->
    <div class="chatbox-input">
      <input
        v-model="userInput"
        type="text"
        :placeholder="activeDocumentId 
          ? 'Ask for feedback on your document...' 
          : 'Please select a document first...'"
        :disabled="!activeDocumentId || isLoading"
        @keyup.enter="sendMessage"
      />
      <button 
        @click="sendMessage"
        :disabled="!activeDocumentId || isLoading"
      >
      Send
      </button>
    </div>


  </div>
</template>

<script setup>
import { ref } from 'vue'

const isOpen = ref(true)
const userInput = ref('')
const isLoading = ref(false)

const messages = ref([
  { role: 'ai', content: 'Upload a resume and ask for feedback!' }
])

// selected documentID
const activeDocumentId = ref(null)
const activeDocumentName = ref('')

function setActiveDocument(doc) {
  activeDocumentId.value = doc.id
  activeDocumentName.value = doc.title
  isOpen.value = true
}


function closeChat() {
  isOpen.value = false
}

async function sendMessage() {
  const input = userInput.value.trim()
  if (!input || isLoading.value) return

  if (!activeDocumentId.value) {
    messages.value.push({
      role: 'ai',
      content: 'Please select a document to analyze first.'
    })
    return
  }

  messages.value.push({
    role: 'user',
    content: input
  })

  userInput.value = ''
  isLoading.value = true

  try {
    const res = await fetch('/api/ai/chat', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        message: input,
        documentIds: [activeDocumentId.value],
        history: messages.value
      })
    })

    if (!res.ok) throw new Error(`HTTP ${res.status}`)

    const data = await res.json()

    messages.value.push({
      role: 'ai',
      content: data.reply || 'No response from AI'
    })

  } catch (err) {
    console.error(err)

    messages.value.push({
      role: 'ai',
      content: 'Error: Unable to get response from server.'
    })
  } finally {
    isLoading.value = false
  }
}

</script>
<style scoped src="@/assets/css/chatbox.css"></style>