<template>
  <div class="table-outer">
    <div class="table-controls">
      <input
        v-model="searchQuery"
        class="search-input"
        type="search"
        placeholder="Search archive…"
      />
      <label class="filter-toggle">
        <input type="checkbox" v-model="showLosingOnly" />
        Losing only
      </label>
    </div>
    <div class="table-wrap">
    <table>
      <thead>
        <tr>
          <th style="width:40px">#</th>
          <th style="width:70px">Priority</th>
          <th style="width:320px">Mod</th>
          <th style="width:55px">Files</th>
          <th style="width:80px">Conflicts</th>
          <th style="width:80px">W / L</th>
          <th style="width:65px">Status</th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="({ row, idx }) in filteredRows"
          :key="row.name + idx"
          :class="rowClass(row, idx)"
          @click="row.missing ? null : $emit('select', idx)"
        >
          <td>{{ row.missing ? '—' : idx + 1 }}</td>
          <td>
            <PriorityBadge v-if="row.mod" :priority="row.mod.priority" />
            <span v-else>—</span>
          </td>
          <td :title="row.name">{{ truncate(row.name, 43) }}</td>
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
import { ref, computed } from 'vue'
import type { main } from '../../wailsjs/go/models'
import { truncate } from '../utils'
import PriorityBadge from './PriorityBadge.vue'

const searchQuery = ref('')
const showLosingOnly = ref(false)

const props = defineProps<{
  rows: main.DisplayRowDTO[]
  selectedIndex: number
}>()
defineEmits<{ select: [idx: number] }>()

const filteredRows = computed(() => {
  const q = searchQuery.value.toLowerCase()
  return props.rows
    .map((row, idx) => ({ row, idx }))
    .filter(({ row }) => {
      if (q && !row.name.toLowerCase().includes(q)) return false
      if (showLosingOnly.value && (row.mod?.losses ?? 0) === 0) return false
      return true
    })
})

function rowClass(row: main.DisplayRowDTO, idx: number) {
  return {
    'row-selected': idx === props.selectedIndex && !row.missing,
    'row-missing':  row.missing,
    'row-new':      row.unlisted && !row.missing,
    'row-conflict': !row.missing && !row.unlisted && (row.mod?.losses ?? 0) > 0,
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
  padding: 6px 8px;
  background: var(--cp-bg2);
  border-bottom: 1px solid var(--cp-border);
  flex-shrink: 0;
}
.search-input {
  flex: 1;
  background: var(--cp-bg);
  color: var(--cp-fg);
  border: 1px solid var(--cp-border);
  border-radius: 3px;
  padding: 3px 7px;
  font-size: 12px;
  outline: none;
}
.search-input:focus { border-color: var(--cp-focus); }
.filter-toggle {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 12px;
  color: var(--cp-dim);
  cursor: pointer;
  white-space: nowrap;
}
table {
  width: 100%;
  border-collapse: collapse;
  table-layout: fixed;
}
thead tr {
  background: var(--cp-bg2);
  position: sticky;
  top: 0;
}
th, td {
  padding: 4px 6px;
  text-align: left;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  border-bottom: 1px solid var(--cp-border);
}
th { font-weight: 600; font-size: 12px; color: var(--cp-primary); }

.row-clickable { cursor: pointer; }
.row-clickable:hover td { background: rgba(255,255,255,0.04); }
.row-selected td { background: var(--cp-row-selected) !important; }
.row-conflict td { background: var(--cp-row-conflict); }
.row-missing td  { background: var(--cp-row-missing); color: var(--cp-dim); cursor: default; }
.row-new td      { background: var(--cp-row-new); }

.badge {
  font-size: 10px;
  padding: 1px 5px;
  border-radius: 2px;
  font-weight: 600;
}
.badge-missing { background: rgba(255,51,102,0.3); color: var(--cp-error); }
.badge-new     { background: rgba(240,192,0,0.25); color: var(--cp-primary); }
</style>
