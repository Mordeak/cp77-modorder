<template>
  <div class="statusbar">{{ text }}</div>
</template>

<script setup lang="ts">
defineProps<{ text: string }>()
</script>

<style scoped>
.statusbar {
  flex-shrink: 0;
  padding: 3px 10px;
  background: var(--cp-bg2);
  border-top: 1px solid var(--cp-border);
  font-size: 11px;
  color: var(--cp-dim);
}
</style>
