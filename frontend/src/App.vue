<template>
  <el-config-provider>
    <router-view></router-view>
  </el-config-provider>
</template>

<style>
html, body {
  margin: 0;
  padding: 0;
  height: 100%;
}

#app {
  height: 100%;
}
</style>

<script setup>
import { useUserStore } from '@/stores/user'
import { WebSocketService } from '@/utils/websocket'
import { onMounted, onUnmounted } from 'vue'

const userStore = useUserStore()
const wsService = WebSocketService.getInstance()

onMounted(() => {
  if (userStore.isLoggedIn) {
    wsService.connect()
  }
})

onUnmounted(() => {
  wsService.disconnect()
})
</script>
