<template>
  <input
    :id="id"
    v-model="proxyChecked"
    type="checkbox"
    :disabled="disabled"
    :value="value"
    :class="[computedCheckboxSize, computedCheckboxShape, computedCheckboxType]"
    class="border border-grayCust-350 shadow-sm focus:outline-none focus:ring-0 focus:ring-offset-0 transition-shadow cursor-pointer disabled:border-grayCust-160 disabled:bg-grayCust-50 disabled:cursor-not-allowed"
    :aria-checked="checked"
    :aria-disabled="disabled"
    :aria-label="id || 'Checkbox'"
  >
</template>

<script setup>
  import { computed } from 'vue';
  
  const props = defineProps({
    checked: {
      type: [Array, Boolean],
      default: false,
    },
    value: {
      type: [Array, Boolean],
      default: null,
    },
    id: {
      type: String,
      default: null,
    },
    size: {
      type: String,
      default: 'checkbox-sm',
      validator: (value) => ['checkbox-sm', 'checkbox-lg'].includes(value),
    },
    shape: {
      type: String,
      default: 'default',
      validator: (value) => ['default', 'rounded'].includes(value),
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    type: {
      type: String,
      default: 'secondary',
      validator: (value) => ['secondary', 'primary'].includes(value),
    },
  });
  
  const emit = defineEmits(['update:checked']);
  
  const computedCheckboxSize = computed(() => 
    props.size === 'checkbox-lg' ? 'w-5 h-5 rounded-md' : 'w-4 h-4'
  );
  
  const computedCheckboxShape = computed(() => 
    props.shape === 'rounded' ? '!rounded-full' : 'rounded'
  );
  
  const computedCheckboxType = computed(() => 
    props.type === 'primary' ? 'text-primary-800 focus:border-primary-800' : 'text-secondary-900 focus:border-secondary-900'
  );
  
  const proxyChecked = computed({
    get() {
      return props.checked;
    },
    set(val) {
      emit("update:checked", val);
    },
  });
  </script>
  
  
  
  <style scoped>
  .check-border,
  .check-border:focus {
    border-radius: 1.75px;
  }
  </style>