<template>
  <div class="table-outer">
    <div class="table-controls">
      <input v-model="searchQuery" class="search-input" type="search" placeholder="Search archive…" />
      <label class="filter-toggle">
        <input type="checkbox" v-model="showLosingOnly" />
        Display only conflicts
      </label>
    </div>
    <div class="table-wrap">
      <table>
        <thead>
          <tr>
            <th style="width: 24px"></th>
            <th style="width: 32px">#</th>
            <th style="width: 390px">Mod</th>
            <th style="width: 55px">Files</th>
            <th style="width: 80px">Conflicts</th>
            <th style="width: 80px">W / L</th>
            <th style="width: 65px">Status</th>
          </tr>
        </thead>
        <draggable
          v-if="!isFiltered"
          v-model="draggableRows"
          tag="tbody"
          :item-key="(el: RowEntry) => el.row.name"
          handle=".drag-handle"
          :animation="150"
          @end="onReorder"
        >
          <template #item="{ element: { row, idx } }">
            <tr :class="rowClass(row, idx)" @click="row.missing ? null : $emit('select', idx)">
              <td><GripVertical class="drag-handle" :class="{ 'drag-handle--hidden': row.missing }" :size="14" /></td>
              <td>{{ row.missing ? '—' : idx + 1 }}</td>
              <td :title="row.name">{{ row.name }}</td>
              <td>{{ row.mod ? row.mod.fileCount : '—' }}</td>
              <td :class="row.mod && row.mod.conflictCount > 0 ? 'text-error' : ''">
                {{ row.mod ? row.mod.conflictCount : '—' }}
              </td>
              <td v-if="row.mod">
                <span class="text-success">{{ row.mod.wins }}</span>
                {{ ' / ' }}
                <span class="text-error">{{ row.mod.losses }}</span>
              </td>
              <td v-else>—</td>
              <td>
                <span v-if="row.missing" class="badge badge-missing">MISSING</span>
                <span v-else-if="row.unlisted" class="badge badge-new">NEW</span>
              </td>
            </tr>
          </template>
        </draggable>
        <tbody v-else>
          <tr
            v-for="{ row, idx } in filteredRows"
            :key="row.name"
            :class="rowClass(row, idx)"
            @click="row.missing ? null : $emit('select', idx)"
          >
            <td><GripVertical class="drag-handle drag-handle--hidden" :size="14" /></td>
            <td>{{ row.missing ? '—' : idx + 1 }}</td>
            <td :title="row.name">{{ row.name }}</td>
            <td>{{ row.mod ? row.mod.fileCount : '—' }}</td>
            <td :class="row.mod && row.mod.conflictCount > 0 ? 'text-error' : ''">
              {{ row.mod ? row.mod.conflictCount : '—' }}
            </td>
            <td v-if="row.mod">
              <span class="text-success">{{ row.mod.wins }}</span>
              {{ ' / ' }}
              <span class="text-error">{{ row.mod.losses }}</span>
            </td>
            <td v-else>—</td>
            <td>
              <span v-if="row.missing" class="badge badge-missing">MISSING</span>
              <span v-else-if="row.unlisted" class="badge badge-new">NEW</span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import draggable from 'vuedraggable'
import { GripVertical } from 'lucide-vue-next'
import type { main } from '../../wailsjs/go/models'
import { useWails } from '../composables/useWails'

type RowEntry = { row: main.DisplayRowDTO; idx: number }

const searchQuery = ref('')
const showLosingOnly = ref(false)

const props = defineProps<{
  rows: main.DisplayRowDTO[]
  selectedIndex: number
}>()
defineEmits<{ select: [idx: number] }>()

const wails = useWails()

const draggableRows = ref<RowEntry[]>([])
watch(
  () => props.rows,
  (rows) => {
    draggableRows.value = rows.map((row, idx) => ({ row, idx }))
  },
  { immediate: true }
)

const isFiltered = computed(() => searchQuery.value !== '' || showLosingOnly.value)

const filteredRows = computed(() => {
  const q = searchQuery.value.toLowerCase()
  return draggableRows.value.filter(({ row }) => {
    if (q && !row.name.toLowerCase().includes(q)) return false
    if (showLosingOnly.value && (row.mod?.losses ?? 0) === 0) return false
    return true
  })
})

async function onReorder() {
  if (isFiltered.value) return
  await wails.setModlistOrder(draggableRows.value.map((e) => e.row.name))
}

function rowClass(row: main.DisplayRowDTO, idx: number) {
  const losses = row.mod?.losses ?? 0
  const wins = row.mod?.wins ?? 0
  return {
    'row-selected': idx === props.selectedIndex && !row.missing,
    'row-missing': row.missing,
    'row-new': row.unlisted && !row.missing,
    'row-losing-all': !row.missing && !row.unlisted && losses > 0 && wins === 0,
    'row-conflict': !row.missing && !row.unlisted && losses > 0 && wins > 0,
    'row-clickable': !row.missing,
  }
}
</script>

<style scoped>
.table-outer {
  display: flex;
  flex-direction: column;
  height: 100%;
}
.table-wrap {
  flex: 1;
  overflow-y: auto;
}
.table-controls {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 5px 8px;
  background: #0d0d0d;
  border-bottom: 1px solid var(--cp-border);
  flex-shrink: 0;
}
.search-input {
  flex: 1;
  background: var(--cp-input);
  color: var(--cp-fg);
  border: 1px solid var(--cp-border);
  border-radius: 0;
  padding: 3px 7px;
  font-size: 12px;
  outline: none;
  transition:
    border-color 0.15s,
    box-shadow 0.15s;
}
.search-input:focus {
  border-color: var(--cp-focus);
  box-shadow: var(--glow-focus);
}
.filter-toggle {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 11px;
  color: var(--cp-dim);
  cursor: pointer;
  white-space: nowrap;
  text-transform: uppercase;
  letter-spacing: 0.4px;
}
table {
  width: 100%;
  border-collapse: collapse;
  table-layout: fixed;
}
thead tr {
  background: #0d0d0d;
  border-bottom: 2px solid var(--cp-primary);
  position: sticky;
  top: 0;
}
th,
td {
  padding: 4px 6px;
  text-align: left;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  border-bottom: 1px solid var(--cp-border);
}
th {
  font-weight: 700;
  font-size: 11px;
  color: var(--cp-primary);
  text-transform: uppercase;
  letter-spacing: 0.6px;
  text-shadow: var(--glow-primary);
}

.drag-handle {
  color: var(--cp-dim);
  cursor: grab;
  display: block;
  transition: color 0.15s;
}
.drag-handle:hover {
  color: var(--cp-focus);
}
.drag-handle:active {
  cursor: grabbing;
  color: var(--cp-primary);
}
.drag-handle--hidden {
  visibility: hidden;
  cursor: default;
}

.row-clickable {
  cursor: pointer;
}
.row-clickable:hover td {
  background: rgba(255, 255, 255, 0.03);
}
.row-selected td {
  background: var(--cp-row-selected) !important;
  border-left: 2px solid var(--cp-focus) !important;
}
.row-losing-all td {
  color: var(--cp-dim);
  border-left: 2px solid var(--cp-dim);
}
.row-conflict td {
  background: var(--cp-row-conflict);
  border-left: 2px solid var(--cp-error);
}
.row-missing td {
  background: var(--cp-row-missing);
  color: var(--cp-dim);
  cursor: default;
}
.row-new td {
  background: var(--cp-row-new);
  border-left: 2px solid var(--cp-primary);
}

.badge {
  font-size: 10px;
  padding: 1px 5px;
  border-radius: 0;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.3px;
}
.badge-missing {
  background: rgba(255, 51, 102, 0.2);
  color: var(--cp-error);
  border: 1px solid var(--cp-error);
  box-shadow: var(--glow-error);
}
.badge-new {
  background: rgba(240, 192, 0, 0.15);
  color: var(--cp-primary);
  border: 1px solid var(--cp-primary);
  box-shadow: var(--glow-primary);
}
</style>
