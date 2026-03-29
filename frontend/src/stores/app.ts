import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { main } from '../../wailsjs/go/models'

export const useAppStore = defineStore('app', () => {
  // State
  const scanResult = ref<main.ScanResultDTO | null>(null)
  const selectedIndex = ref<number>(-1)
  const modDir = ref('')
  const scanning = ref(false)
  const scanProgress = ref('')
  const scanError = ref('')
  const showConflictGraph = ref(false)
  const showApplyDialog = ref(false)
  const showConflictResolution = ref(false)
  const showRestoreDialog = ref(false)
  const newConflictingMods = ref<string[]>([])
  // MO2
  const modStructure = ref<'default' | 'MO2'>('default')
  const mo2Dir = ref('')
  const mo2Profile = ref('')
  const mo2Profiles = ref<string[]>([])

  // Computed
  const rows = computed(() => scanResult.value?.rows ?? [])
  const conflicts = computed(() => scanResult.value?.conflicts ?? [])
  const summary = computed(() => scanResult.value?.summary ?? '')
  const hasResult = computed(() => scanResult.value !== null)
  const selectedMod = computed(() => {
    if (selectedIndex.value < 0) return null
    const row = rows.value[selectedIndex.value]
    return row?.mod ?? null
  })

  function setScanResult(result: main.ScanResultDTO) {
    scanResult.value = result
    selectedIndex.value = -1
    scanError.value = ''
  }

  // Like setScanResult but re-finds the previously selected mod by name so the
  // detail panel stays open when rows are reordered or stats change.
  function updateScanResult(result: main.ScanResultDTO) {
    const prevName = selectedIndex.value >= 0 ? rows.value[selectedIndex.value]?.name : null
    scanResult.value = result
    if (prevName != null) {
      const newIdx = result.rows.findIndex(r => r.name === prevName)
      selectedIndex.value = newIdx
    }
  }

  function selectRow(idx: number) {
    selectedIndex.value = idx
  }

  return {
    scanResult, selectedIndex, modDir, scanning, scanProgress, scanError,
    showConflictGraph, showApplyDialog, showConflictResolution, showRestoreDialog, newConflictingMods,
    modStructure, mo2Dir, mo2Profile, mo2Profiles,
    rows, conflicts, summary, hasResult, selectedMod,
    setScanResult, updateScanResult, selectRow,
  }
})
