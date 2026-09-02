import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'
import { http } from '../api/http'

export interface Connection {
  id: number
  name: string
  driver: string
  host: string
  port: number
  database: string
  ssl_mode: string
  read_username: string
  write_username: string
  allowed_schemas: string[]
}

export interface TableMeta {
  schema: string
  name: string
  kind: string
}

export const useDatabaseStore = defineStore('database', () => {
  const connections = ref<Connection[]>([])
  const selectedConnId = ref<number | null>(null)
  const databases = ref<string[]>([])
  const selectedPhysicalDatabase = ref('')
  const catalogRoles = ref<string[]>([])
  const tables = ref<TableMeta[]>([])
  const selectedSchema = ref('public')

  const currentConnection = computed(() => connections.value.find((x) => x.id === selectedConnId.value) ?? null)

  const effectiveSchema = computed(() => {
    const c = currentConnection.value
    if (!c) return selectedSchema.value
    if (c.driver === 'mysql') {
      const db = (selectedPhysicalDatabase.value || c.database || '').trim()
      return db || selectedSchema.value
    }
    return selectedSchema.value
  })

  function dbParams(): Record<string, string> {
    const o: Record<string, string> = {}
    if (selectedPhysicalDatabase.value) o.database = selectedPhysicalDatabase.value
    return o
  }

  async function loadConnections() {
    try {
      const { data } = await http.get<{ connections: Connection[] }>('/api/connections')
      connections.value = data.connections || []
      if (connections.value.length === 0) {
        selectedConnId.value = null
        return
      }
      if (!selectedConnId.value) {
        selectedConnId.value = connections.value[0].id
        selectedPhysicalDatabase.value = connections.value[0].database || ''
      }
    } catch (e) {
      console.error('Failed to load connections', e)
    }
  }

  async function loadDatabases() {
    if (!selectedConnId.value) return
    try {
      const { data } = await http.get<{ databases: string[] }>(
        `/api/connections/${selectedConnId.value}/databases`,
        { params: dbParams() }
      )
      databases.value = data.databases || []
    } catch (e) {
      databases.value = []
    }
  }

  async function loadCatalogRoles() {
    if (!selectedConnId.value) return
    try {
      const { data } = await http.get<{ roles: string[] }>(
        `/api/connections/${selectedConnId.value}/catalog/roles`,
        { params: dbParams() },
      )
      catalogRoles.value = data.roles || []
    } catch (e: unknown) {
      catalogRoles.value = []
    }
  }

  async function loadTables() {
    if (!selectedConnId.value) return
    try {
      const { data } = await http.get<{ tables: TableMeta[] }>(
        `/api/connections/${selectedConnId.value}/tables`,
        { params: { schema: effectiveSchema.value, ...dbParams() } },
      )
      tables.value = data.tables || []
    } catch (e) {
      tables.value = []
    }
  }

  // Set up global watchers
  watch(selectedConnId, async () => {
    databases.value = []
    tables.value = []
    catalogRoles.value = []
    if (selectedConnId.value) {
      const c = connections.value.find((x) => x.id === selectedConnId.value)
      if (c?.database) selectedPhysicalDatabase.value = c.database
      await loadDatabases()
      if (databases.value.length && !databases.value.includes(selectedPhysicalDatabase.value)) {
        selectedPhysicalDatabase.value = databases.value[0] || ''
      }
      await loadTables()
      await loadCatalogRoles()
    }
  })

  watch([selectedConnId, selectedSchema], () => {
    void loadTables()
  })

  watch(selectedPhysicalDatabase, async (next, prev) => {
    if (!selectedConnId.value) return
    if (next !== prev) {
      await loadTables()
      await loadCatalogRoles()
    }
  })

  return {
    connections,
    selectedConnId,
    databases,
    selectedPhysicalDatabase,
    catalogRoles,
    tables,
    selectedSchema,
    currentConnection,
    effectiveSchema,
    dbParams,
    loadConnections,
    loadDatabases,
    loadCatalogRoles,
    loadTables
  }
})
