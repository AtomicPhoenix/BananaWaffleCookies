<template>
  <div class="page-container">
    <h2 class="page-title">Settings</h2>

    <div class="form-card">
      
      <!-- Account Info -->
      <div class="section">
        <h3 class="section-title">Account Information</h3>

        <div class="form-group">
          <label>Email</label>
          <input v-model="form.email" type="email" />
        </div>
      </div>

      <!-- Password Reset -->
      <div class="section">
        <h3 class="section-title">Password</h3>

        <div class="form-group">
          <label>New Password</label>
          <input
            v-model="form.password"
            type="password"
            placeholder="Enter new password"
          />
        </div>

        <div class="form-group">
          <label>Confirm Password</label>
          <input
            v-model="form.confirmPassword"
            type="password"
            placeholder="Confirm new password"
          />
        </div>
      </div>

      <!-- Preferences -->
      <div class="section">
        <h3 class="section-title">Preferences</h3>

        <label class="checkbox-group">
          <input type="checkbox" v-model="form.notifications" />
          Enable notifications
        </label>
      </div>

      <!-- Save Button -->
      <button @click="saveSettings">
        Save Changes
      </button>

      <!-- Messages -->
      <p v-if="error" class="error">{{ error }}</p>
      <p v-if="saved" class="success">
        Settings saved successfully!
      </p>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'

const form = reactive({
  email: '',
  password: '',
  confirmPassword: '',
  notifications: false
})

const saved = ref(false)
const error = ref('')

async function saveSettings() {
  error.value = ''

  // Validation
  if (form.password && form.password !== form.confirmPassword) {
    error.value = 'Passwords do not match'
    return
  }

  try {
    const res = await fetch(`/api/settings`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email: form.email, password: form.password })
    })

  if (res.ok) {
      form.confirmPassword = ''
    }
  } catch (err) {
    console.error(err)
    return
  }

  console.log('Saved settings:', form)
  saved.value = true

  setTimeout(() => {
    saved.value = false
  }, 2000)
}

</script>

<style scoped src="@/assets/css/settings.css"></style>
