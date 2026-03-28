<template>
  <main>
    <body>
      <!-- Search Bar to Query Job Information -->
      <div class="search-box">
        <form action="backend.go" method="get">
            <label for="job-search">Search Jobs</label>
            <input type="search" id="site-search" name="query" placeholder="Search Jobs">
            <button type="submit">Go</button>
        </form>
      </div>
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
        <RouterLink to="/" class="pie-templete">
          <img src="/images/pie-template.png" alt="pie">
        </RouterLink>
      </div>
      <!-- Job List -->
       <div class="job-list">
            <!-- Needs to be created of internal database information, template for now -->
            <div class="job-listing">
                <h1> Job Title 1 </h1>
                <p> Last Modified:  </p>
                <p> Deadline: </p>
            </div>
            <h1 class="listing-status-button"> Accepted </h1>
       </div>
    </body>
  </main>
</template>

<style scoped src="@/assets/css/dashboard.css"></style>