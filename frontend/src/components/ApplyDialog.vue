<template>
  <dialog ref="dlg" :class="{ maximized }" @cancel.prevent @click="onBackdropClick">
    <div class="dialog-header">
      <span>Apply — Write modlist.txt</span>
      <div class="header-actions">
        <button
          class="view-btn"
          :class="{ active: !showDiff }"
          title="Plain list view"
          @click="showDiff = false"
        >List</button>
        <button
          class="view-btn"
          :class="{ active: showDiff }"
          title="Diff view — shows position changes since last scan"
          @click="showDiff = true"
        >Diff</button>
        <button
          v-if="showDiff"
          class="view-btn"
          :class="{ active: diffOnly }"
          title="Hide unchanged entries"
          @click="diffOnly = !diffOnly"
        >Changed only</button>
        <button
          class="view-btn maximize-btn"
          :title="maximized ? 'Restore' : 'Maximize'"
          @click="maximized = !maximized"
        ><Minimize2 v-if="maximized" :size="13" /><Maximize2 v-else :size="13" /></button>
      </div>
    </div>

    <div class="dialog-body">
      <div v-if="loading" class="text-dim">Loading preview…</div>
      <div v-else-if="error" class="text-error">{{ error }}</div>

      <template v-else>
        <!-- Plain list -->
        <ol v-if="!showDiff" class="preview-list">
          <li v-for="name in names" :key="name">{{ name }}</li>
        </ol>

        <!-- Diff list -->
        <div v-else class="diff-list">
          <div
            v-for="item in visibleDiffItems"
            :key="item.name"
            class="diff-row"
          >
            <span class="diff-pos text-dim">{{ item.newPos }}</span>
            <span
              class="diff-badge"
              :class="{
                'badge-new':  item.isNew,
                'badge-up':   !item.isNew && item.delta! < 0,
                'badge-down': !item.isNew && item.delta! > 0,
                'badge-same': !item.isNew && item.delta === 0,
              }"
            >
              <template v-if="item.isNew">NEW</template>
              <template v-else-if="item.delta! < 0">↑{{ -item.delta! }}</template>
              <template v-else-if="item.delta! > 0">↓{{ item.delta }}</template>
              <template v-else>—</template>
            </span>
            <span class="diff-name">{{ item.name }}</span>
          </div>

          <template v-if="removedItems.length">
            <div class="removed-header">Removed ({{ removedItems.length }})</div>
            <div
              v-for="name in removedItems"
              :key="'rem-' + name"
              class="diff-row diff-removed"
            >
              <span class="diff-pos text-dim">—</span>
              <span class="diff-badge badge-removed">–</span>
              <span class="diff-name">{{ name }}</span>
            </div>
          </template>

          <div
            v-if="diffOnly && visibleDiffItems.length === 0 && removedItems.length === 0"
            class="text-dim no-changes"
          >No changes since last scan.</div>
        </div>
      </template>
    </div>

    <div class="dialog-footer">
      <button @click="$emit('close')">Close</button>
      <button class="primary" @click="doWrite" :disabled="loading || !!writeError">Write modlist.txt</button>
      <span v-if="writeError" class="text-error" style="font-size: 11px">{{ writeError }}</span>
      <span v-if="written" class="text-success" style="font-size: 11px">Done! File written.</span>
    </div>
  </dialog>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { Maximize2, Minimize2 } from 'lucide-vue-next'
import { useWails } from '../composables/useWails'

const emit = defineEmits<{ close: [] }>()

const wails = useWails()
const dlg = ref<HTMLDialogElement>()

const names = ref<string[]>([])
const current = ref<string[]>([])
const loading = ref(true)
const error = ref('')
const writeError = ref('')
const written = ref(false)

const showDiff = ref(readStoredBoolean('cp77-modorder.apply.showDiff'))
const diffOnly = ref(readStoredBoolean('cp77-modorder.apply.diffOnly'))
const maximized = ref(false)

watch(showDiff, value => storeBoolean('cp77-modorder.apply.showDiff', value))
watch(diffOnly, value => storeBoolean('cp77-modorder.apply.diffOnly', value))

function readStoredBoolean(key: string): boolean {
  try {
    return localStorage.getItem(key) === 'true'
  } catch {
    return false
  }
}

function storeBoolean(key: string, value: boolean) {
  try {
    localStorage.setItem(key, String(value))
  } catch {
    // Ignore unavailable storage; the controls still work for this dialog.
  }
}

onMounted(async () => {
  dlg.value?.showModal()
  try {
    const result = await wails.getApplyPreview()
    names.value = result.names
    current.value = result.current ?? []
  } catch (e: any) {
    error.value = String(e)
  } finally {
    loading.value = false
  }
})

function onBackdropClick(e: MouseEvent) {
  if (e.target === dlg.value) emit('close')
}

const oldPositions = computed(() => {
  const map = new Map<string, number>()
  for (const [i, name] of current.value.entries()) {
    map.set(name, i + 1)
  }
  return map
})

const diffItems = computed(() =>
  names.value.map((name, i) => {
    const newPos = i + 1
    const oldPos = oldPositions.value.get(name)
    return {
      name,
      newPos,
      oldPos,
      // negative = moved earlier (higher priority), positive = moved later
      delta: oldPos === undefined ? null : newPos - oldPos,
      isNew: oldPos === undefined,
    }
  })
)

const visibleDiffItems = computed(() => {
  if (!diffOnly.value) return diffItems.value
  return diffItems.value.filter(item => item.isNew || (item.delta !== null && item.delta !== 0))
})

const removedItems = computed(() => {
  const newSet = new Set(names.value)
  return current.value.filter(n => !newSet.has(n))
})

async function doWrite() {
  writeError.value = ''
  written.value = false
  try {
    await wails.writeModlist()
    written.value = true
  } catch (e: any) {
    writeError.value = String(e)
  }
}
</script>

<style scoped>
dialog {
  width: 540px;
  border: 1px solid var(--cp-primary);
  box-shadow: 0 0 30px rgba(240, 192, 0, 0.2);
  transition: width 0.12s, height 0.12s;
}

dialog.maximized {
  width: 90vw;
  height: 90vh;
  display: flex;
  flex-direction: column;
}
dialog.maximized :deep(.dialog-body) {
  flex: 1;
  min-height: 0;
}
dialog.maximized .preview-list,
dialog.maximized .diff-list {
  max-height: 100%;
  height: 100%;
}

/* ---- Header buttons ---- */
.header-actions {
  display: flex;
  gap: 4px;
  align-items: center;
}
.view-btn {
  font-size: 10px;
  padding: 2px 7px;
  background: var(--cp-input);
  border: 1px solid var(--cp-border);
  color: var(--cp-dim);
  cursor: pointer;
  border-radius: 2px;
  line-height: 1.5;
  font-family: inherit;
  letter-spacing: 0.03em;
}
.view-btn:hover {
  color: var(--cp-fg);
  border-color: var(--cp-primary);
}
.view-btn.active {
  color: var(--cp-primary);
  border-color: var(--cp-primary);
}
.maximize-btn {
  padding: 2px 5px;
  margin-left: 4px;
  display: flex;
  align-items: center;
}

/* ---- Plain list ---- */
.preview-list {
  max-height: 55vh;
  overflow-y: auto;
  padding: 8px 8px 8px 28px;
  font-size: 12px;
  line-height: 1.9;
  font-family: 'Consolas', monospace;
  border: 1px solid var(--cp-border);
  background: #0d0d0d;
}
.preview-list li {
  border-bottom: 1px solid var(--cp-border);
  padding: 1px 0;
}
.preview-list li:last-child {
  border-bottom: none;
}

/* ---- Diff list ---- */
.diff-list {
  max-height: 55vh;
  overflow-y: auto;
  font-size: 12px;
  font-family: 'Consolas', monospace;
  border: 1px solid var(--cp-border);
  background: #0d0d0d;
}
.diff-row {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 2px 8px;
  border-bottom: 1px solid var(--cp-border);
  min-height: 22px;
}
.diff-row:last-child {
  border-bottom: none;
}
.diff-row.diff-removed {
  opacity: 0.65;
}
.diff-pos {
  width: 30px;
  text-align: right;
  flex-shrink: 0;
  font-size: 10px;
}
.diff-badge {
  width: 38px;
  text-align: center;
  flex-shrink: 0;
  font-size: 10px;
  font-weight: 700;
}
.badge-new      { color: var(--cp-focus); }
.badge-up       { color: var(--cp-success); }
.badge-down     { color: var(--cp-error); }
.badge-same     { color: var(--cp-dim); font-weight: 400; }
.badge-removed  { color: var(--cp-error); }

.diff-name {
  flex: 1;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.removed-header {
  padding: 5px 8px 2px;
  font-size: 10px;
  color: var(--cp-dim);
  text-transform: uppercase;
  letter-spacing: 0.06em;
  border-top: 1px solid var(--cp-border);
}
.no-changes {
  padding: 12px 8px;
  font-size: 12px;
}
</style>
