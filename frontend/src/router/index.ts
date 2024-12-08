import { createRouter, createWebHistory } from 'vue-router'
import DoctorView from '../views/DoctorView.vue'
import PatientChat from '../views/PatientChat.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'doctor',
      component: DoctorView
    },
    {
      path: '/patient-chat',
      name: 'patient-chat',
      component: PatientChat
    }
  ]
})

export default router 