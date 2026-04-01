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
        <input v-model="form.first_name" />
      </div>

      <div class="form-group">
        <label>Last Name</label>
        <input v-model="form.last_name" />
      </div>

      <div class="form-group">
        <label>Phone</label>
        <input v-model="form.phone" type="tel" />
      </div>


      <div class="form-group">
        <label>Phone</label>
        <input v-model="form.city" type="text" />
      </div>


      <div class="form-group">
        <label>Phone</label>
        <input v-model="form.state" type="text" />
      </div>


      <div class="form-group">
        <label>Phone</label>
        <input v-model="form.country" type="text" />
      </div>

      <div class="form-group">
        <label>Phone</label>
        <input v-model="form.linkedin_url" type="text" />
      </div>

      <div class="form-group">
        <label>Phone</label>
        <input v-model="form.portfolio_url" type="text" />
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
  first_name: '',
  last_name: '',
  phone: '',
  city: '',
  state: '',
  country: '',
  linkedin_url: '',
  portfolio_url: '',
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

  getProfile()
  const saved = localStorage.getItem('profile')
  if (saved) {
    Object.assign(form, JSON.parse(saved))
  }
})

async function getProfile() {
  try {
      const res = await fetch(`/api/profile`, {method: 'GET'})
      if (res.ok) {
        let profile_data = await res.json()
        form.first_name = profile_data.first_name
        form.last_name = profile_data.last_name
        form.phone = profile_data.phone
        form.city = profile_data.city
        form.state = profile_data.state
        form.country = profile_data.country
        form.linkedin_url = profile_data.linkedin_url
        form.portfolio_url  = profile_data.portfolio_url
        form.summary = profile_data.summary
      }   
  } catch (err) {
    console.error(err)
  }
}

// Send profile data to backend API
async function saveProfile() {
  try {
    const res = await fetch(`/api/profile`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ first_name: form.first_name, 
                             last_name: form.last_name,
                             phone: form.phone,
                             city: form.city,
                             state: form.state,
                             country: form.country,
                             linkedin_url: form.linkedin_url,
                             portfolio_url : form.portfolio_url,
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
