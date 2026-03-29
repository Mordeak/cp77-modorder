<template>
  <div class="toolbar">
    <button @click="$emit('open-folder')" title="Open folder">📁</button>
    <button @click="$emit('rescan')" title="Rescan">🔄</button>
    <button @click="$emit('apply')" :disabled="!hasResult" title="Apply — write modlist.txt">💾</button>
    <button @click="$emit('conflicts')" :disabled="!hasResult" title="Conflict graph">ℹ️</button>
  </div>
</template>

<script setup lang="ts">
defineProps<{ hasResult: boolean }>()
defineEmits<{
  'open-folder': []
  'rescan': []
  'apply': []
  'conflicts': []
}>()
</script>

<style scoped>
.toolbar {
  display: flex;
  gap: 4px;
  padding: 6px 8px;
  background: var(--cp-bg2);
  border-bottom: 1px solid var(--cp-border);
}
.toolbar button {
  font-size: 16px;
  padding: 3px 8px;
  border: none;
  background: transparent;
}
.toolbar button:hover { background: var(--cp-border); }
</style>
