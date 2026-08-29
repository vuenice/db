<script setup lang="ts">
// @ts-nocheck
import { useRoute } from 'vue-router'
import {
  CubeIcon,
  GlobeIcon,
  ViewGridIcon,
  CogIcon,
  CloudUploadIcon,
  CurrencyDollarIcon,
  PuzzleIcon,
  AdjustmentsIcon
} from '@heroicons/vue/outline'

const props = defineProps<{
  isCollapsed: boolean
}>()

const route = useRoute()

const menuItems = [
  { name: 'Sites', icon: CubeIcon, path: '/workbench', activeKeys: ['workbench'] },
  { name: 'Domains', icon: GlobeIcon, path: '/domains', activeKeys: ['domains'] },
  { name: 'Snapshots', icon: ViewGridIcon, path: '/snapshots', activeKeys: ['snapshots'] },
  { name: 'Manage', icon: AdjustmentsIcon, path: '/manage', activeKeys: ['manage'] },
  { name: 'Migrations', icon: CloudUploadIcon, path: '/migrations', activeKeys: ['migrations'] },
  { name: 'Sell', icon: CurrencyDollarIcon, path: '/sell', activeKeys: ['sell'] },
  { name: 'Settings', icon: CogIcon, path: '/settings', activeKeys: ['settings'] },
  { name: 'Integrations', icon: PuzzleIcon, path: '/integrations', activeKeys: ['integrations'] },
]

// Determine if active based on path
const isActive = (item: any) => {
  return route.path.startsWith(item.path) || (route.path === '/' && item.path === '/workbench')
}
</script>

<template>
  <div 
    class="h-full border-r border-grayCust-160 transition-all duration-300 ease-in-out flex flex-col" 
    :class="[isCollapsed ? 'w-14' : 'w-[240px]', 'bg-grayCust-50']"
  >
    <div class="flex-1 overflow-y-auto flex flex-col py-4 [&::-webkit-scrollbar]:w-1 [&::-webkit-scrollbar-track]:bg-transparent [&::-webkit-scrollbar-thumb]:bg-slate-300 [&::-webkit-scrollbar-thumb]:rounded-full">
      <nav class="px-2">
        <ul class="space-y-1">
          <li v-for="item in menuItems" :key="item.name">
            <router-link
              :to="item.path"
              class="flex items-center gap-3 px-3 py-2.5 rounded-md transition-colors text-sm font-medium"
              :class="[
                isActive(item)
                  ? 'bg-white text-emerald-600 shadow-sm border border-gray-200'
                  : 'text-grayCust-640 hover:bg-gray-100 hover:text-gray-900 border border-transparent'
              ]"
              :title="item.name"
            >
              <component :is="item.icon" class="w-5 h-5 shrink-0" :class="isActive(item) ? 'text-emerald-500' : 'text-gray-400'" />
              <span v-if="!isCollapsed" class="truncate">{{ item.name }}</span>
            </router-link>
          </li>
        </ul>
      </nav>
      
      <!-- Complete your profile placeholder -->
      <div v-if="!isCollapsed" class="mt-auto px-4 py-3">
        <div class="border border-gray-200 rounded-md bg-white p-3 shadow-sm relative overflow-hidden">
          <div class="absolute top-0 left-0 right-0 h-1 bg-gradient-to-r from-teal-400 via-emerald-300 to-amber-200"></div>
          <div class="flex justify-between items-center mb-1 mt-1">
            <h4 class="text-sm font-bold text-gray-800">Complete your profile!</h4>
            <button class="text-gray-400 hover:text-gray-600">&times;</button>
          </div>
          <div class="text-xs text-gray-500 mb-2">2/5 steps completed</div>
          
          <ul class="space-y-2 mt-3">
            <li class="flex items-center gap-2 text-xs text-gray-400 line-through">
              <div class="w-4 h-4 rounded-full border border-gray-300 flex items-center justify-center shrink-0">
                <svg class="w-2 h-2 text-gray-400" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"></path></svg>
              </div>
              <span>Add a credit card</span>
            </li>
            <li class="flex items-center gap-2 text-xs text-gray-400 line-through">
              <div class="w-4 h-4 rounded-full border border-gray-300 flex items-center justify-center shrink-0">
                <svg class="w-2 h-2 text-gray-400" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"></path></svg>
              </div>
              <span>Create Your First Site</span>
            </li>
            <li class="flex items-center gap-2 text-xs text-gray-700">
              <div class="w-4 h-4 rounded-full border border-emerald-500 flex items-center justify-center shrink-0"></div>
              <span class="font-medium text-emerald-700">Save Snapshot</span>
            </li>
            <li class="flex items-center gap-2 text-xs text-gray-700">
              <div class="w-4 h-4 rounded-full border border-emerald-500 flex items-center justify-center shrink-0"></div>
              <span class="font-medium text-emerald-700">Map Custom Domain</span>
            </li>
            <li class="flex items-center gap-2 text-xs text-gray-700">
              <div class="w-4 h-4 rounded-full border border-emerald-500 flex items-center justify-center shrink-0"></div>
              <span class="font-medium text-emerald-700">Connect a site</span>
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

