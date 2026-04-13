<template>
  <div class="page-container">
    <h2 class="page-title">Profile</h2>

    <!-- Completion Indicator -->
    <div class="completion-section">
      <div class="completion-header">
        <span>Profile Completion</span>
        <span>{{ completionPercentage }}%</span>
      </div>

      <div class="progress-bar">
        <div
          class="progress-fill"
          :style="{ width: completionPercentage + '%' }"
        ></div>
      </div>
    </div>

    <!-- Profile Form -->
    <div class="form-card">

      <!-- BASIC INFO -->
      <div class="section">
        <h3 class="section-title">Basic Information</h3>

        <div class="form-row">
          <div class="form-group">
            <label>First Name</label>
            <input v-model="form.first_name" />
          </div>

          <div class="form-group">
            <label>Last Name</label>
            <input v-model="form.last_name" />
          </div>
        </div>

        <div class="form-group">
          <label>Phone</label>
          <input v-model="form.phone" type="tel" />
        </div>
      </div>

      <!-- LOCATION -->
      <div class="section">
        <h3 class="section-title">Location</h3>

        <div class="form-row">
          <div class="form-group">
            <label>City</label>
            <input v-model="form.city" />
          </div>

          <div class="form-group">
            <label>State</label>
            <input v-model="form.state" />
          </div>
        </div>

        <div class="form-group">
          <label>Country</label>
          <input v-model="form.country" />
        </div>
      </div>

      <!-- LINKS -->
      <div class="section">
        <h3 class="section-title">Links</h3>

        <div class="form-group">
          <label>LinkedIn URL</label>
          <input v-model="form.linkedin_url" />
        </div>

        <div class="form-group">
          <label>Portfolio URL</label>
          <input v-model="form.portfolio_url" />
        </div>
      </div>

      <!-- SUMMARY -->
      <div class="section">
        <h3 class="section-title">Summary</h3>

        <div class="form-group">
          <textarea v-model="form.summary" rows="4"></textarea>
        </div>
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

  } catch (err) {
    console.error(err)
  } 
}
</script>

<style scoped src="@/assets/css/profile.css"></style>
