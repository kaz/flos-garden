import Vue from 'vue'
import Router from 'vue-router'

import Node from './views/Node.vue'
import CnC from './views/CnC.vue'
import Monitor from './views/Monitor.vue'
import Audit from './views/Audit.vue'
import Backup from './views/Backup.vue'

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
    },
    {
      path: '/audit',
      component: Audit
    },
    {
      path: '/backup',
      component: Backup
    }
  ]
})
