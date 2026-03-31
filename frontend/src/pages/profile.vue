<template>
  <div class="page-container">
    <h2 class="page-title">Profile</h2>

    <!-- Completion Indicator -->
    <div class="completion-section">
      <p class="completion-text">
        Profile Completion: {{ completionPercentage }}%
      </p>

      <div class="progress-bar">
        <div
          class="progress-fill"
          :style="{ width: completionPercentage + '%' }"
        ></div>
      </div>
    </div>

    <!-- Profile Form -->
    <div class="form-card">
      <div class="form-group">
        <label>First Name</label>
        <input v-model="form.firstName" />
      </div>

      <div class="form-group">
        <label>Last Name</label>
        <input v-model="form.lastName" />
      </div>

      <div class="form-group">
        <label>Email</label>
        <input v-model="form.email" type="email" />
      </div>

      <div class="form-group">
        <label>Summary</label>
        <textarea v-model="form.summary"></textarea>
      </div>

      <button @click="saveProfile">
        Save Profile
      </button>
    </div>
  </div>
</template>

<script setup>
import { reactive, computed, onMounted } from 'vue'

const form = reactive({
  firstName: '',
  lastName: '',
  email: '',
  summary: ''
})

// Completion calculation
const completionPercentage = computed(() => {
  const fields = Object.values(form)
  const filled = fields.filter(value => value && value.trim() !== '').length
  return Math.round((filled / fields.length) * 100)
})

// Load saved profile (simulated persistence)
onMounted(() => {
  const saved = localStorage.getItem('profile')
  if (saved) {
    Object.assign(form, JSON.parse(saved))
  }
})


// Send profile data to backend API
async function saveProfile() {
  try {
    const res = await fetch(`/api/profile`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ firstName: form.firstName, 
                             lastName: form.lastName,
                             email: form.email 
                             summary: form.summary })
    })

    if (res.ok) {
      form.email = ''
      form.password = ''
    }   
  } catch (err) {
    console.error(err)
  } 
}
</script>

<style scoped src="@/assets/css/profile.css"></style>
