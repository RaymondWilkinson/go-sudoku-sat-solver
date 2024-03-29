import Vue from 'vue'
import App from './App.vue'

import './assets/styles/index.css';
import '@fortawesome/fontawesome-free/js/all.js';

Vue.config.productionTip = false

new Vue({
  render: h => h(App),
}).$mount('#app')
