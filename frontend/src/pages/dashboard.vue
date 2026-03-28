<template>
  <main>
    <body>
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
                Interested: <br>
                Applied: <br>
                Interview Offered: <br>
                Accepted: <br>
                Rejected: <br>
            </p>
            <RouterLink to="/" class="pie-template">
            <img src="/images/pie-template.png" alt="pie">
            </RouterLink>
        </div>
        <!-- Job List -->
        <div class="job-list">
                <!-- Needs to be created of internal database information, template for now -->
                <div class="job-listing">
                    <div class="left top">[Job Title]   |   [Company]   |   [Location]</div>
                    <div class="left mid jdesc">[Last Modified: ] </div>
                    <div class="left bot jdesc">[Deadline: ]</div>

                    <div class="listing-status-button right mid">Accepted </div>
                </div>
        </div>
      </div>
    </body>
  </main>
</template>

<style scoped src="@/assets/css/dashboard.css"></style>