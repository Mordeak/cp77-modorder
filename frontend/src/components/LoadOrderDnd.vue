<template>
  <div class="dnd-wrap">
    <div v-if="loading" class="text-dim" style="padding:8px;">Loading…</div>
    <draggable
      v-else
      v-model="localOrder"
      item-key="name"
      handle=".drag-handle"
      :animation="150"
      @end="onReorder"
    >
      <template #item="{ element }">
        <div
          class="dnd-row"
          :class="element === anchorName ? 'dnd-row--anchor' : 'dnd-row--conflict'"
        >
          <span class="drag-handle">⠿</span>
          <span class="dnd-name" :title="element">{{ truncate(element, 36) }}</span>
        </div>
      </template>
    </draggable>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import draggable from 'vuedraggable'
import { useWails } from '../composables/useWails'
import { truncate } from '../utils'

const props = defineProps<{ anchorName: string }>()

const wails = useWails()
const localOrder = ref<string[]>([])
const loading = ref(true)

async function load() {
  loading.value = true
  try {
    const group = await wails.getConflictGroup(props.anchorName)
    localOrder.value = group.mods
  } finally {
    loading.value = false
  }
}

onMounted(load)
watch(() => props.anchorName, load)

async function onReorder() {
  await wails.reorderGroup(localOrder.value)
}
</script>

<style scoped>
.dnd-wrap { user-select: none; }
.dnd-row {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 8px;
  border-bottom: 1px solid var(--cp-border);
  cursor: default;
}
.dnd-row--anchor { background: var(--cp-dnd-anchor); }
.dnd-row--conflict { }
.drag-handle {
  color: var(--cp-dim);
  cursor: grab;
  font-size: 16px;
  flex-shrink: 0;
}
.drag-handle:active { cursor: grabbing; }
.dnd-name { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.dnd-row--conflict .dnd-name { color: var(--cp-error); }
.dnd-row--anchor .dnd-name   { color: var(--cp-fg); }
.sortable-drag .dnd-row { background: var(--cp-dnd-drag) !important; }
</style>
