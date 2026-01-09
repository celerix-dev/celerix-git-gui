import type { RouteRecordRaw } from 'vue-router';
import HomeView from '@/views/HomeView.vue';

import GitGuiView from "@/views/GitGuiView.vue";

// Define the routes for the generic routes
export const basicRoutes: RouteRecordRaw[] = [
    {
        path: '/',
        name: 'home',
        component: HomeView,
        children: [
            {
                path: '/',
                name: 'git-gui-index',
                component: GitGuiView
            },
        ]
    }
];
