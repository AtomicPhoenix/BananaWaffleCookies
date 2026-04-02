<template>
    <div class="home-content">
      <div class="title-row">
        <RouterLink to="/" class="title-logo-link" aria-label="Go to home page">
          <img class="title-logo" src="/images/bwc.png" alt="BWC Logo">
        </RouterLink>
        <h1>Banana Waffle<br>Cookies</h1>
      </div>
        <p>Keep applications organized, generate tailored resumes, and manage statuses in one place!</p>
      <div>
        <RouterLink to="/signup" class="button">Get Started!</RouterLink>
      </div>
    </div>
</template>

<style scoped src="@/assets/css/home.css"></style>
