<template>
  <div>
    <h2 class="text-2xl font-bold mb-6">Profile</h2>

    <!-- Completion Indicator -->
    <div class="mb-6">
      <p class="mb-2 font-semibold">
        Profile Completion: {{ completionPercentage }}%
      </p>

      <div class="w-full bg-gray-200 rounded h-4">
        <div
          class="bg-green-500 h-4 rounded"
          :style="{ width: completionPercentage + '%' }"
        ></div>
      </div>
    </div>

    <!-- Profile Form -->
    <div class="bg-white p-6 rounded shadow max-w-lg space-y-4">
      <div>
        <label class="block mb-1">First Name</label>
        <input v-model="form.firstName" class="w-full border p-2 rounded" />
      </div>

      <div>
        <label class="block mb-1">Last Name</label>
        <input v-model="form.lastName" class="w-full border p-2 rounded" />
      </div>

      <div>
        <label class="block mb-1">Email</label>
        <input v-model="form.email" type="email" class="w-full border p-2 rounded" />
      </div>

      <div>
        <label class="block mb-1">Summary</label>
        <textarea v-model="form.summary" class="w-full border p-2 rounded"></textarea>
      </div>

      <button
        @click="saveProfile"
        class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600"
      >
        Save Profile
      </button>
    </div>
  </div>
</template>

<script setup>
import { reactive, computed } from 'vue'

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

function saveProfile() {
  console.log('Saved profile:', form)
}
</script>