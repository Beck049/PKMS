<template>
  <ul class="note-tree">
    <li v-for="node in nodes" :key="node.path" class="note-tree-item">
      <template v-if="node.isDir">
        <div class="folder-header">
          <button class="expand-btn" @click="toggle(node)">
            <span v-if="isExpanded(node)">▼</span>
            <span v-else>▶</span>
          </button>
          <span class="folder-name" @click="emitNodeClick(node.path)">{{ node.name }}</span>
        </div>
        <NoteTree v-if="isExpanded(node) && node.children && node.children.length" :nodes="node.children" @node-click="$emit('node-click', $event)" />
      </template>
      <template v-else>
        <span class="file-name">{{ node.name }}</span>
      </template>
    </li>
  </ul>
</template>

<script setup lang="ts">
import { ref, defineEmits } from 'vue'

interface NoteNode {
  name: string
  path: string
  isDir: boolean
  children?: NoteNode[]
}

defineProps<{ nodes: NoteNode[] }>()

const emit = defineEmits<{ (e: 'node-click', path: string): void }>()

// 用 path 當 key，記錄展開狀態
const expanded = ref<Record<string, boolean>>({})

function toggle(node: NoteNode) {
  expanded.value[node.path] = !expanded.value[node.path]
}
function isExpanded(node: NoteNode) {
  return !!expanded.value[node.path]
}
function emitNodeClick(path: string) {
  emit('node-click', path)
}
</script>

<style scoped>
.note-tree {
  list-style: none;
  margin: 0;
  padding-left: 20px;
}
.note-tree-item {
  padding: 4px 0;
}
.folder-header {
  display: flex;
  align-items: center;
}
.expand-btn {
  background: none;
  border: none;
  cursor: pointer;
  margin-right: 4px;
  font-size: 14px;
  padding: 0 2px;
  color: #666;
}
.folder-name {
  width: 100%;
  font-weight: bold;
  border-radius: 6px;
  color: #A1A1AA;
  padding-left: 4px;
  margin-right: 4px;
  cursor: pointer;
  transition: background 0.15s;
}
.folder-name:hover {
  background: #8C8C9380;
}
.file-name {
  display: block;
  width: 100% - 22px;
  border-radius: 6px;
  color: #A1A1AA;
  padding-left: 4px;
  margin-left: 14px;
  margin-right: 4px;
  cursor: pointer;
  transition: background 0.15s;
}
.file-name:hover {
  background: #8C8C9380;
}
</style> 