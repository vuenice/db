<template>
  <TransitionRoot
    as="template"
    :show="isOpen"
    @close="close"
  >
    <Dialog
      as="div"
      class="fixed inset-0 z-40 flex items-center justify-center"
    >
      <div class="fixed inset-0">
        <TransitionChild
          as="template"
          enter="ease-out duration-300"
          enter-from="opacity-0"
          enter-to="opacity-100"
          leave="ease-in duration-200"
          leave-from="opacity-100"
          leave-to="opacity-0"
        >
          <DialogOverlay
            class="fixed inset-0 z-40 bg-black bg-opacity-50"
          />
        </TransitionChild>
      </div>
      
      <div class=" flex items-center justify-center w-full max-h-screen px-4 sm:px-6">
        <TransitionChild
          as="template"
          enter="ease-out duration-300"
          enter-from="opacity-0 scale-95"
          enter-to="opacity-100 scale-100"
          leave="ease-in duration-200"
          leave-from="opacity-100 scale-100"
          leave-to="opacity-0 scale-95"
        >
          <div
            :class="['w-full mx-auto my-6 md:my-8 text-left transform', bgColor, roundedClass,'shadow-xl  max-h-[95vh] overflow-hidden', modalSizeClass]"

          >
            <div class="flex flex-col h-full max-h-[95vh]">
              <!-- Header -->
              <div
                v-if="$slots.header"
                class="border-b sm:py-4 flex-shrink-0 mx-6"
                :class="extraHeaderClass"
              >
                <slot name="header" />
              </div>
              
              <!-- Main -->
              <div
                class="flex-grow overflow-y-auto custom-scrollbar-v2"
                :class="modalPadding"
              >
                <slot name="content" />
              </div>
              
              <!-- Footer -->
              <div
                v-if="$slots.footer"
                class="border-t p-4 sm:p-6 sm:py-3 pt-3 pb-3 flex-shrink-0"
                :class="extraFooterClass"
              >
                <slot name="footer" />
              </div>
            </div>
          </div>
        </TransitionChild>
      </div>
    </Dialog>
  </TransitionRoot>
</template>

<script setup>
  import { Dialog, DialogOverlay, TransitionChild, TransitionRoot } from '@headlessui/vue'
import { computed, defineEmits, defineProps, ref, watch } from 'vue'
  
  const props = defineProps({
  modelValue: Boolean,
  size: {
    type: String,
    default: '2xl',
    validator: val => ['sm', 'md', 'lg', 'xl', '2xl', '3xl', '4xl', '5xl', '6xl', '7xl','h-auto'].includes(val)
  },
  footer: { type: Boolean, default: true },
  header: { type: Boolean, default: true },
  fullHeader: { type: Boolean, default: true },
  extraHeaderClass: { type: String, default: '' },
  extraFooterClass: { type: String, default: '' },
  modalPadding: { type: String, default: 'p-4 sm:p-6' },
  position: {
    type: String,
    default: 'center',
    validator: val => ['top', 'center'].includes(val)
  },
  allowPadding: { type: Boolean, default: true },
  preventClose: { type: Boolean, default: false },
  bgColor: {
  type: String,
  default: 'bg-white'
},
rounded: {
    type: String,
    default: 'lg' 
  }
})

const emit = defineEmits(['update:modelValue', 'close'])
const isOpen = ref(props.modelValue)
watch(() => props.modelValue, val => { isOpen.value = val })
function close() {
  if (props.preventClose) {
    return
  }

  isOpen.value = false
  emit('update:modelValue', false)
  emit('close')
}

const modalSizeClass = computed(() => {
  return {
    sm: 'max-w-sm',
    md: 'max-w-md',
    lg: 'max-w-lg',
    xl: 'max-w-xl',
    '2xl': 'max-w-2xl',
    '3xl': 'max-w-3xl',
    '4xl': 'max-w-4xl',
    '5xl': 'max-w-5xl',
    '6xl': 'max-w-6xl',
    '7xl': 'max-w-7xl'
  }[props.size] || 'max-w-2xl'
})
const roundedClass = computed(() => {
  // Allows rounded values like 'none', 'sm', '', 'md', 'lg', 'xl', '2xl', '3xl', 'full'
  return `rounded-${props.rounded || 'lg'}`; // Defaults to 'lg' if not set
});

</script>
