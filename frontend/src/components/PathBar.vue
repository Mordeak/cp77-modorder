<template>
  <div class="pathbar">
    <input
      type="text"
      :value="modelValue"
      @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
      @keydown.enter="$emit('scan')"
      placeholder="Paste path to mods folder, e.g. C:/Games/Cyberpunk 2077/archive/pc/mod"
    />
    <button @click="$emit('pick')" title="Browse…">…</button>
    <button class="primary" @click="$emit('scan')" :disabled="scanning">
      {{ scanning ? 'Scanning…' : 'Scan' }}
    </button>
    <span v-if="error" class="msg text-error">{{ error }}</span>
    <span v-else-if="progress" class="msg text-dim">{{ progress }}</span>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  modelValue: string
  scanning: boolean
  progress: string
  error: string
}>()
defineEmits<{
  'update:modelValue': [v: string]
  'scan': []
  'pick': []
}>()
</script>

<style scoped>
.pathbar {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 5px 8px;
  background: var(--cp-bg);
  border-bottom: 1px solid var(--cp-border);
}
.pathbar input { flex: 1; }
.msg { font-size: 11px; white-space: nowrap; }
</style>
