<template>
  <dialog ref="dlg" class="cg-dialog" @cancel.prevent>
    <div class="dialog-header">
      Conflict Graph — {{ conflicts.length }} conflicting files
      <button class="close-btn" @click="$emit('close')" title="Close"><X :size="16" /></button>
    </div>
    <div class="dialog-body">
      <input type="search" v-model="query" placeholder="Filter by mod name or resource…" class="cg-search" />
      <div class="cg-table-wrap">
        <table>
          <thead>
            <tr>
              <th style="width: 420px">Resource</th>
              <th>Mods</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(c, i) in filtered" :key="i">
              <td class="mono text-dim">{{ c.resource }}</td>
              <td class="text-error">{{ c.mods.join('  ✗  ') }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </dialog>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { X } from 'lucide-vue-next'
import type { main } from '../../wailsjs/go/models'

const props = defineProps<{ conflicts: main.ConflictDTO[] }>()
defineEmits<{ close: [] }>()

const dlg = ref<HTMLDialogElement>()
const query = ref('')

onMounted(() => dlg.value?.showModal())

const filtered = computed(() => {
  const q = query.value.toLowerCase()
  if (!q) return props.conflicts
  return props.conflicts.filter(
    (c) => c.resource.toLowerCase().includes(q) || c.mods.some((m) => m.toLowerCase().includes(q))
  )
})
</script>

<style scoped>
.cg-dialog {
  width: 960px;
  max-width: 95vw;
  max-height: 85vh;
  display: flex;
  flex-direction: column;
  border: 1px solid var(--cp-focus);
  box-shadow: 0 0 30px rgba(0, 229, 255, 0.2);
}
.close-btn {
  background: none;
  border: none;
  color: var(--cp-dim);
  cursor: pointer;
  padding: 0 4px;
  display: flex;
  align-items: center;
  transition:
    color 0.15s,
    text-shadow 0.15s;
}
.close-btn:hover {
  color: var(--cp-focus);
  text-shadow: var(--glow-focus);
  box-shadow: none;
  border: none;
}
.dialog-body {
  display: flex;
  flex-direction: column;
  gap: 8px;
  flex: 1;
  overflow: hidden;
}
.cg-search {
  margin-bottom: 4px;
}
.cg-table-wrap {
  overflow-y: auto;
  flex: 1;
}
table {
  width: 100%;
  border-collapse: collapse;
  table-layout: fixed;
}
th,
td {
  padding: 4px 8px;
  border-bottom: 1px solid var(--cp-border);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
th {
  font-weight: 700;
  font-size: 11px;
  color: var(--cp-focus);
  text-shadow: var(--glow-focus);
  text-transform: uppercase;
  letter-spacing: 0.6px;
  background: #0d0d0d;
  border-bottom: 2px solid var(--cp-focus);
  position: sticky;
  top: 0;
}
</style>
