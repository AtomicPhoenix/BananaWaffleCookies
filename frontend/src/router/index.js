import { createRouter, createWebHistory } from "vue-router";
import axios from "axios";
import HomePage from "@/pages/home.vue";
import LibraryPage from "@/pages/library.vue";
import ProfilePage from "@/pages/profile.vue";
import SettingsPage from "@/pages/settings.vue";
import LoginPage from "@/pages/login.vue";
import DashboardPage from "@/pages/dashboard.vue";
import CreateJobPage from "@/pages/create-job.vue";
import ViewJobPage from "@/pages/view-job.vue";
import EditJobPage from "@/pages/edit-job.vue";
import SignupPage from "@/pages/signup.vue";

import JobDetail from "@/pages/job-detail.vue";

async function checkUserAuth(_) {
  try {
    const response = await axios.get("/api/auth", {
      withCredentials: true,
    });
    if (response.data.authenticated) {
      return true; // allow navigation
    } else {
      return { name: "login" }; // redirect to login
    }
  } catch (err) {
    console.log(err);
    return { name: "login" }; // redirect if API fails
  }
}

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "home",
      component: HomePage,
    },
    {
      path: "/library",
      name: "library",
      component: LibraryPage,
      beforeEnter: checkUserAuth,
    },
    {
      path: "/profile",
      name: "profile",
      component: ProfilePage,
      beforeEnter: checkUserAuth,
    },
    {
      path: "/settings",
      name: "settings",
      component: SettingsPage,
      beforeEnter: checkUserAuth,
    },
    {
      path: "/login",
      name: "login",
      component: LoginPage,
    },
    {
      path: "/dashboard",
      name: "dashboard",
      component: DashboardPage,
      beforeEnter: checkUserAuth,
    },
    {
      path: "/jobs",
      children: [
        { path: "create", name: "create-job", component: CreateJobPage },
        { path: ":job_id", component: ViewJobPage },
        { path: ":job_id/detail", name: "job-detail", component: JobDetail },
        { path: ":job_id/edit", name: "edit-job", component: EditJobPage },
      ],
      beforeEnter: checkUserAuth,
    },
    {
      path: "/signup",
      name: "signup",
      component: SignupPage,
    },
    {
      path: "/:pathMatch(.*)*",
      redirect: "/",
    },
  ],
});

export default router;
