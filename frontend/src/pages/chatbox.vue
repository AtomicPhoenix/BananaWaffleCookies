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
    </div>

    <!-- INPUT -->
    <div class="chatbox-input">
      <input
        v-model="userInput"
        type="text"
        placeholder="Ask for feedback on your documents..."
        @keyup.enter="sendMessage"
      />
      <button @click="sendMessage">Send</button>
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

// thought: pass selected docs as props, as of right now it reads all docs on the page
const selectedDocumentIds = ref([])

function closeChat() {
  isOpen.value = false
}

async function sendMessage() {
  const input = userInput.value.trim()
  if (!input || isLoading.value) return

  // push user message immediately
  messages.value.push({
    role: 'user',
    content: input
  })

  userInput.value = ''
  isLoading.value = true

  try {
    const res = await fetch('/api/ai/chat', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        message: input,
        documentIds: selectedDocumentIds.value, 
        history: messages.value
      })
    })

    if (!res.ok) {
      throw new Error(`HTTP ${res.status}`)
    }

    const data = await res.json()

    messages.value.push({
      role: 'ai',
      content: data.reply || 'No response from AI'
    })

  } catch (err) {
    console.error('Chat error:', err)

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