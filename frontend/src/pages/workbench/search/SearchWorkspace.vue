<script setup lang="ts">
import { ref, computed } from 'vue'
import { http } from '../../../api/http'
import { useDatabaseStore, type TableMeta } from '../../../stores/database'
import { storeToRefs } from 'pinia'
import CommonTable from '@/shared/common/CommonTable.vue'
// @ts-ignore
import { SearchIcon } from '@heroicons/vue/outline'

const dbStore = useDatabaseStore()
const { selectedConnId, selectedPhysicalDatabase, tables, currentConnection } = storeToRefs(dbStore)

// --- Search Tab Logic ---
const searchTableListSearch = ref('')
interface SearchActiveTable {
  id: string
  schema: string
  name: string
  columns: { column: string; data_type: string; is_nullable: string }[]
  filterCol: string
  filterOp: string
  filterVal: string
  dataPreview: { columns: string[]; rows: unknown[][]; row_count: number } | null
  loading: boolean
  error: string
}
const searchActiveTables = ref<SearchActiveTable[]>([])

const searchFilteredTables = computed(() => {
  const q = searchTableListSearch.value.trim().toLowerCase()
  if (!q) return tables.value
  return tables.value.filter((t: TableMeta) => `${t.name} ${t.schema}`.toLowerCase().includes(q))
})

const searchTableGroups = computed<{ letter: string; items: TableMeta[] }[]>(() => {
  const items = [...searchFilteredTables.value].sort((a, b) => a.name.localeCompare(b.name))
  const byLetter = new Map<string, TableMeta[]>()
  for (const t of items) {
    const first = (t.name || '').trim().charAt(0).toUpperCase()
    const letter = /^[A-Z]$/.test(first) ? first : '#'
    if (!byLetter.has(letter)) byLetter.set(letter, [])
    byLetter.get(letter)!.push(t)
  }
  const letters = Array.from(byLetter.keys()).sort((a, b) => {
    if (a === '#') return 1
    if (b === '#') return -1
    return a.localeCompare(b)
  })
  return letters.map((l) => ({ letter: l, items: byLetter.get(l)! }))
})

function quoteIdent(driver: string, name: string): string {
  if (driver.toLowerCase() === 'mysql') return '`' + name.replace(/`/g, '``') + '`'
  return '"' + name.replace(/"/g, '""') + '"'
}

function quoteLiteral(v: unknown): string {
  if (v === null || v === undefined) return 'NULL'
  if (typeof v === 'number' && Number.isFinite(v)) return String(v)
  if (typeof v === 'boolean') return v ? '1' : '0'
  const s = String(v)
  return "'" + s.replace(/'/g, "''") + "'"
}

async function addTableToSearch(t: TableMeta) {
  const id = Date.now().toString() + Math.random().toString().slice(2, 6)
  const item: SearchActiveTable = {
    id,
    schema: t.schema,
    name: t.name,
    columns: [],
    filterCol: '',
    filterOp: '=',
    filterVal: '',
    dataPreview: null,
    loading: true,
    error: ''
  }
  searchActiveTables.value.push(item)

  if (selectedConnId.value) {
    try {
      const { data } = await http.get<{ columns: { column: string; data_type: string; is_nullable: string }[] }>(
        `/api/connections/${selectedConnId.value}/columns`,
        { params: { schema: item.schema, table: item.name, ...dbStore.dbParams() } },
      )
      
      // Update the reactive proxy, not the raw object
      const targetItem = searchActiveTables.value.find(t => t.id === id)
      if (targetItem) {
        targetItem.columns = data.columns
        if (targetItem.columns.length > 0) {
          targetItem.filterCol = targetItem.columns[0].column
        }
      }
    } catch (e: any) {
      const targetItem = searchActiveTables.value.find(t => t.id === id)
      if (targetItem) targetItem.error = 'Failed to load columns'
    }
  }
  
  const targetItem = searchActiveTables.value.find(t => t.id === id)
  if (targetItem) targetItem.loading = false
}

function removeSearchTable(id: string) {
  searchActiveTables.value = searchActiveTables.value.filter(t => t.id !== id)
}

async function executeSearchOnTable(item: SearchActiveTable) {
  item.error = ''
  item.dataPreview = null
  if (!selectedConnId.value) return
  if (!item.filterCol) {
    item.error = 'Please select a column.'
    return
  }

  const conn = currentConnection.value
  if (!conn) return

  item.loading = true
  try {
    const q = (n: string) => quoteIdent(conn.driver, n)
    const tableExpr = conn.driver.toLowerCase() === 'mysql' ? `${q(item.schema)}.${q(item.name)}` : `${q(item.schema)}.${q(item.name)}`
    
    const colExpr = q(item.filterCol)
    let op = item.filterOp
    let val = item.filterVal
    let clause = ''

    if (op === '=') {
      clause = `${colExpr} = ${quoteLiteral(val)}`
    } else if (op === 'LIKE') {
      clause = `${colExpr} LIKE ${quoteLiteral(val)}`
    } else if (op === '%LIKE%') {
      clause = `${colExpr} LIKE ${quoteLiteral('%' + val + '%')}`
    } else if (op === 'LIKE%') {
      clause = `${colExpr} LIKE ${quoteLiteral(val + '%')}`
    }

    const sql = `SELECT * FROM ${tableExpr} WHERE ${clause} LIMIT 500;`

    const { data } = await http.post(`/api/connections/${selectedConnId.value}/sql/execute`, {
      sql,
      pool: 'read',
      max_rows: 500,
      database: selectedPhysicalDatabase.value || undefined
    })
    
    item.dataPreview = data.result as any
  } catch (e: unknown) {
    const err = e as { response?: { data?: { error?: string } } }
    item.error = err.response?.data?.error || 'Search failed'
  } finally {
    item.loading = false
  }
}
</script>

<template>
  <div class="flex flex-row flex-1 overflow-hidden h-full w-full">
    <!-- Left Sidebar (20%) -->
    <aside class="w-[20%] min-w-[200px] max-w-[300px] border-r border-[#e5e7eb] flex flex-col bg-[#ffffff] shrink-0 h-full">
      <div class="p-3 border-b border-[#e5e7eb]">
        <input
          v-model="searchTableListSearch"
          class="w-full px-3 py-2 border border-gray-300 rounded text-sm focus:outline-none focus:border-emerald-400"
          type="search"
          placeholder="Search tables…"
          autocomplete="off"
        />
      </div>
      <div class="flex-1 overflow-y-auto custom-scrollbar-v2" role="list">
        <section v-for="g in searchTableGroups" :key="g.letter" role="listitem">
          <div class="font-semibold text-gray-500 bg-gray-50 px-3 py-1 text-sm border-b border-gray-100">{{ g.letter }}</div>
          <div class="px-2 py-1">
            <button
              v-for="t in g.items"
              :key="t.schema + '.' + t.name"
              type="button"
              class="flex justify-between items-center w-full px-2 py-1.5 hover:bg-gray-100 rounded group transition-colors"
              @click="addTableToSearch(t)"
            >
              <div class="flex flex-col text-left overflow-hidden">
                <span class="text-sm font-medium text-gray-700 group-hover:text-emerald-500 truncate">{{ t.name }}</span>
                <span class="text-xs text-gray-400">{{ t.kind === 'view' ? 'view' : 'table' }}</span>
              </div>
              <span class="text-gray-300 group-hover:text-emerald-500 text-lg leading-none font-bold shrink-0 ml-2" title="Add to search workspace">+</span>
            </button>
          </div>
        </section>
      </div>
    </aside>

    <!-- Right Content (80%) -->
    <div class="flex-1 overflow-y-auto p-4 bg-[#f9fafb] h-full">
      <div v-if="searchActiveTables.length === 0" class="flex flex-col items-center justify-center h-full text-gray-400">
        <SearchIcon class="w-12 h-12 mb-3 text-gray-300" />
        <p>Click the '+' icon on a table in the sidebar to add it to your search workspace.</p>
      </div>
      <div v-else class="space-y-6">
        <div v-for="item in searchActiveTables" :key="item.id" class="bg-white border border-[#e5e7eb] rounded-md shadow-sm overflow-hidden flex flex-col">
          <div class="bg-gray-50 border-b border-[#e5e7eb] px-4 py-3 flex items-center justify-between">
            <h3 class="text-sm font-semibold text-gray-800">{{ item.schema }}.{{ item.name }}</h3>
            <button type="button" @click="removeSearchTable(item.id)" class="text-gray-400 hover:text-red-500 transition-colors p-1 rounded hover:bg-gray-200" title="Remove table">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
            </button>
          </div>
          
          <div class="p-4 border-b border-[#e5e7eb] flex items-center gap-3 flex-wrap bg-white">
            <select v-model="item.filterCol" class="border border-gray-300 rounded px-3 py-1.5 text-sm bg-white min-w-[150px] focus:outline-none focus:border-emerald-400 focus:ring-1 focus:ring-emerald-400">
              <option v-for="col in item.columns" :key="col.column" :value="col.column">{{ col.column }}</option>
            </select>
            <select v-model="item.filterOp" class="border border-gray-300 rounded px-3 py-1.5 text-sm bg-white focus:outline-none focus:border-emerald-400 focus:ring-1 focus:ring-emerald-400">
              <option value="=">=</option>
              <option value="LIKE">LIKE</option>
              <option value="%LIKE%">%LIKE%</option>
              <option value="LIKE%">LIKE%</option>
            </select>
            <input v-model="item.filterVal" type="text" class="border border-gray-300 rounded px-3 py-1.5 text-sm min-w-[200px] flex-1 focus:outline-none focus:border-emerald-400 focus:ring-1 focus:ring-emerald-400" placeholder="Search value..." @keydown.enter="executeSearchOnTable(item)" />
            <button type="button" @click="executeSearchOnTable(item)" class="bg-emerald-500 hover:bg-emerald-600 text-white px-4 py-1.5 rounded text-sm font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-1 min-w-[90px]" :disabled="item.loading">
              <svg v-if="item.loading" class="animate-spin h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              <span>{{ item.loading ? 'Searching' : 'Search' }}</span>
            </button>
          </div>
          
          <div class="p-0 bg-white relative overflow-hidden">
            <div v-if="item.error" class="text-red-500 text-sm py-4 px-4">{{ item.error }}</div>
            <div v-else-if="item.dataPreview" class="overflow-x-auto w-full custom-scrollbar-v2">
              <div class="min-w-[500px]">
                <CommonTable 
                  :fields="item.dataPreview.columns.map(c => ({key: c, label: c}))" 
                  :items="item.dataPreview.rows.map(row => { const obj: any = {}; item.dataPreview!.columns.forEach((c, i) => { obj[c] = row[i] }); return obj; })" 
                  bordered striped hover 
                />
              </div>
            </div>
            <div v-else-if="!item.loading" class="text-gray-400 text-sm py-8 px-4 text-center italic bg-gray-50">
              Run a search to see data for {{ item.name }}.
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
