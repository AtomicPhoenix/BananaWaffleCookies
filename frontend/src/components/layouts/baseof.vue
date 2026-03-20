<template>
  <div class="base">
    <!-- Header Navigation Bar-->
    <!-- Main Header Wrapper (now sticky via site-header class) -->
    <div class="header site-header">
        <nav>
            <!-- Logo -->
            <a href="/" class="nav-logo">
                <img src="/public/images/temp-logo.png" alt="BWC Logo">
            </a>

            <!-- Title -->
            <h1>Banana Waffle Cookies</h1>

            <!-- Navigation Links -->
            <div class="navigation" id="navbar">
                <a href="/" class="navlink">Dashboard</a>
                <a href="/" class="navlink">Document Library</a>
                <a href="/" class="navlink">Profile</a>	
                <a href="/" class="navlink">Settings</a>
            </div>
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
