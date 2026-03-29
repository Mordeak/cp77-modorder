<template>
  <div class="detail">
    <div class="detail-name">{{ mod.name }}</div>

    <div class="detail-stats">
      <div class="stat-row"><span class="stat-label">Files</span><span>{{ mod.fileCount }}</span></div>
      <div class="stat-row">
        <span class="stat-label">Conflicting files</span><span>{{ mod.conflictCount }}</span>
      </div>
      <div class="stat-row"><span class="stat-label">Wins</span><span class="text-success">{{ mod.wins }}</span></div>
      <div class="stat-row"><span class="stat-label">Losses</span><span class="text-error">{{ mod.losses }}</span></div>
      <div class="stat-row">
        <span class="stat-label">Priority</span>
        <PriorityBadge :priority="mod.priority" />
      </div>
    </div>

    <div class="detail-actions">
      <button @click="openPrioForm">Set Priority</button>
      <button @click="clearPriority" :disabled="mod.priority === 0">Clear Priority</button>
    </div>

    <div v-if="mod.conflictsWith?.length" class="detail-section section-conflicts">
      <div class="section-title section-title--conflict">Conflicts</div>
      <div class="conflict-list">
        <div
          v-for="entry in groupedConflicts"
          :key="entry.opponent"
          class="conflict-entry"
        >
          <span class="text-error">{{ entry.opponent }}</span>
          <span class="conflict-count">{{ entry.count }}</span>
        </div>
        <div v-if="mod.hasMore" class="text-dim" style="padding: 4px 0;">
          …and {{ mod.moreCount }} more
        </div>
      </div>
    </div>

    <div v-if="mod.conflictCount > 0" class="detail-section section-loadorder">
      <div class="section-title section-title--loadorder">Load Order — drag to reorder</div>
      <LoadOrderDnd :anchor-name="mod.name" />
    </div>

    <!-- Priority form -->
    <dialog ref="prioDialog" class="prio-dialog">
      <div class="dialog-header">Set Priority — {{ mod.name }}</div>
      <div class="dialog-body">
        <input
          ref="prioInput"
          type="text"
          v-model="prioText"
          placeholder="1–99, 0 = unset"
          @keydown.enter="applyPriority"
        />
      </div>
      <div class="dialog-footer">
        <button @click="prioDialog!.close()">Cancel</button>
        <button class="primary" @click="applyPriority">Apply</button>
      </div>
    </dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import type { main } from '../../wailsjs/go/models'
import { useWails } from '../composables/useWails'
import PriorityBadge from './PriorityBadge.vue'
import LoadOrderDnd from './LoadOrderDnd.vue'

const props = defineProps<{
  mod: main.ModDTO
  selectedIndex: number
  rows: main.DisplayRowDTO[]
}>()

const wails = useWails()

const prioDialog = ref<HTMLDialogElement>()
const prioInput = ref<HTMLInputElement>()
const prioText = ref('')

const groupedConflicts = computed(() => {
  const map = new Map<string, number>()
  for (const pair of props.mod.conflictsWith ?? []) {
    map.set(pair.opponent, (map.get(pair.opponent) ?? 0) + 1)
  }
  return Array.from(map.entries()).map(([opponent, count]) => ({ opponent, count }))
})

function openPrioForm() {
  prioText.value = props.mod.priority > 0 ? String(props.mod.priority) : ''
  prioDialog.value?.showModal()
  setTimeout(() => prioInput.value?.focus(), 50)
}

async function applyPriority() {
  const p = parseInt(prioText.value)
  if (isNaN(p) || p < 0 || p > 99) return
  prioDialog.value?.close()
  await wails.setPriority(props.mod.name, p)
}

async function clearPriority() {
  await wails.setPriority(props.mod.name, 0)
}
</script>

<style scoped>
.detail {
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.detail-name {
  font-size: 15px;
  font-weight: 700;
  color: var(--cp-primary);
  word-break: break-all;
}
.detail-stats {
  display: flex;
  flex-direction: column;
  gap: 3px;
}
.stat-row {
  display: flex;
  gap: 8px;
  font-size: 12px;
}
.stat-label {
  color: var(--cp-dim);
  width: 130px;
  flex-shrink: 0;
}
.detail-actions { display: flex; gap: 8px; }
.detail-section {
  display: flex;
  flex-direction: column;
  gap: 4px;
  border-radius: 4px;
  padding: 8px;
}
.section-conflicts {
  background: rgba(255, 51, 102, 0.05);
  border: 1px solid rgba(255, 51, 102, 0.2);
}
.section-loadorder {
  background: rgba(0, 229, 255, 0.05);
  border: 1px solid rgba(0, 229, 255, 0.15);
}
.section-title {
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  padding-bottom: 5px;
  margin-bottom: 2px;
}
.section-title--conflict {
  color: var(--cp-error);
  border-bottom: 1px solid rgba(255, 51, 102, 0.3);
}
.section-title--loadorder {
  color: var(--cp-focus);
  border-bottom: 1px solid rgba(0, 229, 255, 0.2);
}
.conflict-list { max-height: 260px; overflow-y: auto; }
.conflict-entry {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 4px 0;
  border-bottom: 1px solid var(--cp-border);
}
.conflict-count {
  font-size: 10px;
  font-weight: 600;
  padding: 1px 6px;
  border-radius: 10px;
  background: rgba(255, 51, 102, 0.25);
  color: var(--cp-error);
  flex-shrink: 0;
}
.prio-dialog { min-width: 300px; }
.prio-dialog input { width: 100%; }
</style>
