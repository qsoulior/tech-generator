import { createRouter, createWebHistory } from "vue-router"
import AuthView from "@/views/AuthView.vue"
import ProjectListView from "@/views/ProjectListView.vue"
import ProjectView from "@/views/ProjectView.vue"
import TemplateView from "@/views/TemplateView.vue"
import TaskListView from "@/views/TaskListView.vue"

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
      path: "/projects/:projectID",
      name: "project",
      props: (route) => ({ projectID: Number(route.params.projectID) }),
      component: ProjectView,
    },
    {
      path: "/projects/:projectID/template/:templateID",
      name: "template",
      props: (route) => ({
        projectID: Number(route.params.projectID),
        templateID: Number(route.params.templateID),
      }),
      component: TemplateView,
    },
    {
      path: "/projects/:projectID/template/:templateID/tasks",
      name: "taskList",
      props: (route) => ({
        projectID: Number(route.params.projectID),
        templateID: Number(route.params.templateID),
      }),
      component: TaskListView,
    },
  ],
})

export default router
