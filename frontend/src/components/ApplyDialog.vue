<template>
  <dialog ref="dlg" @cancel.prevent>
    <div class="dialog-header">Apply — Write modlist.txt</div>
    <div class="dialog-body">
      <div v-if="loading" class="text-dim">Loading preview…</div>
      <div v-else-if="error" class="text-error">{{ error }}</div>
      <ol v-else class="preview-list">
        <li v-for="name in preview" :key="name">{{ name }}</li>
      </ol>
    </div>
    <div class="dialog-footer">
      <button @click="$emit('close')">Cancel</button>
      <button class="primary" @click="doWrite" :disabled="loading || !!writeError">
        Write modlist.txt
      </button>
      <span v-if="writeError" class="text-error" style="font-size:11px;">{{ writeError }}</span>
      <span v-if="written" class="text-success" style="font-size:11px;">Done! File written.</span>
    </div>
  </dialog>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useWails } from '../composables/useWails'

defineEmits<{ close: [] }>()

const wails = useWails()
const dlg = ref<HTMLDialogElement>()
const preview = ref<string[]>([])
const loading = ref(true)
const error = ref('')
const writeError = ref('')
const written = ref(false)

onMounted(async () => {
  dlg.value?.showModal()
  try {
    const result = await wails.getApplyPreview()
    preview.value = result.names
  } catch (e: any) {
    error.value = String(e)
  } finally {
    loading.value = false
  }
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
dialog { width: 520px; }
.preview-list {
  max-height: 55vh;
  overflow-y: auto;
  padding-left: 24px;
  font-size: 12px;
  line-height: 1.8;
  font-family: 'Consolas', monospace;
}
</style>
