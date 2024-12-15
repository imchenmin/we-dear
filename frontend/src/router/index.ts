import { createRouter, createWebHistory } from 'vue-router'
import MainLayout from '@/layouts/MainLayout.vue'
import DoctorView from '@/views/DoctorView.vue'
import PatientChat from '@/views/PatientChat.vue'
import DepartmentView from '@/views/DepartmentView.vue'
import LoginView from '@/views/LoginView.vue'
import { useUserStore } from '@/stores/user'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: LoginView
    },
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

router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  
  if (to.path === '/login') {
    if (userStore.token) {
      next('/')
    } else {
      next()
    }
    return
  }

  if (!userStore.token) {
    next('/login')
    return
  }

  next()
})

export default router 