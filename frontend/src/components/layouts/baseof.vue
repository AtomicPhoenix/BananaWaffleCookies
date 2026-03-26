<template>
  <div class="base">
    <!-- Header Navigation Bar-->
    <!-- Main Header Wrapper (now sticky via site-header class) -->
    <div class="header site-header">
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
            <div v-if="isOpen" ref="menuRef" class="hamburger">
              <a href="#" class="navlink">Dashboard</a>
              <a href="#" class="navlink">Library</a>
              <a href="#" class="navlink">Profile</a>
              <a href="#" class="navlink">Settings</a>
            </div>
            <!-- Mobile Menu Hamburger -->
            <button @click.stop="toggleMenu" class="hamburger" aria-label="Menu" aria-expanded="false">☰</button>
        </nav>
    </div>

    <!-- Page Content -->
    <main class="main-page">
      <slot/>
    </main>

    <!-- Footer -->
    <footer class="bg-gray-100 text-center p-4">
      © 2026 Banana Waffle Cookies
    </footer>

  </div>
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
