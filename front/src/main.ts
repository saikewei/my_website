import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import 'element-plus/dist/index.css' // 确保这行存在

const app = createApp(App)

app.use(createPinia())
app.use(router)

app.mount('#app')
