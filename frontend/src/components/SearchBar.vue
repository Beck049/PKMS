<template>
  <div class="search-bar" :class="[`search-bar--${size}`]">
    <div class="search-input-wrapper">
      <input
        v-model="searchValue"
        type="text"
        :placeholder="placeholder"
        class="search-input"
        @input="handleInput"
        @keyup.enter="handleSearch"
      />
      <button
        class="search-button"
        @click="handleSearch"
        :disabled="!searchValue.trim()"
      >
        <svg
          class="search-icon"
          width="16"
          height="16"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <circle cx="11" cy="11" r="8"></circle>
          <path d="m21 21-4.35-4.35"></path>
        </svg>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

interface Props {
  size?: 'small' | 'large'
  placeholder?: string
  modelValue?: string
}

const props = withDefaults(defineProps<Props>(), {
  size: 'large',
  placeholder: '搜尋...',
  modelValue: ''
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
  'search': [value: string]
}>()

const searchValue = ref(props.modelValue)

const handleInput = () => {
  emit('update:modelValue', searchValue.value)
}

const handleSearch = () => {
  if (searchValue.value.trim()) {
    emit('search', searchValue.value.trim())
  }
}
</script>

<style scoped>
.search-bar {
  display: flex;
  align-items: center;
}

.search-input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
  background: transparent;
  border: 1px solid #52525b;
  border-radius: 8px;
  overflow: hidden;
  transition: all 0.2s ease;
}

.search-input-wrapper:focus-within {
  box-shadow: 0 0 0 2px rgba(98, 0, 255, 0.25);
}

.search-input {
  border: none;
  outline: none;
  background: transparent;
  color: #333;
  font-family: inherit;
}

.search-input::placeholder {
  color: #999;
}

.search-button {
  display: flex;
  align-items: center;
  background: #626262;
  justify-content: center;
  border: none;
  color: white;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.search-button:hover:not(:disabled) {
  background: #5f00b3;
}

.search-button:disabled {
  cursor: not-allowed;
}

.search-icon {
  display: block;
}

/* Small size */
.search-bar--small .search-input-wrapper {
  height: 32px;
}

.search-bar--small .search-input {
  padding: 0 12px;
  font-size: 14px;
  width: 200px;
}

.search-bar--small .search-button {
  width: 32px;
  height: 32px;
}

.search-bar--small .search-icon {
  width: 14px;
  height: 14px;
}

/* Large size */
.search-bar--large .search-input-wrapper {
  height: 48px;
}

.search-bar--large .search-input {
  padding: 0 16px;
  font-size: 16px;
  width: 300px;
}

.search-bar--large .search-button {
  width: 48px;
  height: 48px;
}

.search-bar--large .search-icon {
  width: 18px;
  height: 18px;
}
</style> 