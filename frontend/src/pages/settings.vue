<template>
  <div>
    <h2 class="text-2xl font-bold mb-6">Settings</h2>

    <div class="bg-white p-6 rounded shadow max-w-lg space-y-6">
      
      <!-- Account Info -->
      <div>
        <h3 class="font-semibold mb-3">Account Information</h3>

        <div class="mb-3">
          <label class="block mb-1">First Name</label>
          <input v-model="form.firstName" class="w-full border p-2 rounded" />
        </div>

        <div class="mb-3">
          <label class="block mb-1">Last Name</label>
          <input v-model="form.lastName" class="w-full border p-2 rounded" />
        </div>

        <div>
          <label class="block mb-1">Email</label>
          <input v-model="form.email" type="email" class="w-full border p-2 rounded" />
        </div>
      </div>

      <!-- Password Reset -->
      <div>
        <h3 class="font-semibold mb-3">Password</h3>

        <div class="mb-3">
          <label class="block mb-1">New Password</label>
          <input
            v-model="form.password"
            type="password"
            class="w-full border p-2 rounded"
            placeholder="Enter new password"
          />
        </div>

        <div>
          <label class="block mb-1">Confirm Password</label>
          <input
            v-model="form.confirmPassword"
            type="password"
            class="w-full border p-2 rounded"
            placeholder="Confirm new password"
          />
        </div>
      </div>

      <!-- Preferences -->
      <div>
        <h3 class="font-semibold mb-2">Preferences</h3>
        <label class="flex items-center gap-2">
          <input type="checkbox" v-model="form.notifications" />
          Enable notifications
        </label>
      </div>

      <!-- Save Button -->
      <button
        @click="saveSettings"
        class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600"
      >
        Save Changes
      </button>

      <!-- Messages -->
      <p v-if="error" class="text-red-500 text-sm">{{ error }}</p>
      <p v-if="saved" class="text-green-600 text-sm">
        Settings saved successfully!
      </p>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'

const form = reactive({
  firstName: '',
  lastName: '',
  email: '',
  password: '',
  confirmPassword: '',
  notifications: false
})

const saved = ref(false)
const error = ref('')

function saveSettings() {
  error.value = ''

  // Simple validation for password match
  if (form.password && form.password !== form.confirmPassword) {
    error.value = 'Passwords do not match'
    return
  }

  console.log('Saved settings:', form)

  saved.value = true

  setTimeout(() => {
    saved.value = false
  }, 2000)
}
</script>