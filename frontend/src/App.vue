<template>
  <div class="app-layout">
    <div class="app-top">
      <Toolbar
        :has-result="store.hasResult"
        :mod-structure="store.modStructure"
        @open-folder="onOpenFolder"
        @rescan="onRescan"
        @apply="store.showApplyDialog = true"
        @conflicts="store.showConflictGraph = true"
        @group="wails.groupConflicts()"
        @restore="store.showRestoreDialog = true"
        @toggle-mo2="store.modStructure = store.modStructure === 'MO2' ? 'default' : 'MO2'"
      />
      <PathBar
        v-model="store.modDir"
        :scanning="store.scanning"
        :progress="store.scanProgress"
        :error="store.scanError"
        :mod-structure="store.modStructure"
        :mo2-dir="store.mo2Dir"
        :mo2-profile="store.mo2Profile"
        :mo2-profiles="store.mo2Profiles"
        @scan="onScan"
        @pick="onPick"
        @update:mo2-dir="store.mo2Dir = $event"
        @update:mo2-profile="store.mo2Profile = $event"
        @pick-m-o2="onPickMO2"
        @scan-m-o2="onScanMO2"
        @load-profiles="wails.getMO2Profiles($event)"
      />
    </div>

    <div class="app-split">
      <div class="app-left">
        <ModListTable :rows="store.rows" :selected-index="store.selectedIndex" @select="store.selectRow" />
      </div>
      <div class="app-right">
        <DetailPanel
          v-if="store.selectedMod"
          :mod="store.selectedMod"
          :selected-index="store.selectedIndex"
          :rows="store.rows"
        />
        <div v-else class="app-placeholder">Select a mod to see details.</div>
      </div>
    </div>

    <StatusBar :text="store.summary || 'No folder loaded.'" />

    <ConflictGraph
      v-if="store.showConflictGraph"
      :conflicts="store.conflicts"
      @close="store.showConflictGraph = false"
    />
    <ApplyDialog v-if="store.showApplyDialog" @close="store.showApplyDialog = false" />
    <ConflictResolutionDialog
      v-if="store.showConflictResolution"
      :mods="store.newConflictingMods"
      @close="onConflictResolutionClose"
    />
    <RestoreDialog v-if="store.showRestoreDialog" @close="store.showRestoreDialog = false" />
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useAppStore } from './stores/app'
import { useWails } from './composables/useWails'
import Toolbar from './components/Toolbar.vue'
import PathBar from './components/PathBar.vue'
import StatusBar from './components/StatusBar.vue'
import ModListTable from './components/ModListTable.vue'
import DetailPanel from './components/DetailPanel.vue'
import ConflictGraph from './components/ConflictGraph.vue'
import ApplyDialog from './components/ApplyDialog.vue'
import ConflictResolutionDialog from './components/ConflictResolutionDialog.vue'
import RestoreDialog from './components/RestoreDialog.vue'

const store = useAppStore()
const wails = useWails()

onMounted(async () => {
  wails.registerEvents()
  await wails.loadConfig()
  if (store.modStructure === 'MO2') {
    if (store.mo2Dir) {
      await wails.getMO2Profiles(store.mo2Dir)
    }
    if (store.mo2Dir && store.mo2Profile) {
      await wails.scanMO2(store.mo2Dir, store.mo2Profile)
    }
  } else if (store.modDir) {
    await wails.runScan(store.modDir)
  }
})

async function onOpenFolder() {
  if (store.modStructure === 'MO2') {
    await onPickMO2()
  } else {
    const dir = await wails.pickFolder()
    if (dir) await wails.runScan(dir)
  }
}

async function onRescan() {
  if (store.modStructure === 'MO2' && store.mo2Dir && store.mo2Profile) {
    await wails.scanMO2(store.mo2Dir, store.mo2Profile)
  } else {
    await wails.runScan()
  }
}

async function onScan() {
  await wails.runScan(store.modDir)
}

async function onPick() {
  await wails.pickFolder()
}

async function onPickMO2() {
  const dir = await wails.pickFolder()
  if (dir) {
    store.mo2Dir = dir
    await wails.getMO2Profiles(dir)
  }
}

async function onScanMO2() {
  await wails.scanMO2(store.mo2Dir, store.mo2Profile)
}

function onConflictResolutionClose() {
  store.showConflictResolution = false
  store.newConflictingMods = []
}
</script>

<style scoped>
.app-layout {
  display: flex;
  flex-direction: column;
  height: 100vh;
  overflow: hidden;
}
.app-top {
  flex-shrink: 0;
}
.app-split {
  flex: 1;
  display: flex;
  overflow: hidden;
}
.app-left {
  width: 60%;
  overflow: hidden;
  border-right: 1px solid var(--cp-border);
}
.app-right {
  flex: 1;
  overflow-y: auto;
}
.app-placeholder {
  padding: 24px;
  color: var(--cp-dim);
}
</style>
