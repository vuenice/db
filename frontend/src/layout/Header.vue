<script setup lang="ts">
// @ts-nocheck
import {
  MenuIcon,
  ChatIcon,
  BookOpenIcon,
  VideoCameraIcon,
  QuestionMarkCircleIcon,
} from '@heroicons/vue/outline'
import MyTeamList from './MyTeamList.vue'

defineEmits(['toggle-sidebar'])
</script>

<template>
  <header class="flex justify-between items-center w-full bg-primary-900 py-2 px-4 h-[50px]">
    <div class="h-[24px] flex items-center">
      <router-link to="/" class="flex items-center gap-2">
        <!-- Placeholder for logo, using text for now if image is missing -->
        <div class="text-white font-bold text-xl tracking-wider flex items-center">
          <svg class="w-6 h-6 mr-1" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M12 2L2 7L12 12L22 7L12 2Z" fill="currentColor"/>
            <path d="M2 17L12 22L22 17" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M2 12L12 17L22 12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
          W
        </div>
      </router-link>
    </div>
    <div class="items-center gap-3 flex">
      <button class="border border-white/20 hover:bg-white/10 text-white rounded p-1.5 focus:outline-none transition-colors">
        <ChatIcon class="w-4 h-4" />
      </button>

      <button class="border border-white/20 hover:bg-white/10 text-white rounded p-1.5 focus:outline-none transition-colors">
        <BookOpenIcon class="w-4 h-4" />
      </button>

      <button class="border border-white/20 hover:bg-white/10 text-white rounded p-1.5 focus:outline-none transition-colors">
        <VideoCameraIcon class="w-4 h-4" />
      </button>

      <button class="border border-white/20 hover:bg-white/10 text-white rounded p-1.5 focus:outline-none transition-colors">
        <QuestionMarkCircleIcon class="w-4 h-4" />
      </button>

      <button
        class="lg:hidden flex border border-white/20 hover:bg-white/10 text-white rounded p-1.5 focus:outline-none transition-colors"
        @click="$emit('toggle-sidebar')"
      >
        <MenuIcon class="w-4 h-4" />
      </button>

      <div class="h-[24px] w-[1px] bg-white/20 hidden lg:block" />

      <MyTeamList class="hidden lg:block" />

      <div class="h-[24px] w-[1px] bg-white/20 hidden lg:block" />

      <!-- User Profile placeholder -->
      <button class="hidden lg:flex w-8 h-8 rounded-full bg-emerald-600 border border-white/20 items-center justify-center text-white text-sm font-semibold">
        Y
      </button>
    </div>
  </header>
</template>
