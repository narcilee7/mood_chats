<template>
  <div class="login">
    <button @click="loginWithGitHub">使用 GitHub 登录</button>
  </div>
</template>

<script setup>
import { useRouter, useRoute } from 'vue-router'
import { ref, onMounted } from 'vue'

const router = useRouter()
const route = useRoute()

const loginWithGitHub = () => {
  window.location.href = 'http://localhost:8081/api/login' // 后端登录跳转入口
}

// 登录回调处理
onMounted(async () => {
  const token = route.query.token
  if (token) {
    localStorage.setItem('token', token)
    // 可以请求 /me 接口拿到用户信息
    const res = await fetch('http://localhost:8081/api/me', {
      headers: {
        Authorization: `Bearer ${token}`
      }
    })
    const user = await res.json()
    console.log('用户信息', user)
    // 处理登录状态，跳转页面
    router.replace('/')
  }
})
</script>
