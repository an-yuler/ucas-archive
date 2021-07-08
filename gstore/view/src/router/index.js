import Vue from 'vue'
import VueRouter from 'vue-router'
import Home from '../views/Home.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    // redirect: '/query1',
    name: 'Home',
    component: Home
  },
  {
    path: '/query1',
    name: 'Query1',
    component: () => import('../views/Query1.vue')
  },
  {
    path: '/query2',
    name: 'Query2',
    component: () => import('../views/Query2.vue')
  },
  {
    path: '/query3',
    name: 'Query3',
    component: () => import('../views/Query3.vue')
  }
]

const router = new VueRouter({
  mode: 'history',
  routes
})

export default router
