<template>
  <button class="create-note-button" @click="handleClick">
    <svg
      class="plus-icon"
      width="20"
      height="20"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      stroke-width="2"
      stroke-linecap="round"
      stroke-linejoin="round"
    >
      <line x1="12" y1="5" x2="12" y2="19"></line>
      <line x1="5" y1="12" x2="19" y2="12"></line>
    </svg>
    <span class="button-text">建立新筆記</span>
  </button>
</template>

<script setup lang="ts">
interface Props {
  disabled?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  disabled: false
})

const emit = defineEmits<{
  'click': [event: MouseEvent]
}>()

const handleClick = (event: MouseEvent) => {
  if (!props.disabled) {
    emit('click', event)
  }
}
</script>

<style scoped>
.create-note-button {
  display: flex;
  width: 100%;
  align-items: center;
  justify-content: center;
  gap: 8px;
  height: 48px;
  padding: 0 16px;
  background: #6064a9;
  border: none;
  border-radius: 8px;
  color: white;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  font-family: inherit;
}

.create-note-button:hover:not(:disabled) {
  background: #7a7fd8;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(146, 150, 235, 0.3);
}

.create-note-button:active:not(:disabled) {
  transform: translateY(0);
  box-shadow: 0 2px 8px rgba(146, 150, 235, 0.3);
}

.create-note-button:disabled {
  background: #ccc;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

.plus-icon {
  flex-shrink: 0;
}

.button-text {
  white-space: nowrap;
}
</style> 