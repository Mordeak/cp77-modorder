import { EventsOn } from '../../wailsjs/runtime/runtime'
import * as App from '../../wailsjs/go/main/App'
import { useAppStore } from '../stores/app'

export function useWails() {
  const store = useAppStore()

  function registerEvents() {
    EventsOn('scan:progress', (msg: string) => {
      store.scanProgress = msg
    })
  }

  async function loadConfig() {
    const cfg = await App.GetConfig()
    store.modDir = cfg.modDir
  }

  async function pickFolder(): Promise<string> {
    const dir = await App.PickFolder()
    if (dir) store.modDir = dir
    return dir
  }

  async function runScan(dir?: string) {
    const target = dir ?? store.modDir
    if (!target) return
    store.scanning = true
    store.scanError = ''
    try {
      const result = await App.Scan(target)
      store.setScanResult(result)
      if (result.hasModlist) {
        const newConflicting = result.rows
          .filter(r => r.unlisted && r.mod && r.mod.conflictCount > 0)
          .map(r => r.name)
        if (newConflicting.length > 0) {
          store.newConflictingMods = newConflicting
          store.showConflictResolution = true
        }
      }
    } catch (e: any) {
      store.scanError = String(e)
    } finally {
      store.scanning = false
    }
  }

  async function setPriority(name: string, p: number) {
    try {
      const result = await App.SetPriority(name, p)
      store.updateScanResult(result)
    } catch (e: any) {
      store.scanError = String(e)
    }
  }

  async function reorderGroup(names: string[]) {
    try {
      const result = await App.ReorderConflictGroup(names)
      store.updateScanResult(result)
    } catch (e: any) {
      store.scanError = String(e)
    }
  }

  async function writeModlist() {
    await App.WriteModlist()
  }

  async function getApplyPreview() {
    return App.GetApplyPreview()
  }

  async function getConflictGroup(name: string) {
    return App.GetConflictGroup(name)
  }

  return {
    registerEvents, loadConfig, pickFolder, runScan,
    setPriority, reorderGroup, writeModlist, getApplyPreview, getConflictGroup,
  }
}
