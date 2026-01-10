import 'halfmoon/css/halfmoon.min.css';
import 'halfmoon/css/cores/halfmoon.modern.css';
import 'bootstrap/dist/js/bootstrap.bundle.min.js';
import '@tabler/icons-webfont/dist/tabler-icons.min.css';
import './assets/main.css';

import '@/services/color-scheme.ts';

import { createApp } from 'vue';
import { createPinia } from 'pinia';

import App from './App.vue';
import router from './router/index';

const app = createApp(App);

app.use(createPinia());
app.use(router);

app.mount('#app');