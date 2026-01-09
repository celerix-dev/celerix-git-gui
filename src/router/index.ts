import { createRouter, createWebHistory } from 'vue-router';
import { allRoutes } from './routes';

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: allRoutes,
    linkActiveClass: 'active',
    linkExactActiveClass: 'active',
});

export default router;
