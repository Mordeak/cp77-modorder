<template>
  <dialog ref="dlg" class="crd-dialog" @cancel.prevent>
    <div class="dialog-header">
      New Conflicting Mods — {{ stepIndex + 1 }} of {{ props.mods.length }}
    </div>
    <div class="dialog-body">
      <p class="crd-hint">
        <span class="text-error">{{ currentMod }}</span> conflicts with existing mods.<br/>
        Drag to set the load order — top loads first and wins conflicts.
      </p>
      <LoadOrderDnd :anchor-name="currentMod" />
    </div>
    <div class="dialog-footer">
      <button @click="advance">{{ isLastStep ? 'Done' : 'Skip' }}</button>
      <button v-if="!isLastStep" class="primary" @click="advance">Next →</button>
    </div>
  </dialog>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import LoadOrderDnd from './LoadOrderDnd.vue'

const props = defineProps<{ mods: string[] }>()
const emit = defineEmits<{ close: [] }>()

const dlg = ref<HTMLDialogElement>()
const stepIndex = ref(0)

const currentMod = computed(() => props.mods[stepIndex.value])
const isLastStep = computed(() => stepIndex.value === props.mods.length - 1)

onMounted(() => dlg.value?.showModal())

function advance() {
  if (isLastStep.value) {
    dlg.value?.close()
    emit('close')
  } else {
    stepIndex.value++
  }
}
</script>

<style scoped>
.crd-dialog { width: 480px; }
.crd-hint {
  font-size: 12px;
  color: var(--cp-dim);
  line-height: 1.6;
  margin-bottom: 10px;
}
</style>
