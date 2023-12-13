import {createApp} from 'vue'
import App from './App.vue'
import './style.css';
import VueCropper from 'vue-cropper'; 
import 'vue-cropper/dist/index.css'
import ElementPlus from "element-plus"
import 'element-plus/dist/index.css'


createApp(App).use(VueCropper).use(ElementPlus).mount('#app')
