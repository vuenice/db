<script setup lang="ts">
// @ts-nocheck
import { Popover, PopoverButton, PopoverPanel } from '@headlessui/vue'
import { ChevronDownIcon, CogIcon } from '@heroicons/vue/outline'

const myTeam = {
  heading: 'My Teams',
  teams: [
    { id: 1, name: 'My Team', image: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80' },
  ],
}

const selectedOption = myTeam.teams[0]
</script>

<template>
  <Popover class="relative">
    <PopoverButton
      class="flex items-center gap-2 cursor-pointer rounded-md bg-white/10 py-[3px] px-2.5 text-left text-white border border-white/20 focus:outline-none focus-visible:ring-2 focus-visible:ring-white/75"
    >
      <div class="flex items-center gap-2">
        <span class="flex items-center gap-2">
          <div
            class="w-6 h-6 shrink-0 rounded-full bg-emerald-100 flex items-center justify-center text-emerald-700 font-bold text-xs"
          >
            MT
          </div>
          <span class="block truncate text-sm font-medium">{{ selectedOption?.name }}</span>
        </span>
        <ChevronDownIcon class="w-4 h-4 text-gray-300" aria-hidden="true" />
      </div>
    </PopoverButton>

    <transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="translate-y-1 opacity-0"
      enter-to-class="translate-y-0 opacity-100"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="translate-y-0 opacity-100"
      leave-to-class="translate-y-1 opacity-0"
    >
      <PopoverPanel
        class="absolute right-0 z-50 mt-2 w-56 origin-top-right rounded-md bg-white shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none"
      >
        <div class="mb-1">
          <div class="flex items-center justify-between gap-2 px-4 pt-3 pb-2">
            <span class="text-grayCust-700 font-medium text-sm">{{ myTeam.heading }}</span>
            <button class="flex items-center gap-1 text-sm text-grayCust-980 hover:text-black">
              <span class="text-xs">+ Create New</span>
            </button>
          </div>
          <div class="max-h-[200px] overflow-y-auto custom-scrollbar-v2 pb-1">
            <div
              v-for="team in myTeam.teams"
              :key="team.id"
              class="py-2 px-4 cursor-pointer relative hover:bg-grayCust-100 flex items-center text-sm font-normal gap-2 justify-between bg-emerald-50 text-grayCust-700 font-semibold"
            >
              <div class="flex items-center gap-3">
                <div
                  class="w-6 h-6 shrink-0 rounded-full bg-emerald-100 flex items-center justify-center text-emerald-700 font-bold text-xs circle-shadow"
                >
                  MT
                </div>
                <span class="block truncate w-[110px]">{{ team.name }}</span>
              </div>
              <div class="flex items-center gap-2">
                <CogIcon class="w-4 h-4 text-gray-400 hover:text-gray-600" />
              </div>
            </div>
          </div>
        </div>
      </PopoverPanel>
    </transition>
  </Popover>
</template>

<style scoped>
.circle-shadow {
  box-shadow: 0 0 0 1px #fff, 0 0 0 2px #11bf85;
}
.custom-scrollbar-v2::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar-v2::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar-v2::-webkit-scrollbar-thumb {
  background-color: #cbd5e1;
  border-radius: 20px;
}
</style>
