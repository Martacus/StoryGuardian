import { createRouter, createWebHashHistory } from 'vue-router'
import StoryBoard from "../views/story/StoryBoard.vue";
import AddStory from "../views/story/AddStory.vue";

const routes = [
  {
    path: "/",
    name: "Portal",
    component: StoryBoard,
  },
  {
    path: "/portal",
    name: "Portal",
    component: StoryBoard
  },
  {
    path: '/portal/story/add',
    name: 'AddStory',
    component: AddStory
  }
];

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router