<template>
  <div class="app-shell">
    <HeaderLayout />
    <RouterView />
    <FooterLayout />
  </div>
</template>

<script setup>
import { RouterView } from 'vue-router'
import HeaderLayout from '@/components/layouts/header.vue'
import FooterLayout from '@/components/layouts/footer.vue'
</script>