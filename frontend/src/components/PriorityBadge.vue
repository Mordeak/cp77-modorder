<template>
  <span v-if="priority > 0" class="badge-prio">{{ priority }}</span>
  <span v-else class="badge-unset">—</span>
</template>

<script setup lang="ts">
defineProps<{ priority: number }>()
</script>

<style scoped>
.badge-prio {
  background: rgba(240,192,0,0.25);
  color: var(--cp-primary);
  font-size: 11px;
  padding: 1px 6px;
  border-radius: 10px;
  font-weight: 600;
}
.badge-unset { color: var(--cp-dim); }
</style>
