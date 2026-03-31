<template>
    <div class="dashboard">
      <!-- Search Bar to Query Job Information -->
      <div class="search-box">
        <form action="backend.go" method="get">
            <label class="search-label" for="job-search">Search Jobs  </label>
            <input class="search-bar" type="search" id="site-search" name="query" placeholder="Enter Job Information Here">
            <button class="search-unicode" type="submit">⌕</button>
        </form>
      </div>
      <div class="job-info">
        <!-- Pie Chart Overview of Job Information, Later Sprint Implementation -->
        <div class="overview">
            <h1 class="overview-title"> Overview </h1>
            <p class="overview-items"> 
                Interested: 34<br>
                Applied: 13<br>
                Interview Offered: 2<br>
                Accepted: 1<br>
                Rejected: 1<br>
            </p>
            <RouterLink to="/" class="pie-template">
            <img src="/images/pie-template.png" alt="pie">
            </RouterLink>
        </div>
        <!-- Job List -->
        <div class="job-list">
                <!-- Field to add jobs to listing-->
                <div class="create-job">
                  <router-link class="create-job-button" to="/create-job">
                    Create New Job Application
                  </router-link>
                </div>
                <!-- Needs to be created of internal database information, template for now -->
                <div class="job-listing">
                    <div class="left top">[Job Title]   |   [Company]   |   [Location]</div>
                    <div class="left mid jdesc">[Last Modified: ] </div>
                    <div class="left bot jdesc">[Deadline: ]</div>

                    <div class="listing-status-button right mid" id="accepted">Accepted </div>
                </div>

                <div class="job-listing">
                    <div class="left top">[Job Title]   |   [Company]   |   [Location]</div>
                    <div class="left mid jdesc">[Last Modified: ] </div>
                    <div class="left bot jdesc">[Deadline: ]</div>

                    <div class="listing-status-button right mid" id="interview">Interview </div>
                </div>

                <div class="job-listing">
                    <div class="left top">[Job Title]   |   [Company]   |   [Location]</div>
                    <div class="left mid jdesc">[Last Modified: ] </div>
                    <div class="left bot jdesc">[Deadline: ]</div>

                    <div class="listing-status-button right mid" id="applied">Applied </div>
                </div>
        </div>
      </div>
    </div>
</template>
<style scoped src="@/assets/css/dashboard.css"></style>
<script setup>
import { ref } from 'vue'

const username = ref('')
const onSubmit = () => console.log(username.value)
</script>