<template>
  <div class="pathbar">
    <!-- MO2 mode -->
    <template v-if="modStructure === 'MO2'">
      <input
        type="text"
        :value="mo2Dir"
        @input="$emit('update:mo2Dir', ($event.target as HTMLInputElement).value)"
        @blur="mo2Dir && $emit('loadProfiles', mo2Dir)"
        placeholder="MO2 instance folder…"
      />
      <button @click="$emit('pickMO2')" title="Browse MO2 instance…">…</button>
      <select
        class="profile-select"
        :value="mo2Profile"
        @change="$emit('update:mo2Profile', ($event.target as HTMLSelectElement).value)"
        :disabled="!mo2Profiles.length"
      >
        <option value="" disabled>{{ mo2Profiles.length ? 'Select profile…' : 'No profiles found' }}</option>
        <option v-for="p in mo2Profiles" :key="p" :value="p">{{ p }}</option>
      </select>
      <button class="primary" @click="$emit('scanMO2')" :disabled="scanning || !mo2Dir || !mo2Profile">
        {{ scanning ? 'Scanning…' : 'Scan' }}
      </button>
    </template>

    <!-- Default mode -->
    <template v-else>
      <input
        type="text"
        :value="modelValue"
        @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
        @keydown.enter="$emit('scan')"
        placeholder="Paste path to mods folder, e.g. C:/Games/Cyberpunk 2077/archive/pc/mod"
      />
      <button @click="$emit('pick')" title="Browse…">…</button>
      <button class="primary" @click="$emit('scan')" :disabled="scanning">
        {{ scanning ? 'Scanning…' : 'Scan' }}
      </button>
    </template>

    <span v-if="error" class="msg text-error">{{ error }}</span>
    <span v-else-if="progress" class="msg text-dim">{{ progress }}</span>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  modStructure: 'default' | 'MO2'
  modelValue: string
  scanning: boolean
  progress: string
  error: string
  mo2Dir: string
  mo2Profile: string
  mo2Profiles: string[]
}>()
defineEmits<{
  'update:modelValue': [v: string]
  'scan': []
  'pick': []
  'update:mo2Dir': [v: string]
  'update:mo2Profile': [v: string]
  'pickMO2': []
  'scanMO2': []
  'loadProfiles': [instanceDir: string]
}>()
</script>

<style scoped>
.pathbar {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 5px 8px;
  background: #0D0D0D;
  border-bottom: 1px solid var(--cp-border);
}
.pathbar input { flex: 1; }
.profile-select {
  background: var(--cp-input);
  color: var(--cp-fg);
  border: 1px solid var(--cp-border);
  padding: 3px 6px;
  font-size: 12px;
  min-width: 160px;
  outline: none;
}
.profile-select:focus { border-color: var(--cp-focus); }
.profile-select:disabled { color: var(--cp-dim); }
.msg {
  font-size: 10px;
  white-space: nowrap;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  font-weight: 600;
}
</style>
