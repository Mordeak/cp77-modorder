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
  background: #0d0d0d;
  border-top: 2px solid var(--cp-border);
  font-size: 10px;
  color: var(--cp-dim);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}
</style>
