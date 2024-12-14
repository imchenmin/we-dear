import { createRouter, createWebHistory } from 'vue-router'
import MainLayout from '@/layouts/MainLayout.vue'
import DoctorView from '@/views/DoctorView.vue'
import PatientChat from '@/views/PatientChat.vue'
import DepartmentView from '@/views/DepartmentView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      component: MainLayout,
      children: [
        {
          path: '',
          name: 'doctor',
          component: DoctorView
        },
        {
          path: 'patient-chat',
          name: 'patient-chat',
          component: PatientChat
        },
        {
          path: 'departments',
          name: 'departments',
          component: DepartmentView
        }
      ]
    }
  ]
})

export default router 