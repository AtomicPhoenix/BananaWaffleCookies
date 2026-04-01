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

<!--
type Profile struct {
	Id                int       `json:"id"`
	UserID            int       `json:"user_id"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	Phone             string    `json:"phone"`
	City              string    `json:"city"`
	State             string    `json:"state"`
	Country           string    `json:"country"`
	LinkedinURL       string    `json:"linkedin_url"`
	PortfolioURL      string    `json:"portfolio_url"`
	Summary           string    `json:"summary"`
	CompletionPercent int       `json:"completion_percent"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
-->
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
