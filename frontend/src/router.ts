import {createRouter, createWebHashHistory} from "vue-router";
import Home from './pages/Home.vue'

const routes = [
    {
        path: '/',
        name: 'Home',
        component: Home,
    },
    // {
    //     path: '/createproject',
    //     name: 'Create new project',
    //     component: CreateProject,
    // },
]

const router = createRouter({
    history: createWebHashHistory(),
    routes,
})

export default router