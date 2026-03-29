<template>
  <dialog ref="dlg" @cancel.prevent>
    <div class="dialog-header">Restore Backup</div>
    <div class="dialog-body">
      <div v-if="loading" class="text-dim">Loading backups…</div>
      <div v-else-if="error" class="text-error">{{ error }}</div>
      <div v-else-if="backups.length === 0" class="text-dim">No backups found in modlist.old/.</div>
      <ul v-else class="backup-list">
        <li
          v-for="name in backups"
          :key="name"
          class="backup-entry"
          :class="{ 'backup-entry--selected': name === selected }"
          @click="selected = name"
          @dblclick="doRestore"
        >{{ name }}</li>
      </ul>
    </div>
    <div class="dialog-footer">
      <button @click="$emit('close')">Cancel</button>
      <button class="primary" @click="doRestore" :disabled="!selected || restoring">
        Restore
      </button>
      <span v-if="restoreError" class="text-error" style="font-size:11px;">{{ restoreError }}</span>
    </div>
  </dialog>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useWails } from '../composables/useWails'

const emit = defineEmits<{ close: [] }>()

const wails = useWails()
const dlg = ref<HTMLDialogElement>()
const backups = ref<string[]>([])
const selected = ref('')
const loading = ref(true)
const error = ref('')
const restoring = ref(false)
const restoreError = ref('')

onMounted(async () => {
  dlg.value?.showModal()
  try {
    backups.value = await wails.listBackups()
  } catch (e: any) {
    error.value = String(e)
  } finally {
    loading.value = false
  }
})

async function doRestore() {
  if (!selected.value) return
  restoring.value = true
  restoreError.value = ''
  try {
    await wails.restoreBackup(selected.value)
    emit('close')
  } catch (e: any) {
    restoreError.value = String(e)
  } finally {
    restoring.value = false
  }
}
</script>

<style scoped>
dialog {
  width: 520px;
  border: 1px solid var(--cp-primary);
  box-shadow: 0 0 30px rgba(240,192,0,0.2);
}
.backup-list {
  list-style: none;
  padding: 0;
  margin: 0;
  max-height: 55vh;
  overflow-y: auto;
  border: 1px solid var(--cp-border);
  background: #0D0D0D;
}
.backup-entry {
  padding: 6px 10px;
  font-size: 12px;
  font-family: 'Consolas', monospace;
  border-bottom: 1px solid var(--cp-border);
  cursor: pointer;
  color: var(--cp-dim);
}
.backup-entry:last-child { border-bottom: none; }
.backup-entry:hover { background: rgba(255,255,255,0.03); color: var(--cp-fg); }
.backup-entry--selected {
  background: var(--cp-row-selected) !important;
  border-left: 2px solid var(--cp-focus);
  color: var(--cp-fg);
}
</style>
