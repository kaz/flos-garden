import Vue from 'vue'
import Router from 'vue-router'

import Node from './views/Node.vue'
import CnC from './views/CnC.vue'
import Monitor from './views/Monitor.vue'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/node',
      component: Node
    },
    {
      path: '/cnc',
      component: CnC
    },
    {
      path: '/monitor',
      component: Monitor
    }
  ]
})
