import { createRouter, createWebHistory } from "vue-router";


const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "home",
      component: () => import("../views/HomeView.vue"),
    },
    {
      path: "/about",
      name: "about",
      component: () => import("../views/AboutView.vue"),
    },
    {
      path: "/status",
      name: "status",
      component: () => import("../views/StatusView.vue"),
    },
    {
      path: "/chat",
      name: "chat",
      component: () => import("../views/ChatView.vue"),
    }
  ]
})

export default router;