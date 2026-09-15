<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { useRouter, RouterLink } from 'vue-router'
import { useAuthStore } from '../../stores/auth'

const router = useRouter()
const auth = useAuthStore()

const error = ref('')
const loading = ref(false)

const errorHint = computed(() =>
  error.value.includes('email and password required')
    ? 'Please fill Host, Database, and Database username. Database password is optional.'
    : error.value.includes('Access denied for user')
    ? 'hint: Incorrect Database username or Database password'
    : '',
)

const conn = ref({
  connection_name: '',
  driver: 'postgres' as 'postgres' | 'mysql',
  host: '127.0.0.1',
  port: 5432,
  database: '',
  ssl_mode: 'disable',
  db_username: '',
  db_password: '',
  use_ssh: false,
  ssh_host: '',
  ssh_port: 22,
  ssh_user: '',
  ssh_password: '',
  ssh_key: ''
})

watch(
  () => conn.value.driver,
  (d) => {
    conn.value.port = d === 'mysql' ? 3306 : 5432
  },
)

async function submit() {
  error.value = ''
  loading.value = true
  try {
    if (!conn.value.connection_name.trim()) {
      error.value = 'Connection label is required'
      return
    }
    if (!conn.value.host.trim() || !conn.value.database.trim() || !conn.value.db_username.trim()) {
      error.value = 'Host, Database, and Database username are required'
      return
    }
    await auth.register({
      connection_name: conn.value.connection_name.trim(),
      driver: conn.value.driver,
      host: conn.value.host,
      port: conn.value.port,
      database: conn.value.database,
      ssl_mode: conn.value.ssl_mode,
      read_username: conn.value.db_username,
      read_password: conn.value.db_password,
      use_ssh: conn.value.use_ssh,
      ssh_host: conn.value.ssh_host,
      ssh_port: conn.value.ssh_port,
      ssh_user: conn.value.ssh_user,
      ssh_password: conn.value.ssh_password,
      ssh_key: conn.value.ssh_key,
    })
    conn.value.db_password = ''
    await auth.loadPublicHealth()
    const redir = router.currentRoute.value.query.redirect
    await router.push(typeof redir === 'string' && redir ? redir : '/')
  } catch (e: unknown) {
    const err = e as { response?: { data?: { error?: string } } }
    const msg = err.response?.data?.error || 'Request failed'
    error.value = msg
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-[#f9fafb] text-[#111827] p-4">
    <form class="w-full max-w-[480px] p-8 rounded-xl bg-[#ffffff] border border-[#e5e7eb] flex flex-col gap-3" @submit.prevent="submit">
      <h1 class="m-0 text-2xl">ChatDB</h1>
      <p class="m-0 text-[#6b7280] text-[0.9rem]">Register new connection</p>
      <p class="m-0 px-[0.85rem] py-[0.75rem] rounded-lg border border-[#e5e7eb] bg-[#f9fafb] text-[#6b7280] text-[0.8rem] leading-[1.45]">
        For the admin, login username and password are the same as the database username and password.
        Other users are created by the admin and cannot register here.
      </p>
      <label class="flex flex-col gap-1 text-[0.85rem] text-[#6b7280]"
        >Connection label
        <input v-model="conn.connection_name" type="text" required placeholder="e.g. production" class="px-[0.6rem] py-2 rounded-md border border-[#e5e7eb] bg-[#f9fafb] text-[#111827]" />
      </label>
      <label class="flex flex-col gap-1 text-[0.85rem] text-[#6b7280]"
        >Driver
        <select v-model="conn.driver" class="px-[0.6rem] py-2 rounded-md border border-[#e5e7eb] bg-[#f9fafb] text-[#111827]">
          <option value="postgres">PostgreSQL</option>
          <option value="mysql">MySQL / MariaDB</option>
        </select>
      </label>
      <label class="flex flex-col gap-1 text-[0.85rem] text-[#6b7280]">Host <input v-model="conn.host" required class="px-[0.6rem] py-2 rounded-md border border-[#e5e7eb] bg-[#f9fafb] text-[#111827]" /></label>
      <label class="flex flex-col gap-1 text-[0.85rem] text-[#6b7280]">Port <input v-model.number="conn.port" type="number" class="px-[0.6rem] py-2 rounded-md border border-[#e5e7eb] bg-[#f9fafb] text-[#111827]" /></label>
      <label class="flex flex-col gap-1 text-[0.85rem] text-[#6b7280]"
        >Database Name
        <input v-model="conn.database" required placeholder="Default database" class="px-[0.6rem] py-2 rounded-md border border-[#e5e7eb] bg-[#f9fafb] text-[#111827]" />
      </label>
      <label v-if="conn.driver === 'postgres'" class="flex flex-col gap-1 text-[0.85rem] text-[#6b7280]"
        >SSL mode <input v-model="conn.ssl_mode" placeholder="disable" class="px-[0.6rem] py-2 rounded-md border border-[#e5e7eb] bg-[#f9fafb] text-[#111827]"
      /></label>
      <label class="flex flex-col gap-1 text-[0.85rem] text-[#6b7280]"
        >Database username
        <input v-model="conn.db_username" required autocomplete="off" placeholder="e.g. root" class="px-[0.6rem] py-2 rounded-md border border-[#e5e7eb] bg-[#f9fafb] text-[#111827]" />
      </label>
      <label class="flex flex-col gap-1 text-[0.85rem] text-[#6b7280]"
        >Database password
        <input
          v-model="conn.db_password"
          type="password"
          autocomplete="new-password"
          placeholder="leave empty for no password"
          class="px-[0.6rem] py-2 rounded-md border border-[#e5e7eb] bg-[#f9fafb] text-[#111827]"
        />
      </label>

      <!-- SSH Toggle -->
      <label class="flex items-center gap-2 text-[0.85rem] text-[#6b7280] mt-2 cursor-pointer">
        <input type="checkbox" v-model="conn.use_ssh" class="w-4 h-4 rounded border-[#e5e7eb] text-[#059669] focus:ring-[#059669]">
        Use SSH Tunnel
      </label>

      <template v-if="conn.use_ssh">
        <div class="p-4 bg-[#f3f4f6] rounded-md flex flex-col gap-3 mt-1">
          <p class="m-0 text-[#111827] font-medium text-[0.85rem]">SSH Configuration</p>
          <label class="flex flex-col gap-1 text-[0.85rem] text-[#6b7280]">SSH Host <input v-model="conn.ssh_host" required placeholder="e.g. 192.168.1.100" class="px-[0.6rem] py-2 rounded-md border border-[#e5e7eb] bg-[#ffffff] text-[#111827]" /></label>
          <label class="flex flex-col gap-1 text-[0.85rem] text-[#6b7280]">SSH Port <input v-model.number="conn.ssh_port" type="number" required class="px-[0.6rem] py-2 rounded-md border border-[#e5e7eb] bg-[#ffffff] text-[#111827]" /></label>
          <label class="flex flex-col gap-1 text-[0.85rem] text-[#6b7280]">SSH Username <input v-model="conn.ssh_user" required placeholder="e.g. root" class="px-[0.6rem] py-2 rounded-md border border-[#e5e7eb] bg-[#ffffff] text-[#111827]" /></label>
          <label class="flex flex-col gap-1 text-[0.85rem] text-[#6b7280]">SSH Password <input v-model="conn.ssh_password" type="password" placeholder="leave empty if using key" class="px-[0.6rem] py-2 rounded-md border border-[#e5e7eb] bg-[#ffffff] text-[#111827]" /></label>
          <label class="flex flex-col gap-1 text-[0.85rem] text-[#6b7280]">SSH Private Key (Optional) <textarea v-model="conn.ssh_key" rows="3" placeholder="-----BEGIN PRIVATE KEY-----..." class="px-[0.6rem] py-2 rounded-md border border-[#e5e7eb] bg-[#ffffff] text-[#111827] font-mono text-xs"></textarea></label>
        </div>
      </template>

      <p v-if="error" class="text-[#ef4444] m-0 text-[0.85rem]">{{ error }}</p>
      <p v-if="errorHint" class="text-[#d29922] -mt-1 mb-0 text-[0.85rem]">{{ errorHint }}</p>
      <button type="submit" class="mt-2 p-[0.6rem] border-none rounded-md bg-[#059669] text-white font-semibold cursor-pointer disabled:opacity-60 disabled:cursor-default" :disabled="loading">{{ loading ? '…' : 'Continue' }}</button>
      <p v-if="auth.hasUsers !== false" class="m-0 text-[0.85rem] text-center">
        <RouterLink to="/login" class="text-[#2563eb] no-underline hover:underline">Already have an account? Sign in</RouterLink>
      </p>
    </form>
  </div>
</template>
