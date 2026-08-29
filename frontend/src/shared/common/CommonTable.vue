<template>
  <div :class="['text-center w-full bg-white overflow-x-auto custom-scrollbar-v2', tableResponsiveClass, tableClasses, props.tableHeight, props.borderRounded ? 'rounded-lg' : '', props.bordered ? 'border border-grayCust-160' : '']">
    <table
      class="w-full "
      :class="[ props.borderRounded ? 'rounded-lg' : '']"
    >
      <thead :class="['bg-grayCust-50 border-b']">
        <tr>
          <th
            v-if="props.selectTable"
            :class="['pl-8 py-3 w-10 cursor-pointer text-grayCust-600 font-medium text-xs border-b', props.bordered ? 'border border-grayCust-160' : '']"
            scope="col"
          >
            <Checkbox
              id="select_all_checkbox"
              :checked="selectAll"
              size="checkbox-sm"
              @change="toggleSelectAll"
            />
          </th>
          <th
            v-for="(field, colindex) in tableFields"
            :key="field.key"
            :class="[
              'px-4 py-3 text-grayCust-500 font-medium text-xs border-b',
              props.bordered ? 'border border-grayCust-160' : '',
              field.sortable ? 'cursor-pointer' : 'cursor-default',
              field.headerClass,
              field.hide ? 'hidden' : '',
              { 'sorted-asc': localSortBy === field.key && !localSortDesc, 'sorted-desc': localSortBy === field.key && localSortDesc }
            ]"
            role="columnheader"
            scope="col"
            tabindex="0"
            :style="field.width ? `min-width: ${field.width}; max-width: ${field.width};` : 'min-width: auto; max-width: auto;'"
            :aria-colindex="colindex + 1"
            :aria-sort="ariaSort[field.key]"
            @click="field.sortable && sortTable(field.key)"
          >
            <div
              class="flex items-center gap-2 uppercase whitespace-nowrap"
              :class="[field.headerClass, field.key === 'action' ? 'justify-center' : 'justify-start']"
            >
              <slot
                :name="`head_${field.key}`"
                :label="field.label"
              >
                {{ field.label }}
              </slot>
              <span v-if="localSortBy === field.key && field.sortable">
                <OutlineSortDescendingIcon
                  v-if="localSortDesc"
                  class="h-5 w-5 text-grayCust-500"
                />
                <OutlineSortAscendingIcon
                  v-else
                  class="h-5 w-5 text-grayCust-500"
                />
              </span>
            </div>
          </th>
        </tr>
      </thead>
      
      <!-- Loading State Body -->
      <tbody
        v-if="props.busy"
        class="bg-white"
      >
        <slot name="table-busy">
          <CommonTableSkeletonRow 
            v-for="i in 3" 
            :key="i"
            :fields="tableFields" 
            :select-table="props.selectTable"
            :row-height="props.rowHeight"
          />
        </slot>
      </tbody>
      
      <!-- Empty State -->
      <tbody
        v-else-if="!localItems.length && props.showEmpty"
        class="bg-white rounded-lg h-[100px]"
        :class="props.borderRounded && props.bordered ? '' : 'border-b border-grayCust-160'"
      >
        <tr>
          <td
            :colspan="tableFields.length + (props.selectTable ? 1 : 0)"
            class="px-4 py-2 text-center text-grayCust-500"
          >
            <slot name="empty">
              <span class="text-grayCust-750">No data available</span>
            </slot>
          </td>
        </tr>
      </tbody>
      
      <!-- Table Data -->
      <tbody v-else>
        <slot name="tbody-start" />
        <template
          v-for="(item, itemIndex) in localItems"
          :key="item.id ?? itemIndex"
        >
          <slot
            name="before-row"
            :item="item"
            :item-index="itemIndex"
          />
          <tr
            :class="[
              item.trClass,
              item.selected || activeRow === item ? 'bg-indigo-50' : (props.striped ? 'even:bg-grayCust-50' : ''),
              props.hover && !(item.selected || activeRow === item) ? 'hover:bg-grayCust-100 transition-colors duration-150' : '',
              props.bordered && props.borderRounded ? 'last:border-b-0' : '',
              // props.rowBordered && items.data.length > 1 && index !== items.data.length - 1 
              //   ? 'border-b border-grayCust-160' 
              //   : '',
              props.rowBordered && !props.bordered ? 'border-b border-grayCust-160' : '',
              props.rowBordered && props.bordered && localItems.length > 1
                ? 'border-b border-grayCust-160'
                : '',
            ]"
            :style="item.trStyle || ''"
            @dblclick="dbRowClicked($event, item)"
            @click="rowClicked($event, item)"
          >
            <td
              v-if="props.selectTable"
              class="pl-8 py-4 text-start"
              :class="[item.disabled ? 'bg-grayCust-50' : '', props.bordered ? 'border border-grayCust-160' : '']"
            >
              <Checkbox
                :checked="item.selected"
                :disabled="item.disabled"
                size="checkbox-sm"
                @change="rowSelectionChanged(item, !item.selected)"
              />
            </td>
            <td
              v-for="field in tableFields"
              :key="field.key"
              :class="['px-4 py-4 text-start relative', field.class, item.disabled ? 'bg-grayCust-50' : '', field.hide ? 'hidden' : '', props.bordered ? 'border border-grayCust-160' : '', 'hover:bg-grayCust-200 hover:shadow-inner transition-colors duration-150 cursor-pointer']"
              :style="field.width ? `min-width: ${field.width}; max-width: ${field.width};` : 'min-width: auto; max-width: auto;'"
            >
              <slot
                :name="`cell_${field.key}`"
                :item="item"
                :index="itemIndex"
                :value="getFormattedFieldValue(field, item)"
                :unformatted="item[field.key]"
                :field="field"
                :disabled="item.disabled"
              >
                {{ getFormattedFieldValue(field, item) }}
              </slot>
            </td>
          </tr>
          <slot
            name="after-row"
            :item="item"
            :item-index="itemIndex"
          />
        </template>
        <slot name="tbody-end" />
      </tbody>
    </table>
  </div>
</template>

<script setup>
import Checkbox from "@/shared/common/CommonCheckbox.vue";
import CommonTableSkeletonRow from "@/shared/common/CommonTableSkeleton.vue";
import { computed, ref, watch } from 'vue';

const props = defineProps({
  sortBy: {
    type: String,
    default: "",
  },
  sortDesc: {
    type: Boolean,
    default: true,
  },
  sortable: {
    type: Boolean,
    default: false,
  },
  responsive: {
    type: [String, Boolean],
    default: "",
  },
  busy: {
    type: Boolean,
    default: false,
  },
  hover: {
    type: Boolean,
    default: false,
  },
  striped: {
    type: Boolean,
    default: false,
  },
  bordered: {
    type: Boolean,
    default: false,
  },
  small: {
    type: Boolean,
    default: false,
  },
  fields: {
    type: Array,
    required: true,
  },
  items: {
    type: Array,
    required: true,
  },
  perPage: {
    type: Number,
    default: 10,
  },
  rowBordered: {
    type: Boolean,
    default: false
  },
  
  tableHeight: {
    type: String,
    default: '',
    required: false,
  },
  selectTable: {
    type: Boolean,
    default: false,
  },
  showEmpty: {
    type: Boolean,
    default: true,
  },
  borderRounded:{
    type: Boolean,
    default: false,
  },
  hide: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(['update:sortBy', 'update:sortDesc', 'sortChanged', 'rowClicked', 'rowDbClicked', 'selectionChanged', 'page-change']);

const activeRow = ref(null);
const localSortBy = ref(props.sortBy);
const localSortDesc = ref(props.sortDesc);
const selectAll = ref(false);
const localItems = ref(Array.isArray(props.items) ? [...props.items] : (Array.isArray(props.items?.data) ? [...props.items.data] : []));

const ariaSort = computed(() => {
  return props.fields.reduce((result, field) => {
    result[field.key] = field.sortable
      ? field.key === localSortBy.value
        ? localSortDesc.value
          ? "descending"
          : "ascending"
        : "none"
      : null;
    return result;
  }, {});
});

const tableClasses = computed(() => {
  return [
    "w-full border-collapse",
    props.striped ? "even:bg-grayCust-50" : "",
    // props.bordered ? "border border-grayCust-160" : "",
    // props.borderRounded ? "rounded-lg" : "",
  ].join(" ").trim();
});

// const tableItems = computed(() => {
//   if (!props.sortable) return localItems.value;

//   return [...localItems.value].sort((a, b) => {
//     const aValue = a[localSortBy.value];
//     const bValue = b[localSortBy.value];
//     if (aValue < bValue) return localSortDesc.value ? 1 : -1;
//     if (aValue > bValue) return localSortDesc.value ? -1 : 1;
//     return 0;
//   });
// });

const tableFields = computed(() => {
  if (!props.fields.length && localItems.value[0]) {
    return Object.keys(localItems.value[0]).map((key) => ({
      key,
      label: key,
      sortable: true,
    }));
  }
  return props.fields;
});

// const tableResponsiveClass = computed(() => {
//   if (!props.responsive) return "";
//   return props.responsive === true
//     ? "" 
//     : `overflow-x-auto sm:overflow-x-visible ${props.responsive}:overflow-x-auto`;
// });

const checkSelectAllState = () => {
  if (!localItems.value.length) {
    selectAll.value = false;
    return;
  }
  
  const items = localItems.value?.filter(item => !item.disabled) ?? [];
  selectAll.value = items.length > 0 && items.every(item => item.selected);
};

watch(() => props.items, (newItems) => {
  localItems.value = Array.isArray(newItems) ? [...newItems] : (Array.isArray(newItems?.data) ? [...newItems.data] : []);
  checkSelectAllState();
}, { deep: true, immediate: true });

const getFormattedFieldValue = (field, item) => {
  if (field.formatter && typeof field.formatter === "function") {
    return field.formatter(item[field.key], field.key, item);
  }
  if (field.key?.includes(".")) {
    return field.key.split(".").reduce((obj, key) => obj?.[key], item);
  }
  return item[field.key];
};

const sortTable = (key) => {
  if (localSortBy.value === key) {
    localSortDesc.value = !localSortDesc.value;
  } else {
    localSortBy.value = key;
    localSortDesc.value = false;
  }
  emit("update:sortBy", localSortBy.value);
  emit("update:sortDesc", localSortDesc.value);
  emit("sortChanged", localSortBy.value, localSortDesc.value);
};

const rowClicked = (e, item) => {
  activeRow.value = item;
  if (e.target.nodeName !== "INPUT") {
    emit("rowClicked", item);
  }
};

const dbRowClicked = (e, item) => {
  if (e.target.nodeName !== "INPUT") {
    emit("rowDbClicked", item);
  }
};


const toggleSelectAll = () => {
  selectAll.value = !selectAll.value;
  localItems.value.filter(item => !item.disabled).forEach((item) => {
    item.selected = selectAll.value;
  });
  emit("selectionChanged", localItems.value.filter((item) => item.selected));
};

const rowSelectionChanged = (item, selected) => {
  if (item.disabled) return;
  item.selected = selected;
  checkSelectAllState();
  emit("selectionChanged", localItems.value.filter((item) => item.selected));
};
</script>


<style scoped>
.scrollbar-hide::-webkit-scrollbar {
  display: none;
}
</style>

