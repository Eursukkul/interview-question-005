import { createRouter, createWebHistory } from 'vue-router'
import TicketReceiveView from '../views/TicketReceiveView.vue'
import QueueDisplayView from '../views/QueueDisplayView.vue'
import QueueResetView from '../views/QueueResetView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', name: 'ticket-receive', component: TicketReceiveView },
    { path: '/queue', name: 'queue-display', component: QueueDisplayView },
    { path: '/reset', name: 'queue-reset', component: QueueResetView }
  ]
})

export default router
