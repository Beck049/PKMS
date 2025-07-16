<template>
    <div class="tag-dropdown">
      <div class="dropdown-header" @click="toggleDropdown">
        <span>標籤列表</span>
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
        <div class="searchbar-center">
          <SearchBar
            v-model="search"
            size="small"
            placeholder="搜尋標籤"
          />
        </div>
        <ul class="tag-list">
          <li
            v-for="tag in filteredTags"
            :key="tag.id"
            class="tag-item"
            @click="onTagClick(tag.name)"
          >
            {{ tag.name }}
          </li>
        </ul>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { ref, computed, onMounted } from 'vue'
  import SearchBar from './SearchBar.vue'
  
  interface Tag {
    id: number
    name: string
  }
  
  const isOpen = ref(false)
  const search = ref('')
  const tags = ref<Tag[]>([])
  const loading = ref(false)
  
  const emit = defineEmits<{
    (e: 'tag-click', tag: string): void
  }>()
  
  const toggleDropdown = () => {
    isOpen.value = !isOpen.value
  }
  
  const fetchTags = async () => {
    loading.value = true
    try {
      const res = await fetch('http://localhost:8080/api/tags')
      if (res.ok) {
        tags.value = await res.json()
      }
    } catch (e) {
      // 可加錯誤處理
    }
    loading.value = false
  }
  
  const filteredTags = computed(() =>
    tags.value.filter(tag =>
      tag.name.toLowerCase().includes(search.value.toLowerCase())
    )
  )
  
  function onTagClick(tag: string) {
    emit('tag-click', tag)
  }
  
  onMounted(fetchTags)
  </script>
  
  <style scoped>
  .tag-dropdown {
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
  .tag-list {
    list-style: none;
    margin: 8px 0 0 0;
    padding: 0 16px;
    max-height: 200px;
    overflow-y: auto;
  }
  .tag-item {
    color: #A1A1AA;
    padding: 6px 6px;
    font-size: 14px;
    border-radius: 6px;
    transition: background 0.15s;
  }
  .tag-item:last-child {
    border-bottom: none;
  }
  .tag-item:hover {
    background: #8C8C9380;
    cursor: pointer;
  }
  .searchbar-center {
    display: flex;
    justify-content: center;
    align-items: center;
    margin-bottom: 8px;
  }
  </style>