/**
 * 应用入口文件
 * 初始化Vue应用实例，注册全局插件并挂载到DOM
 */
import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import './styles/main.css'

// 创建Vue应用实例
const app = createApp(App)

// 注册路由插件
app.use(router)

// 将应用挂载到index.html中id为app的DOM元素
app.mount('#app')
