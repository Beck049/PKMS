<template>
    <textarea
      id="edit-area"
      v-model="modelValueProxy"
      class="edit-area"
      rows="20"
      spellcheck="false"
    ></textarea>
</template>

<script setup lang="ts">
import { computed } from 'vue'
const props = defineProps<{ modelValue: string }>()
const emit = defineEmits(['update:modelValue'])

const modelValueProxy = computed({
  get: () => props.modelValue,
  set: v => emit('update:modelValue', v)
})
</script>

<style scoped>
.edit-area {
  width: 100%;
  height: 100%;
  min-height: 300px;
  font-family: 'Fira Mono', 'Consolas', 'Menlo', monospace;
  font-size: 15px;
  padding: 12px;
  resize: vertical;
  background: transparent;
  color: #D4D4D8;
  box-sizing: border-box;
}
.edit-area:focus {
  outline: none;
  border: none;
  box-shadow: none;
}
</style>
