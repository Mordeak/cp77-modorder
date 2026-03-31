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

    <div
      class="mo2-toggle"
      :class="{ active: modStructure === 'MO2' }"
      @click="$emit('toggle-mo2')"
      title="Toggle MO2 Mod Organizer 2 support (experimental)"
    >
      <span class="mo2-tag">MO2 <em>experimental</em></span>
      <div class="toggle-track">
        <div class="toggle-thumb"></div>
      </div>
      <span class="btn-label">{{ modStructure === 'MO2' ? 'MO2' : 'Default' }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{ hasResult: boolean; modStructure: 'default' | 'MO2' }>()
defineEmits<{
  'open-folder': []
  rescan: []
  apply: []
  conflicts: []
  group: []
  restore: []
  'toggle-mo2': []
}>()
</script>

<style scoped>
.toolbar {
  display: flex;
  align-items: center;
  gap: 2px;
  padding: 5px 8px;
  background: #0d0d0d;
  border-bottom: 2px solid var(--cp-primary);
  box-shadow: 0 2px 12px rgba(240, 192, 0, 0.15);
}
.toolbar button {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 3px;
  padding: 6px 10px;
  border: 1px solid transparent;
  background: transparent;
  text-transform: none;
  letter-spacing: 0;
  transition:
    border-color 0.15s,
    box-shadow 0.15s;
}
.toolbar button:hover {
  border-color: var(--cp-primary);
  box-shadow: var(--glow-primary);
  color: var(--cp-fg);
}
.toolbar button:disabled:hover {
  border-color: transparent;
  box-shadow: none;
}
.btn-icon {
  font-size: 15px;
  line-height: 1;
}
.btn-label {
  font-size: 10px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  color: var(--cp-dim);
}
.toolbar button:hover .btn-label {
  color: var(--cp-fg);
}
.toolbar button:disabled .btn-label {
  opacity: 0.4;
}

/* MO2 toggle — pushed to the far right */
.mo2-toggle {
  margin-left: auto;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  padding: 4px 10px;
  border: 1px solid transparent;
  cursor: pointer;
  user-select: none;
  border-left: 1px solid var(--cp-border);
  transition: border-color 0.15s;
}
.mo2-toggle:hover {
  border-color: var(--cp-focus);
}
.mo2-toggle.active {
  border-color: var(--cp-focus);
  box-shadow: 0 0 6px rgba(0, 229, 255, 0.25);
}

.mo2-tag {
  font-size: 9px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  color: var(--cp-dim);
  white-space: nowrap;
}
.mo2-tag em {
  font-style: normal;
  color: var(--cp-focus);
  opacity: 0.7;
}
.mo2-toggle.active .mo2-tag {
  color: var(--cp-fg);
}

.toggle-track {
  width: 28px;
  height: 13px;
  border-radius: 7px;
  background: var(--cp-input);
  border: 1px solid #444;
  position: relative;
  transition:
    background 0.2s,
    border-color 0.2s;
}
.mo2-toggle.active .toggle-track {
  background: rgba(0, 229, 255, 0.15);
  border-color: var(--cp-focus);
}
.toggle-thumb {
  width: 9px;
  height: 9px;
  border-radius: 50%;
  background: #555;
  position: absolute;
  top: 1px;
  left: 2px;
  transition:
    transform 0.2s,
    background 0.2s;
}
.mo2-toggle.active .toggle-thumb {
  transform: translateX(14px);
  background: var(--cp-focus);
}

.mo2-toggle .btn-label {
  color: var(--cp-dim);
}
.mo2-toggle.active .btn-label {
  color: var(--cp-focus);
}
</style>
