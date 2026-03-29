<template>
  <div class="toolbar">
    <button @click="$emit('open-folder')" title="Open folder">
      <span class="btn-icon">📁</span><span class="btn-label">Open</span>
    </button>
    <button @click="$emit('rescan')" title="Rescan">
      <span class="btn-icon">🔄</span><span class="btn-label">Rescan</span>
    </button>
    <button @click="$emit('apply')" :disabled="!hasResult" title="Apply — write modlist.txt">
      <span class="btn-icon">💾</span><span class="btn-label">Apply</span>
    </button>
    <button @click="$emit('conflicts')" :disabled="!hasResult" title="Conflict graph">
      <span class="btn-icon">ℹ️</span><span class="btn-label">Conflicts</span>
    </button>
    <button @click="$emit('group')" :disabled="!hasResult" title="Group related mods together">
      <span class="btn-icon">🔗</span><span class="btn-label">Group</span>
    </button>
    <button @click="$emit('restore')" title="Restore a modlist backup">
      <span class="btn-icon">↩️</span><span class="btn-label">Restore</span>
    </button>
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
  'restore': []
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
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
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
.btn-icon { font-size: 15px; line-height: 1; }
.btn-label { font-size: 10px; text-transform: uppercase; letter-spacing: 0.5px; color: var(--cp-dim); }
.toolbar button:hover .btn-label { color: var(--cp-fg); }
.toolbar button:disabled .btn-label { opacity: 0.4; }
</style>
