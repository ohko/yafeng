import { createApp } from 'vue'
import components from './components.js';
import { router } from './routes.js'
import './style.css'
import { createPinia } from 'pinia'
import App from './App.vue'
import Antd from 'ant-design-vue';
import 'ant-design-vue/dist/reset.css';
import * as antIcons from '@ant-design/icons-vue';
import { post, get } from '@/utils'
import { message } from 'ant-design-vue';

window.post = post
window.get = get
window.message = message

const pinia = createPinia()
const app = createApp(App)

app.config.globalProperties.$timeformat = (value) => { return new Date(value).toLocaleString().replaceAll('/', '-') };
app.use(pinia)
app.use(components)
app.use(router)
app.use(Antd)
Object.keys(antIcons).forEach(key => { app.component(key, antIcons[key]) })
app.mount('body')
