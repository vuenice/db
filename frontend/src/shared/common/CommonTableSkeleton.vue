<template>
  <tr
    class="border-b border-grayCust-160 last:border-b-0 animate-pulse w-full"
    :class="[rowHeight]"
  >
    <td
      v-if="selectTable" 
      class="pl-10 py-4 w-10"
    >
      <div class="h-4 w-4 bg-grayCust-160 rounded" />
    </td>
    <td
      v-for="(field, index) in fields" 
      :key="index"
      class="px-4 py-3 w-fit"
      :class="[field.class]"
      :style="field.style ||''"
    >
      <div 
        class="h-6 bg-grayCust-160 rounded" 
      />
    </td>
  </tr>
</template>

<script setup>
defineProps({
  fields: {
    type: Array,
    required: true
  },
  selectTable: {
    type: Boolean,
    default: false
  },
  rowHeight: {
    type: String,
    default: 'h-11'
  }
});
</script>

