import {createRouter, createWebHashHistory} from "vue-router";
import Home from './pages/Home.vue'
import Dashboard from "@/pages/Dashboard.vue";
import EntityPage from "@/pages/EntityPage.vue";
import RelationPage from "@/pages/RelationPage.vue";
import EntityCollectionPage from "@/pages/EntityCollectionPage.vue";

const routes = [
    {
        path: '/',
        name: 'Home',
        component: Home,
    },
    {
        path: '/dashboard',
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
    },
    {
        path: '/entities',
        name: 'Entities',
        component: EntityCollectionPage,
    }
];

const router = createRouter({
    history: createWebHashHistory(),
    routes,
})

export default router