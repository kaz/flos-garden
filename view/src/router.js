import Vue from 'vue'
import Router from 'vue-router'

import Node from './views/Node.vue'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/node',
      component: Node
    }
  ]
})
