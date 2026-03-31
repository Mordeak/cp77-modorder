<template>
  <span v-if="priority > 0" class="badge-prio">{{ priority }}</span>
  <span v-else class="badge-unset">—</span>
</template>

<script setup lang="ts">
defineProps<{ priority: number }>()
</script>

<style scoped>
.badge-prio {
  background: rgba(240, 192, 0, 0.12);
  color: var(--cp-primary);
  font-size: 10px;
  padding: 1px 5px;
  border-radius: 0;
  font-weight: 700;
  border: 1px solid var(--cp-primary);
  box-shadow: var(--glow-primary);
  letter-spacing: 0.3px;
}
.badge-unset {
  color: var(--cp-dim);
}
</style>
