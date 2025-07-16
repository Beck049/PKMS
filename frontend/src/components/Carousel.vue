<template>
  <div class="carousel-container">
    <div class="arrow" @click="prev" :class="{ disabled: isFirst }">&lt;</div>
    <div class="carousel-viewport">
      <div class="carousel-track">
        <ArticleBlock
          v-for="(item, i) in visibleItems"
          :key="i"
          :item="item"
          class="carousel-article"
        />
      </div>
    </div>
    <div class="arrow" @click="next" :class="{ disabled: isLast }">&gt;</div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, nextTick, watch } from 'vue'
import ArticleBlock from './ArticleBlock.vue';

const props = defineProps<{ 
  arr: number[], 
  articles: any[]
}>()

const pageSize = 8
const currentIndex = ref(0)

const items = computed(() => 
  props.arr.map(i => props.articles[i])
)

const totalPages = computed(() =>
  Math.ceil(items.value.length / pageSize)
)

const isFirst = computed(() => currentIndex.value === 0)
const isLast = computed(() => currentIndex.value >= totalPages.value - 1)

const visibleItems = computed(() => {
  const start = currentIndex.value * pageSize
  return items.value.slice(start, start + pageSize)
})

const firstRow = computed(() => visibleItems.value.slice(0, 4))
const secondRow = computed(() => visibleItems.value.slice(4, 8))

const firstRowBlocks = ref([])
const firstRowWidth = ref(0)

function updateFirstRowWidth() {
  nextTick(() => {
    // 取第一排第一個 ArticleBlock 的 offsetWidth
    const el = (firstRowBlocks.value[0] as any)?.$el
    if (el) {
      firstRowWidth.value = el.offsetWidth
    }
  })
}

watch([firstRow, currentIndex], updateFirstRowWidth, { immediate: true })
onMounted(updateFirstRowWidth)

function prev() {
  if (!isFirst.value) currentIndex.value--
}
function next() {
  if (!isLast.value) currentIndex.value++
}
</script>

<style scoped>
.carousel-container {
  display: flex;
  align-items: center;
  width: 100%;
}
.arrow {
  width: 40px;
  text-align: center;
  font-size: 2rem;
  cursor: pointer;
  user-select: none;
  transition: color 0.2s;
}
.arrow.disabled {
  color: #ccc;
  pointer-events: none;
}
.carousel-viewport {
  overflow: hidden;
  flex: 1;
}
.carousel-track {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  width: 100%;
  transition: transform 0.5s cubic-bezier(.4,1.3,.6,1);
}
.carousel-article {
  flex: 0 1 calc((100% - 52px) / 4); /* 3 gaps * 12px + 8 margins * 2px = 52px */
  min-width: 0px;
  box-sizing: border-box;
}
</style>