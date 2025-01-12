import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '@/views/LoginView.vue'
import PatientLoginView from '@/views/PatientLoginView.vue'
import DoctorView from '@/views/DoctorView.vue'
import PatientChat from '@/views/PatientChat.vue'
import { useUserStore } from '@/stores/user'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/patient/login'
    },
    {
      path: '/login',
      redirect: '/patient/login'
    },
    {
      path: '/doctor/login',
      name: 'doctorLogin',
      component: LoginView
    },
    {
      path: '/patient/login',
      name: 'patientLogin',
      component: PatientLoginView
    },
    {
      path: '/doctor',
      name: 'doctor',
      component: DoctorView,
      meta: { requiresAuth: true, role: 'doctor' }
    },
    {
      path: '/patient/chat',
      name: 'patientChat',
      component: PatientChat,
      meta: { requiresAuth: true, role: 'patient' }
    }
  ]
})

router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  console.log(`Navigating from ${from.path} to ${to.path}`)
  console.log('Current user role:', userStore.getRole())

  // 登录页面不需要验证
  if (to.path.endsWith('/login')) {
    console.log('Accessing login page, no authentication required')
    next()
    return
  }

  // 需要认证但未登录
  if (to.meta.requiresAuth && !userStore.isLoggedIn()) {
    console.log('Authentication required but not logged in')
    next(userStore.getRole() === 'doctor' ? '/doctor/login' : '/patient/login')
    return
  }

  // 角色不匹配
  if (to.meta.role && to.meta.role !== userStore.getRole()) {
    console.log(`Role mismatch: expected ${to.meta.role}, found ${userStore.getRole()}`)
    // 如果是医生访问患者页面，重定向到医生页面
    if (userStore.isDoctor && to.path.startsWith('/patient')) {
      console.log('Doctor trying to access patient page, redirecting to doctor page')
      next('/doctor')
      return
    }
    // 如果是患者访问医生页面，重定向到患者聊天页面
    if (userStore.isPatient && to.path.startsWith('/doctor')) {
      console.log('Patient trying to access doctor page, redirecting to patient chat page')
      next('/patient/chat')
      return
    }
  }

  console.log('Navigation allowed')
  next()
})

export default router 