<template>
    <div class="signup-page">
      <h1>Sign Up</h1>
      <!-- Form wrapper that prevents default page reload on submit -->
      <form @submit.prevent="handleSignup" class="signup-form">

      <!-- Email field container -->
      <div class="input-field">
        <!-- Label describing the input -->
        <label for="email">Email: </label>

        <!-- Input bound to reactive state using v-model -->
        <input
          id="email"
          type="text"
          v-model="form.email"
          required
          class="input-text"
        />
      </div>

      <!-- Password field container -->
      <div class="input-field">
        <!-- Label describing the input -->
        <label for="password">Password: </label>

        <!-- Password input (masked) bound to reactive state -->
        <input
          id="password"
          type="password"
          v-model="form.password"
          required
          class="input-text"
        />
      </div>

      <!-- Submit button to trigger login logic -->
      <div class="submit">
        <button class="register-button" type="submit">Register</button>
      </div>

      <!-- Login with OAuth -->
      <div class="oauth">
        <!-- OAuth Logic Goes Here -->
        <router-link class="oauth-button" to="/">
          Sign up with Github <img src="/images/github.png" class="github-logo" alt="Github">
        </router-link>
      </div>

    </form>
  </div>
</template>

<script setup>
// Import reactive utility for form state //
import { reactive } from 'vue'

// Reactive object to store form input values //
const form = reactive({
  email: '',
  password: ''
})

// Function that runs when the form is submitted //
async function handleSignup() {
  try {
    const res = await fetch(`/api/signup`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email: form.email, password: form.password })
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

<style scoped src="@/assets/css/signup.css"></style>
