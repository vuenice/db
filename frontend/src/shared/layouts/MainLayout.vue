<script setup lang="ts">
import { ref } from 'vue'
import Header from '../common/Header/Header.vue'
import Sidebar from '../common/Sidebar/Sidebar.vue'

const isCollapsed = ref(false)

const toggleSidebar = () => {
  isCollapsed.value = !isCollapsed.value
}
</script>

<template>
  <div class="relative h-screen flex flex-col overflow-hidden text-gray-900 bg-white">
    <Header class="flex-none z-[10]" @toggle-sidebar="toggleSidebar" />
    <div class="flex flex-1 overflow-hidden">
      <Sidebar
        :is-collapsed="isCollapsed"
        class="h-full z-10 bg-grayCust-50 border-r border-gray-200 transition-all duration-300 ease-in-out"
      />
      <div class="flex-1 overflow-hidden transition-all duration-300 ease-in-out relative">
        <div class="h-full overflow-y-auto overflow-x-hidden p-6 bg-white">
          <slot />
        </div>
      </div>
    </div>
  </div>
</template>

