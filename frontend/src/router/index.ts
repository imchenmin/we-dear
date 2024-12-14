import { createRouter, createWebHistory } from 'vue-router'
import DoctorView from '../views/DoctorView.vue'
import PatientChat from '../views/PatientChat.vue'
import DoctorForm from '../components/DoctorForm.vue'
import PatientForm from '../components/PatientForm.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
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
    },
    {
      path: '/doctors/new',
      name: 'newDoctor',
      component: DoctorForm
    },
    {
      path: '/patients/new',
      name: 'newPatient',
      component: PatientForm
    }
  ]
})

export default router 