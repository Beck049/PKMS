<template>
  <div class="popup-mask" @click.self="onClose">
    <div class="popup">
      <h2>建立新筆記</h2>
      <form @submit.prevent="onSubmit">
        <div class="form-group">
          <label for="title">Title</label>
          <input id="title" v-model="title" type="text" required autofocus />
        </div>
        <div class="form-group">
          <label for="path">Path</label>
          <input id="path" v-model="path" type="text" required />
        </div>
        <div class="form-group">
          <label for="type">Type</label>
          <select id="type" v-model="type">
            <option value="markdown">Markdown</option>
            <option value="index">Index</option>
            <option value="draw">Draw</option>
          </select>
        </div>
        <div class="form-group">
          <label for="tags">Tags</label>
          <input id="tags" v-model="tags" type="text" placeholder="以逗號分隔" />
        </div>
        <div class="form-actions">
          <button type="button" @click="onClose">取消</button>
          <button type="submit">確定</button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, watchEffect, computed } from 'vue'

const props = defineProps<{ modelValue: boolean, hierarchy?: string }>()
const emit = defineEmits(['close', 'submit'])

const title = ref('')
const path = ref('')
const type = ref('markdown')
const tags = ref('')

// 當 popup 開啟或 hierarchy 變動時，預設 path 為 hierarchy
watch(() => props.hierarchy, (val) => {
  path.value = val || ''
}, { immediate: true })

// 根據 title 自動產生 path（在 hierarchy 後面）
watch(title, (val) => {
  let slug = val
    .trim()
    .replace(/[^a-z0-9\u4e00-\u9fa5\s-]/gi, '')
    .replace(/\s+/g, '-')

    if(val != '') {
        if(type.value == 'markdown'){
            slug = '/' + slug + '.md'
        }
    } else {
        slug = ''
    }
  path.value = (props.hierarchy || '') + slug
})

function onClose() {
  emit('close')
  // 清空欄位
  title.value = ''
  path.value = ''
  type.value = 'markdown'
  tags.value = ''
}

function onSubmit() {
  emit('submit', {
    title: title.value,
    path: path.value,
    type: type.value,
    tags: tags.value
      .split(',')
      .map(t => t.trim())
      .filter(t => t.length > 0)
  })
  onClose()
}
</script>

<style scoped>
.popup-mask {
  position: fixed;
  z-index: 1000;
  left: 0; top: 0; right: 0; bottom: 0;
  background: rgba(0,0,0,0.25);
  display: flex;
  align-items: center;
  justify-content: center;
}
.popup {
  background: #fff;
  border-radius: 10px;
  box-shadow: 0 2px 16px rgba(0,0,0,0.15);
  padding: 32px 28px 24px 28px;
  min-width: 320px;
  max-width: 90vw;
}
.form-group {
  margin-bottom: 18px;
  display: flex;
  flex-direction: column;
}
.form-group label {
  font-weight: 500;
  margin-bottom: 6px;
}
.form-group input,
.form-group select {
  padding: 7px 10px;
  border: 1px solid #bbb;
  border-radius: 5px;
  font-size: 15px;
}
.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 10px;
}
.form-actions button {
  padding: 7px 18px;
  border-radius: 5px;
  border: none;
  font-size: 15px;
  cursor: pointer;
}
.form-actions button[type="button"] {
  background: #eee;
  color: #444;
}
.form-actions button[type="submit"] {
  background: #38bdf8;
  color: #fff;
}
</style>
