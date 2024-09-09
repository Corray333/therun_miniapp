import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import i18n from './i18n'
import Vue3Lottie from 'vue3-lottie'
import PrimeVue from 'primevue/config'
import Aura from '@primevue/themes/aura'

const app = createApp(App)



app.use(createPinia())
app.use(router)
app.use(i18n)
app.use(Vue3Lottie)
app.use(PrimeVue, {
    theme: {
        preset: Aura,
        options:{
            darkModeSelector: '#app'
        }
    }
});

app.mount('#app')
