<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, watchEffect, computed } from 'vue'
import SearchBar from '@/components/SearchBar.vue'
import CreateNoteButton from '@/components/CreateNoteButton.vue'
import TagDropdown from '@/components/TagDropdown.vue'
import NoteHierarchy from '@/components/NoteHierarchy.vue'
import ContentSection from '@/components/ContentSection.vue'
import CreatePopup from '@/components/CreatePopup.vue'
import { useRouter } from 'vue-router'
const router = useRouter()

const leftWidth = ref(300)
const dragging = ref(false)
const maxLeftWidth = ref(0)
const hierarchy = ref<string[]>(['My Notes']) // 'My Notes' 取代 '/app/articles'
// on change then search
const searchText = ref('')
const tagList = ref<string[]>([])
const path = ref<string>('')
const articles = ref<any[]>([])
const allTags = ref<string[]>([])
const display = ref<Record<string, number[]>>({})
const showCreatePopup = ref(false)

const filteredDisplay = computed(() => {
  return Object.entries(display.value).filter(([k, v]) => 
    (tagList.value.length === 0 || k === 'pin' || tagList.value.includes(k)) && v.length > 0
  )
})

const updateMaxLeftWidth = () => {
  maxLeftWidth.value = window.innerWidth * 0.3
  if (leftWidth.value > maxLeftWidth.value) {
    leftWidth.value = maxLeftWidth.value
  }
}

const startDragging = (e: MouseEvent) => {
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
  if (newWidth > maxLeftWidth.value) newWidth = maxLeftWidth.value
  if (newWidth < 100) newWidth = 100
  leftWidth.value = newWidth
}

const handleSearch = (value: string) => {
  console.log('搜尋:', value)
}

const handleCreateNote = () => {
  showCreatePopup.value = true
}

async function onCreateSubmit(data: any) {
  // 這裡可以處理送出後的邏輯
  try {
    const res = await fetch('/api/articles', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data)
    })
    if (!res.ok) throw new Error('API error')
    const result = await res.json()
    // 成功時，goto http://localhost:3000/content/id
    router.push(`/content/${result.article_id}`)
    console.log('建立成功', result)
  } catch (err) {
    console.error('建立失敗', err)
  }
  showCreatePopup.value = false
}

function handleTagClick(tag: string) {
  if (!tagList.value.includes(tag)) {
    tagList.value.push(tag)
  }
}

function removeTag(tag: string) {
  tagList.value = tagList.value.filter(t => t !== tag)
}

// 點擊 BreadCrumbs pop
function handleBreadcrumbClick(idx: number) {
  // 只保留 array 前 idx+1 項 其他捨去
  hierarchy.value = hierarchy.value.slice(0, idx + 1)
  path.value = constructPath()
}

// 點擊 NoteTree 時，更新 BreadCrumbs
function handleNodeClick(selectPath: string) {
  // 去除 '/app/articles' 前綴
  let cleanPath = selectPath.replace(/^\/app\/articles\/?/, '')
  const parts = cleanPath.split('/').filter(Boolean)
  hierarchy.value = ['My Notes', ...parts]
  path.value = constructPath()
}

// 用 Hierarchy 組成 Path
function constructPath(){
  // hierarchy.value[0] 是 'My Notes'，不參與 path 組合
  return hierarchy.value.length > 1 ? '/' + hierarchy.value.slice(1).join('/') : ''
}

async function fetchArticles() {
  const tagParam = tagList.value.join(',')
  const queryParam = searchText.value
  const pathParam = path.value.slice(1)
  const url = `http://localhost:8080/api/search?tag=${tagParam}&path=${pathParam}&query=${queryParam}`
  const res = await fetch(url)
  const data = await res.json()
  articles.value = data.articles || []
  allTags.value = data.allTags || []
}

watchEffect(() => {
  fetchArticles()
})

watchEffect(() => {
  // 構建 display 物件
  const d: Record<string, number[]> = { pin: [] }
  allTags.value.forEach(tag => { d[tag] = [] })
  articles.value.forEach((item, i) => {
    if (item.pin) d.pin.push(i)
    if (Array.isArray(item.tags)) {
      item.tags.forEach((tag: string) => {
        if (d[tag]) d[tag].push(i)
      })
    }
  })
  display.value = d
})

onMounted(() => {
  updateMaxLeftWidth()
  window.addEventListener('resize', updateMaxLeftWidth)
  window.addEventListener('mousemove', onDrag)
  window.addEventListener('mouseup', stopDragging)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', updateMaxLeftWidth)
  window.removeEventListener('mousemove', onDrag)
  window.removeEventListener('mouseup', stopDragging)
})
</script>

<template>
  <main class="container">
    <div
      class="left-pane"
      :style="{ width: leftWidth + 'px', maxWidth: maxLeftWidth + 'px' }"
    >
      <ul class="left-content">
        <li class="search-item">
          <SearchBar 
            v-model="searchText"
            @search="handleSearch"
            placeholder="搜尋筆記..."
          />
        </li>
        <li class="create-note-item">
          <CreateNoteButton @click="handleCreateNote" />
        </li>
        <li class="note-hierarchy">
          <NoteHierarchy @node-click="handleNodeClick" />
        </li>
        <li class="tag-drop-down">
          <TagDropdown @tag-click="handleTagClick" />
        </li>
      </ul>
    </div>
    <div
      class="divider"
      @mousedown="startDragging"
    ></div>
    <div class="right-pane">
      <div class="breadcrumbs">
        <template v-for="(item, idx) in hierarchy" :key="item">
          <div class="breadcrumb-item" @click="handleBreadcrumbClick(idx)">{{ item }}</div>
          <span v-if="idx !== hierarchy.length - 1" class="breadcrumb-sep">&gt;</span>
        </template>
      </div>
      <div class="tags">
        <div v-for="tag in tagList" :key="tag" class="tag-chip">
          <span>{{ tag }}</span>
          <button class="tag-remove" @click="removeTag(tag)">x</button>
        </div>
      </div>
      <div>
        <ContentSection :filteredDisplay="filteredDisplay" :articles="articles" />
      </div>
    </div>
    <CreatePopup v-if="showCreatePopup"
      v-model:modelValue="showCreatePopup"
      :hierarchy="Array.isArray(hierarchy) ? hierarchy.slice(1).join('/') : ''"
      @close="showCreatePopup = false"
      @submit="onCreateSubmit" />
  </main>
</template>

<style scoped>
.container {
  display: flex;
  height: 100vh;
  width: 100vw;
  overflow: hidden;
  background-color: #27272A;
}
.left-pane {
  min-width: 100px;
  max-width: 30vw;
  height: 100%;
  transition: background 0.2s;
  box-sizing: border-box;
  padding: 16px;
  /* */
  overflow-y: auto;
  display: flex;
  flex-direction: column;
}
.divider {
  width: 6px;
  cursor: col-resize;
  background: #5D5D64;
  height: 100%;
  z-index: 2;
}
.right-pane {
  flex: 1;
  height: 100%;
  box-sizing: border-box;
  padding: 16px;
  overflow-y: auto;
}

.left-content {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.search-item {
  width: 100%;
}

.create-note-item {
  width: 100%;
}

.tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin: 12px 0px;
}
.tag-chip {
  display: flex;
  align-items: center;
  border-radius: 12px;
  padding: 4px 10px 4px 12px;
  font-size: 14px;
  border: 1px solid #71717A;
  color: #7a7fd8;
}
.tag-remove {
  background: none;
  border: none;
  color: #888;
  margin-left: 6px;
  font-size: 14px;
  cursor: pointer;
  border-radius: 50%;
  width: 18px;
  height: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.15s;
}
.tag-remove:hover {
  background: #7a7fd8;
  color: #222;
}
.breadcrumbs {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-bottom: 12px;
  font-size: 24px;
}
.breadcrumb-item {
  font-weight: bolder;
  color: #A1A1AA;
  border-radius: 8px;
  padding: 4px 16px;
  cursor: pointer;
  transition: background 0.15s;
}
.breadcrumb-item:hover {
  color: #D4D4D8;
}
.breadcrumb-sep {
  color: #888;
  margin: 0 2px;
}
</style>
