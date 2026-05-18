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
    {
      path: "/:pathMatch(.*)*",
      name: "notFound",
      component: () => import("@/views/NotFoundView.vue"),
    },
  ],
})

// On initial navigation, a guard redirect sometimes doesn't reach history.replaceState,
// so the address bar keeps showing the originally requested URL. Reconcile it after every nav.
router.afterEach((to) => {
  const current = window.location.pathname + window.location.search + window.location.hash
  if (current !== to.fullPath) {
    window.history.replaceState(window.history.state, "", to.fullPath)
  }
})

router.beforeEach(async (to) => {
  const authStore = useAuthStore()

  if (to.name === "auth") {
    try {
      await authStore.ensureLoaded()
      const redirect = to.query.redirect
      return typeof redirect === "string" && redirect.startsWith("/") ? redirect : { name: "projectList" }
    } catch {
      return true
    }
  }

  if (!to.meta.requiresAuth) return true

  try {
    await authStore.ensureLoaded()
    return true
  } catch {
    return { name: "auth", query: { redirect: to.fullPath } }
  }
})

export default router
