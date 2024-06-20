import {createRouter, createWebHashHistory} from "vue-router";
import Home from './pages/Home.vue'
import Dashboard from "@/pages/Dashboard.vue";
import EntityPage from "@/pages/EntityPage.vue";
import RelationPage from "@/pages/RelationPage.vue";

const routes = [
    {
        path: '/',
        name: 'Home',
        component: Home,
    },
    {
        path: '/dashboard/:id',
        name: 'Dashboard',
        component: Dashboard,
    },
    {
        path: '/entity/:id',
        name: 'Entity',
        component: EntityPage,
    },
    {
        path: '/relation/:id',
        name: 'Relation',
        component: RelationPage,
    }
];

const router = createRouter({
    history: createWebHashHistory(),
    routes,
})

export default router