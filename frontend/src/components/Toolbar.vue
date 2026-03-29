<template>
  <div class="toolbar">
    <button @click="$emit('open-folder')" title="Open folder">📁</button>
    <button @click="$emit('rescan')" title="Rescan">🔄</button>
    <button @click="$emit('apply')" :disabled="!hasResult" title="Apply — write modlist.txt">💾</button>
    <button @click="$emit('conflicts')" :disabled="!hasResult" title="Conflict graph">ℹ️</button>
    <button @click="$emit('group')" :disabled="!hasResult" title="Group related mods together">🔗</button>
  </div>
</template>

<script setup lang="ts">
defineProps<{ hasResult: boolean }>()
defineEmits<{
  'open-folder': []
  'rescan': []
  'apply': []
  'conflicts': []
  'group': []
}>()
</script>

<style scoped>
.toolbar {
  display: flex;
  gap: 2px;
  padding: 5px 8px;
  background: #0D0D0D;
  border-bottom: 2px solid var(--cp-primary);
  box-shadow: 0 2px 12px rgba(240,192,0,0.15);
}
.toolbar button {
  font-size: 15px;
  padding: 4px 10px;
  border: 1px solid transparent;
  background: transparent;
  text-transform: none;
  letter-spacing: 0;
  transition: border-color 0.15s, box-shadow 0.15s;
}
.toolbar button:hover {
  border-color: var(--cp-primary);
  box-shadow: var(--glow-primary);
  color: var(--cp-fg);
}
.toolbar button:disabled:hover { border-color: transparent; box-shadow: none; }
</style>
