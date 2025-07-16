<template>
  <div>
    <div v-for="([key, arr]) in filteredDisplay" :key="key" class="content-group">
      <h2 class="content-title">{{ key }}</h2>
      <div>
        <Carousel :arr="arr" :articles="articles" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineProps } from 'vue'
import Carousel from './Carousel.vue'

defineProps<{ 
  filteredDisplay: [string, number[]][], 
  articles: any[]
}>()
</script>

<style scoped>
.content-group {
  padding-bottom: 32px;
}
.content-title {
  font-size: 1.3rem;
  font-weight: bold;
  color: #D4D4D8;
  border-bottom: 2px solid #5D5D64;
  padding-left: 16px;
  margin-bottom: 16px;
  letter-spacing: 1px;
}
</style>