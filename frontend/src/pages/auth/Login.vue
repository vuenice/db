<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute, RouterLink } from 'vue-router'
import { useAuthStore } from '../../stores/auth'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()

const labels = ref<string[]>([])
const connectionName = ref('')
const username = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)
const pageLoading = ref(true)

onMounted(async () => {
  error.value = ''
  try {
    labels.value = await auth.loadConnectionLabels()
    if (labels.value.length === 0) {
      const r = route.query.redirect
      await router.replace(
        typeof r === 'string' && r
          ? { path: '/register', query: { redirect: r } }
          : { path: '/register' },
      )
      return
    }
    connectionName.value = labels.value[0] ?? ''
  } catch (e: unknown) {
    const err = e as { response?: { data?: { error?: string } } }
    error.value = err.response?.data?.error || 'Failed to load connection labels'
  } finally {
    pageLoading.value = false
  }
})

async function submit() {
  error.value = ''
  loading.value = true
  try {
    if (!connectionName.value.trim()) {
      error.value = 'Select a connection label'
      return
    }
    await auth.login(connectionName.value, username.value, password.value)
    await auth.loadPublicHealth()
    const redir = route.query.redirect
    await router.push(typeof redir === 'string' && redir ? redir : '/')
  } catch (e: unknown) {
    const err = e as { response?: { data?: { error?: string } } }
    error.value = err.response?.data?.error || 'Request failed'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-[#f9fafb] text-[#111827] p-4">
    <p v-if="pageLoading" class="m-0 text-[#6b7280] text-[0.9rem]">Loading…</p>
    <form v-else class="w-full max-w-[400px] p-8 rounded-xl bg-[#ffffff] border border-[#e5e7eb] flex flex-col gap-3" @submit.prevent="submit">
      <h1 class="m-0 text-2xl">ChatDB</h1>
      <p class="m-0 text-[#6b7280] text-[0.9rem]">PostgreSQL &amp; MySQL viewer</p>
      <label class="flex flex-col gap-1 text-[0.85rem] text-[#6b7280]">
        Connection label
        <select v-model="connectionName" required class="px-[0.6rem] py-2 rounded-md border border-[#e5e7eb] bg-[#f9fafb] text-[#111827]">
          <option v-for="l in labels" :key="l" :value="l">{{ l }}</option>
        </select>
      </label>
      <label class="flex flex-col gap-1 text-[0.85rem] text-[#6b7280]">
        Username
        <input v-model="username" type="text" required autocomplete="username" class="px-[0.6rem] py-2 rounded-md border border-[#e5e7eb] bg-[#f9fafb] text-[#111827]" />
      </label>
      <label class="flex flex-col gap-1 text-[0.85rem] text-[#6b7280]">
        Password
        <input v-model="password" type="password" autocomplete="current-password" class="px-[0.6rem] py-2 rounded-md border border-[#e5e7eb] bg-[#f9fafb] text-[#111827]" />
      </label>

      <p v-if="error" class="text-[#ef4444] m-0 text-[0.85rem]">{{ error }}</p>
      <button type="submit" class="mt-2 p-[0.6rem] border-none rounded-md bg-[#059669] text-white font-semibold cursor-pointer disabled:opacity-60 disabled:cursor-default" :disabled="loading || pageLoading">{{ loading ? '…' : 'Login' }}</button>
      <p v-if="auth.hasUsers !== false" class="m-0 text-[0.85rem] text-center">
        <RouterLink to="/register" class="text-[#2563eb] no-underline hover:underline">Register new connection</RouterLink>
      </p>
    </form>
  </div>
</template>
