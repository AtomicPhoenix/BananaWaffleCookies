<template>
    <div class="login-page">
      <h1>Login</h1>
      <!-- Form wrapper that prevents default page reload on submit -->
      <form @submit.prevent="handleLogin">

      <!-- Username field container -->
      <div class="input-field">
        <!-- Label describing the input -->
        <label for="username">Username</label>

        <!-- Input bound to reactive state using v-model -->
        <input
          id="username"
          type="text"
          v-model="form.username"
          required
        />
      </div>

      <!-- Password field container -->
      <div class="input-field">
        <!-- Label describing the input -->
        <label for="password">Password</label>

        <!-- Password input (masked) bound to reactive state -->
        <input
          id="password"
          type="password"
          v-model="form.password"
          required
        />
      </div>

      <!-- Login with OAuth -->
      <div class="oauth-button"> 
        <button>IDK</button>
      </div>
      <!-- Submit button to trigger login logic -->
      <div>
        <button type="submit">Login</button>
      </div>

    </form>
  </div>
</template>

<script setup>
// Import reactive utility for form state //
import { reactive } from 'vue'

// Reactive object to store form input values //
const form = reactive({
  username: '',
  password: ''
})

// Function that runs when the form is submitted //
function handleLogin() {
  // Place Authentication Logic here //
  

  // Example: clear form after submission //
  form.username = ''
  form.password = ''
}
</script>

<style scoped src="@/assets/css/login.css"></style>