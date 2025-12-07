import { createRouter, createWebHistory } from "vue-router"
import AuthView from "@/views/AuthView.vue"
import ProjectListView from "@/views/ProjectListView.vue"

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "auth",
      component: AuthView,
    },
    {
      path: "/projects",
      name: "projectList",
      component: ProjectListView,
    },
  ],
})

export default router
