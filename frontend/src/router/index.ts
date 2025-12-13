import { createRouter, createWebHistory } from "vue-router"
import AuthView from "@/views/AuthView.vue"
import ProjectListView from "@/views/ProjectListView.vue"
import ProjectView from "@/views/ProjectView.vue"
import TemplateView from "@/views/TemplateView.vue"

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
    {
      path: "/projects/:id",
      name: "project",
      props: (route) => ({ id: Number(route.params.id) }),
      component: ProjectView,
    },
    {
      path: "/templates/:id",
      name: "template",
      props: (route) => ({ id: Number(route.params.id) }),
      component: TemplateView,
    },
  ],
})

export default router
