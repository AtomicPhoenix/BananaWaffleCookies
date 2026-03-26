<template>
  <body>
    <!-- Header Navigation Bar-->
    <!-- Main Header Wrapper (now sticky via site-header class) -->
    <header class="header">
        <nav>
            <!-- Logo -->
            <a href="/" class="nav-logo">
                <img src="/public/images/bwc.png" alt="BWC Logo">
            </a>

            <!-- Title -->
            <h1>Banana Waffle Cookies</h1>

            <!-- Navigation Links -->
            <div class="navigation" id="navbar">
                <a href="#" class="navlink">Dashboard</a>
                <a href="#" class="navlink">Document Library</a>
                <a href="#" class="navlink">Profile</a>	
                <a href="#" class="navlink">Settings</a>
            </div>

            <!-- Mobile Dropdown -->
            <div v-if="isOpen" ref="menuRef" class="hamburger-menu">
              <a href="#" class="moblink">Dashboard</a>
              <a href="#" class="moblink">Library</a>
              <a href="#" class="moblink">Profile</a>
              <a href="#" class="moblink">Settings</a>
            </div>

            <!-- Mobile Menu Hamburger -->
            <button @click.stop="toggleMenu" class="hamburger" aria-label="Menu" aria-expanded="false">☰</button>
        </nav>
      </header>

    <!-- Page Content -->
    <main class="main-page">
      <slot/>
    </main>

    <!-- Footer -->
    <footer class="footer">
      <hr>
      <div class="footer-row">
        <!-- Footer Logo and Title -->
        <div>
          <div>
            <a href="/" class="footer-logo">
                <img src="/public/images/bwc.png" alt="BWC Logo">
            </a>
            <h1>Banana Waffle Cookies</h1>
          </div>
        </div>
        <!-- Footer Quick Connect Links -->
        <div class="footer-column">
          <span> Connect </span>
          <a href="https://github.com/AtomicPhoenix/BananaWaffleCookies" class="footer-info">Github</a> <!-- Private Link, placeholder -->
          <a href="#" class="footer-info">LinkedIn</a> <!-- No Linkedin -->
          <a href="https://discord.com/" class="footer-link">Discord</a> <!-- No Server Link -->
        </div>
        <div class="footer-column">
          <span class="footer-header">Connect Information</span>
          <span class="footer-info">support@bwc.org</span>
          <span class="footer-info">XXX-XXX-XXXX</span>
        </div>
      </div>
      <div class="copyright">
        © 2026 Banana Waffle Cookies
      </div>
    </footer>
  </body>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'

const isOpen = ref(false)
const menuRef = ref(null)

const toggleMenu = () => {
  isOpen.value = !isOpen.value
}

const handleClickOutside = (event) => {
  if (
    menuRef.value &&
    !menuRef.value.contains(event.target)
  ) {
    isOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>
