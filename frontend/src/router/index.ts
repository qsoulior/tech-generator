import { createRouter, createWebHistory } from "vue-router"
import { useAuthStore } from "@/stores/auth"

declare module "vue-router" {
  interface RouteMeta {
    requiresAuth?: boolean
  }
}

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "auth",
      component: () => import("@/views/AuthView.vue"),
    },
    {
      path: "/projects",
      name: "projectList",
      component: () => import("@/views/ProjectListView.vue"),
      meta: { requiresAuth: true },
    },
    {
      path: "/projects/:projectID",
      name: "project",
      props: (route) => ({ projectID: Number(route.params.projectID) }),
      component: () => import("@/views/ProjectView.vue"),
      meta: { requiresAuth: true },
    },
    {
      path: "/projects/:projectID/template/:templateID",
      name: "template",
      props: (route) => ({
        projectID: Number(route.params.projectID),
        templateID: Number(route.params.templateID),
      }),
      component: () => import("@/views/TemplateView.vue"),
      meta: { requiresAuth: true },
    },
    {
      path: "/projects/:projectID/template/:templateID/tasks",
      name: "taskList",
      props: (route) => ({
        projectID: Number(route.params.projectID),
        templateID: Number(route.params.templateID),
      }),
      component: () => import("@/views/TaskListView.vue"),
      meta: { requiresAuth: true },
    },
    {
      path: "/projects/:projectID/template/:templateID/tasks/:taskID",
      name: "task",
      props: (route) => ({
        projectID: Number(route.params.projectID),
        templateID: Number(route.params.templateID),
        taskID: Number(route.params.taskID),
      }),
      component: () => import("@/views/TaskView.vue"),
      meta: { requiresAuth: true },
    },
  ],
})

router.beforeEach(async (to) => {
  if (!to.meta.requiresAuth) return true
  const authStore = useAuthStore()
  try {
    await authStore.ensureLoaded()
    return true
  } catch {
    return { name: "auth" }
  }
})

export default router
