<template>
  <div class="content-view">
    <div v-if="loading" class="loading">
      <p>Loading content...</p>
    </div>
    
    <div v-else-if="error" class="error">
      <h2>Error</h2>
      <p>{{ error }}</p>
    </div>
    
    <div v-else-if="content" class="content">
      <div class="content-header">
        <div class="header-left">
          <a class="workspace-link" href="/">
            <span class="workspace-circle"></span>
            <span class="workspace-text">我的工作空間</span>
          </a>
        </div>
        <div class="header-center">
          <span class="pin" v-if="content.pin">📌</span>
          <h1>{{ content.title }}</h1>
          <button class="info-btn">i</button>
        </div>
        <div class="header-right">
          <span class="ref-count">References: {{ content.ref_count }}</span>
          <div class="modify-buttons">
            <button class="update-btn" @click="onUpdate">Update</button>
            <button class="delete-btn" @click="onDelete">Delete</button>
          </div>
        </div>
      </div>
      <div class="content-body">
        <template v-if="!leftClosed">
          <div
            class="left-pane"
            :style="{ width: leftWidth + 'px', minWidth: (innerWidth * 0.3) + 'px' }"
          >
            <ContentEditBlock v-model="content.rawdata" />
          </div>
        </template>
        <div class="divider" @mousedown="startDragging">
          <template v-if="leftClosed">
            <button class="restore-btn" @click.stop="restoreLeft">&gt;</button>
          </template>
          <template v-else-if="rightClosed">
            <button class="restore-btn" @click.stop="restoreRight">&lt;</button>
          </template>
        </div>
        <template v-if="!rightClosed">
          <div class="right-pane">
            <ContentRenderBlock v-model="content.rawdata" />
          </div>
        </template>
      </div>
    </div>
    
    <div v-else class="no-content">
      <p>No content found</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import ContentRenderBlock from '@/components/ContentRenderBlock.vue'
import ContentEditBlock from '@/components/ContentEditBlock.vue'
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { useRoute, useRouter } from 'vue-router'

interface ContentData {
  id: number
  title: string
  path: string
  ref_count: number
  pin: boolean
  rawdata: string
}

const route = useRoute()
const router = useRouter()
const content = ref<ContentData | null>(null)
const originalContent = ref<ContentData | null>(null)
const loading = ref(true)
const error = ref<string | null>(null)

const leftWidth = ref(window.innerWidth * 0.3)
const dragging = ref(false)
const leftClosed = ref(false)
const rightClosed = ref(false)

const innerWidth = ref(window.innerWidth)

const handleResize = () => {
  innerWidth.value = window.innerWidth
}

const fetchContent = async () => {
  try {
    loading.value = true
    error.value = null
    
    const id = route.params.id as string
    const response = await fetch(`/api/content/${id}`)
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }
    
    const data = await response.json()
    content.value = data
    originalContent.value = JSON.parse(JSON.stringify(data))
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'An error occurred while fetching content'
  } finally {
    loading.value = false
  }
}

const startDragging = (e: MouseEvent) => {
  if (leftClosed.value || rightClosed.value) return
  dragging.value = true
  document.body.style.cursor = 'col-resize'
}

const stopDragging = () => {
  dragging.value = false
  document.body.style.cursor = ''
}

const onDrag = (e: MouseEvent) => {
  if (!dragging.value) return
  let newWidth = e.clientX
  const minWidth = window.innerWidth * 0.3
  const maxWidth = window.innerWidth * 0.7
  if (newWidth < minWidth) {
    leftClosed.value = true
    leftWidth.value = 0
    dragging.value = false
    document.body.style.cursor = ''
    return
  }
  if (newWidth > maxWidth) {
    rightClosed.value = true
    leftWidth.value = window.innerWidth
    dragging.value = false
    document.body.style.cursor = ''
    return
  }
  leftWidth.value = newWidth
}

const restoreLeft = () => {
  leftClosed.value = false
  leftWidth.value = window.innerWidth * 0.35
}
const restoreRight = () => {
  rightClosed.value = false
  leftWidth.value = window.innerWidth * 0.65
}

function stripFrontmatter(src: string): string {
  if (src.startsWith('---')) {
    const end = src.indexOf('---', 3)
    if (end !== -1) {
      return src.slice(end + 3).replace(/^\s*\n/, '')
    }
  }
  return src
}

async function onUpdate() {
  if (!content.value || !originalContent.value) return
  const id = content.value.id
  // 只將有變動的欄位放入 body
  const body: Record<string, any> = {}
  if (content.value.title !== originalContent.value.title) {
    body.title = content.value.title
  }
  if (content.value.pin !== originalContent.value.pin) {
    body.pin = content.value.pin
  }
  if (content.value.rawdata !== originalContent.value.rawdata) {
    body.content = stripFrontmatter(content.value.rawdata)
  }
  if (Object.keys(body).length === 0) {
    console.log('內容未變更，不需更新')
    return
  }
  try {
    const res = await fetch(`/api/articles/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body)
    })
    if (!res.ok) throw new Error('Update failed')
    const result = await res.json()
    console.log('更新成功', result)
    await fetchContent()
  } catch (err) {
    console.error('更新失敗', err)
  }
}

function onDelete() {
  if (!content.value) return
  const id = content.value.id
  fetch(`/api/articles/${id}`, {
    method: 'DELETE'
  })
    .then(res => {
      if (!res.ok) throw new Error('Delete failed')
      return res.json()
    })
    .then(result => {
      console.log('刪除成功', result)
      // 可選：導回首頁
      router.push(`/`)
    })
    .catch(err => {
      console.error('刪除失敗', err)
    })
}

onMounted(() => {
  fetchContent()
  window.addEventListener('mousemove', onDrag)
  window.addEventListener('mouseup', stopDragging)
  window.addEventListener('resize', handleResize)
})
onBeforeUnmount(() => {
  window.removeEventListener('mousemove', onDrag)
  window.removeEventListener('mouseup', stopDragging)
  window.removeEventListener('resize', handleResize)
})
</script>

<style scoped>
.content-view {
  width: 100vw;
  height: 100%;
  padding: 0;
  margin: 0;
  box-sizing: border-box;
  background-color: #27272A;
}

.loading {
  text-align: center;
  padding: 40px;
  color: #666;
}

.error {
  color: #d32f2f;
  padding: 20px;
  border: 1px solid #d32f2f;
  border-radius: 4px;
  background-color: #ffebee;
}

.content-header {
  display: flex;
  align-items: center;
  height: 50px;
  font-size: 14px;
  padding: 0 16px;
  box-sizing: border-box;
  border-bottom: 1px solid #5D5D64;
}
.header-left {
  display: flex;
  align-items: center;
  flex: 0 0 auto;
  min-width: 0;
}
.header-center {
  display: flex;
  align-items: center;
  justify-content: center;
  flex: 1 1 0;
  min-width: 0;
  text-align: center;
  gap: 8px;
}
.header-center h1 {
  font-size: 16px;
  margin: 0;
  font-weight: bold;
  color: #D4D4D8;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.info-btn {
  background: transparent;
  border: solid 1px #5D5D64;
  border-radius: 50%;
  width: 14px;
  height: 14px;
  font-size: 10px;
  color: #5D5D64;
  cursor: pointer;
  margin-left: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.header-right {
  display: flex;
  align-items: center;
  flex: 0 0 auto;
  gap: 12px;
  min-width: 0;
  margin-left: auto;
}
.ref-count {
  background: #e3f2fd;
  padding: 2px 8px;
  border-radius: 4px;
  color: #1976d2;
  font-size: 14px;
}
.modify-buttons {
  display: flex;
  gap: 6px;
}

.content-body {
  display: flex;
  height: 100vh;
  min-height: 300px;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.03);
  overflow: hidden;
}

.left-pane {
  min-width: 100px;
  height: 100%;
  box-sizing: border-box;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
}

.divider {
  width: 6px;
  min-width: 6px;
  max-width: 6px;
  height: 100%;
  background: #5D5D64;
  cursor: col-resize;
  z-index: 2;
  transition: background 0.2s;
  align-self: stretch;
  display: flex;
  align-items: center;
  justify-content: center;
}

.restore-btn {
  background: #fff;
  border: 1px solid #bbb;
  border-radius: 50%;
  width: 28px;
  height: 28px;
  font-size: 18px;
  cursor: pointer;
  box-shadow: 0 2px 8px rgba(0,0,0,0.05);
  display: flex;
  align-items: center;
  justify-content: center;
}

.right-pane {
  flex: 1;
  height: 100%;
  box-sizing: border-box;
  padding: 16px;
  overflow-y: auto;
}

.rawdata-section {
  margin-bottom: 30px;
}

.rawdata-section h3 {
  color: #555;
  margin-bottom: 10px;
  border-bottom: 1px solid #eee;
  padding-bottom: 5px;
}

.rawdata-section pre {
  background-color: #f5f5f5;
  padding: 15px;
  border-radius: 4px;
  overflow-x: auto;
  font-size: 14px;
  line-height: 1.4;
  white-space: pre-wrap;
  word-wrap: break-word;
}

.no-content {
  text-align: center;
  padding: 40px;
  color: #666;
}

.update-btn {
  background: #e3f2fd;
  color: #1565c0;
  border: 1.5px solid #1565c0;
  border-radius: 8px;
  padding: 8px 20px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.15s, color 0.15s, border 0.15s;
}
.update-btn:hover {
  background: #bbdefb;
  color: #0d47a1;
  border-color: #0d47a1;
}

.delete-btn {
  background: #ffebee;
  color: #c62828;
  border: 1.5px solid #c62828;
  border-radius: 8px;
  padding: 8px 20px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.15s, color 0.15s, border 0.15s;
}
.delete-btn:hover {
  background: #ffcdd2;
  color: #b71c1c;
  border-color: #b71c1c;
}
.workspace-link {
  display: flex;
  align-items: center;
  text-decoration: none;
  margin-right: 16px;
  border-radius: 6px;
}
.workspace-link:hover {
  background-color: #8C8C9380;
}
.workspace-circle {
  width: 25px;
  height: 25px;
  border-radius: 50%;
  background: #9238f8;
  display: inline-block;
  margin-right: 8px;
}
.workspace-text {
  font-size: 16px;
  font-weight: bolder;
  color: #D4D4D8;
}
</style> 