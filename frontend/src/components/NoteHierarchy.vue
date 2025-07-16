<template>
  <div class="note-hierarchy-dropdown">
    <div class="dropdown-header" @click="toggleDropdown">
      <span>我的筆記</span>
      <svg
        class="arrow"
        :class="{ open: isOpen }"
        width="20"
        height="20"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <polyline points="9 18 15 12 9 6" />
      </svg>
    </div>
    <div v-if="isOpen" class="dropdown-menu">
      <div v-if="loading" class="loading">載入中...</div>
      <template v-else>
        <NoteTree :nodes="hierarchy" @node-click="onNodeClick"/>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import NoteTree from './NoteTree.vue'

interface NoteNode {
  id: number
  path: string
  title: string
  children?: NoteNode[]
}

const isOpen = ref(false)
const loading = ref(false)
const hierarchy = ref<NoteNode[]>([])

const emit = defineEmits<{ (e: 'node-click', path: string): void }>()

function onNodeClick(path: string) {
  emit('node-click', path)
}

const toggleDropdown = () => {
  isOpen.value = !isOpen.value
  if (isOpen.value && hierarchy.value.length === 0) {
    fetchHierarchy()
  }
}

const fetchHierarchy = async () => {
  loading.value = true
  try {
    const res = await fetch('http://localhost:8080/api/hierarchy')
    if (res.ok) {
      hierarchy.value = await res.json()
    }
  } catch (e) {
    // 可加錯誤處理
  }
  loading.value = false
}
</script>

<style scoped>
.note-hierarchy-dropdown {
  width: 100%;
  font-size: 16px;
}
.dropdown-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: transparent;
  color: #A1A1AA;
  padding: 12px 16px;
  border-radius: 8px;
  cursor: pointer;
  user-select: none;
  transition: background 0.2s;
}
.dropdown-header:hover {
  background: #8C8C9380;
}
.arrow {
  transition: transform 0.2s;
}
.arrow.open {
  transform: rotate(90deg);
}
.dropdown-menu {
  margin-top: 8px;
  background: transparent;
  border: solid 1px #52525B;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
  padding: 12px 0;
}
.loading {
  text-align: center;
  color: #888;
  padding: 16px 0;
}
</style> 