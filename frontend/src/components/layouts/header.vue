<template>
  <!-- Header Navigation Bar-->
    <!-- Main Header Wrapper (now sticky via site-header class) -->
    <header class="header">
        <nav>
            <!-- Logo -->
        <RouterLink to="/" class="nav-logo">
          <img src="/images/bwc.png" alt="BWC Logo">
        </RouterLink>

            <!-- Title -->
            <h1>Banana Waffle Cookies</h1>

            <!-- Navigation Links -->
            <div class="navigation" id="navbar">
              <RouterLink to="/dashboard" class="navlink">Dashboard</RouterLink>
              <RouterLink to="/library" class="navlink">Document Library</RouterLink>
              <RouterLink to="/profile" class="navlink">Profile</RouterLink>
              <RouterLink to="/settings" class="navlink">Settings</RouterLink>
            </div>

            <!-- Mobile Dropdown -->
            <div v-if="isOpen" ref="menuRef" class="hamburger-menu">
              <RouterLink to="/dashboard" class="moblink" @click="isOpen = false">Dashboard</RouterLink>
              <RouterLink to="/library" class="moblink" @click="isOpen = false">Library</RouterLink>
              <RouterLink to="/profile" class="moblink" @click="isOpen = false">Profile</RouterLink>
              <RouterLink to="/settings" class="moblink" @click="isOpen = false">Settings</RouterLink>
              <button @click.stop="signOut" class="signout-button mobdropdown">Sign Out</button>
            </div>

            <!-- Log Out Button -->
            <button @click.stop="signOut" class="signout-button desktop">Sign Out</button>
            <!-- Mobile Menu Hamburger -->
            <button @click.stop="toggleMenu" class="hamburger" aria-label="Menu" aria-expanded="false">☰</button>
        </nav>
      </header>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'

const isOpen = ref(false)
const menuRef = ref(null)

const signOut = () => {
  //insert logic to sign out here
  console.log("Signed Out!")
}

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
