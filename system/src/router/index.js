import { createRouter, createWebHashHistory } from 'vue-router'
import StoryBoard from "../views/StoryBoard.vue";

const routes = [
  {
    path: "/",
    name: "Portal",
    component: StoryBoard,
  },
  {
    path: "/portal",
    name: "Portal",
    component: StoryBoard,
  },
  // {
  //   path: '/',
  //   name: 'Home',
  //   component: Home
  // },
  // {
  //   path: '/portal',
  //   name: 'Portal',
  //   // route level code-splitting
  //   // this generates a separate chunk (about.[hash].js) for this route
  //   // which is lazy-loaded when the route is visited.
  //   component: () => import(/* webpackChunkName: "about" */ '../views/About.vue')
  // }
];

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router