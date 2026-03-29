<template>
  <dialog ref="dlg" class="cg-dialog" @cancel.prevent>
    <div class="dialog-header">
      Conflict Graph — {{ conflicts.length }} conflicting files
      <button class="close-btn" @click="$emit('close')">✕</button>
    </div>
    <div class="dialog-body">
      <input
        type="search"
        v-model="query"
        placeholder="Filter by mod name or resource…"
        class="cg-search"
      />
      <div class="cg-table-wrap">
        <table>
          <thead>
            <tr>
              <th style="width:420px">Resource</th>
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
import type { main } from '../../wailsjs/go/models'

const props = defineProps<{ conflicts: main.ConflictDTO[] }>()
defineEmits<{ close: [] }>()

const dlg = ref<HTMLDialogElement>()
const query = ref('')

onMounted(() => dlg.value?.showModal())

const filtered = computed(() => {
  const q = query.value.toLowerCase()
  if (!q) return props.conflicts
  return props.conflicts.filter(c =>
    c.resource.toLowerCase().includes(q) ||
    c.mods.some(m => m.toLowerCase().includes(q))
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
}
.dialog-header { display: flex; justify-content: space-between; align-items: center; }
.close-btn { background: none; border: none; color: var(--cp-dim); font-size: 16px; cursor: pointer; padding: 0; }
.close-btn:hover { color: var(--cp-fg); }
.dialog-body { display: flex; flex-direction: column; gap: 8px; flex: 1; overflow: hidden; }
.cg-search { margin-bottom: 4px; }
.cg-table-wrap { overflow-y: auto; flex: 1; }
table { width: 100%; border-collapse: collapse; table-layout: fixed; }
th, td {
  padding: 4px 8px;
  border-bottom: 1px solid var(--cp-border);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
th { font-weight: 600; font-size: 12px; color: var(--cp-primary); background: var(--cp-bg2); position: sticky; top: 0; }
</style>
