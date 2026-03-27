<template>
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
