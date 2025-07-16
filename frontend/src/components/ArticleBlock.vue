<template>
    <div class="carousel-item" @click="goToContent(item.id)">
        <div class="item-header">
            <div class="title-pin">
                <span class="title">{{ item.title }}</span>
                <span v-if="item.pin" class="pin-icon">📌</span>
            </div>
            <button class="more-btn">⋯</button>
        </div>
        <div class="item-tags">
            <span v-for="tag in item.tags" :key="tag" class="tag">{{ tag }}</span>
        </div>
    </div>
</template>
  
<script setup lang="ts">
    import { defineProps } from 'vue'
    import { useRouter } from 'vue-router'
    
    defineProps<{ item: any }>()
    const router = useRouter()

    function goToContent(id: number) {
    router.push(`/content/${id}`)
    }
</script>
<style scoped>
.carousel-item {
    width: 224px;
    min-height: 120px;
    box-sizing: border-box;
    padding: 12px 12px 8px 12px;
    border: 1px solid #71717A;
    margin: 0 2px;
    background: #303036;
    text-align: left;
    border-radius: 12px;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    box-shadow: 0 2px 8px rgba(0,0,0,0.03);
    cursor: pointer;
}
.item-header {
    display: flex;
    color: #D4D4D8;
    justify-content: space-between;
    align-items: flex-start;
}
.title-pin {
    display: flex;
    align-items: center;
    gap: 6px;
    font-weight: bold;
    font-size: 14px;
}
.pin-icon {
    color: #f7b801;
    font-size: 1.1em;
}
.more-btn {
    background: none;
    border: none;
    font-size: 1.3em;
    color: #888;
    cursor: pointer;
    padding: 0 4px;
    border-radius: 4px;
    transition: background 0.15s;
}
.more-btn:hover {
    background: #eee;
}
.item-tags {
    margin-top: 18px;
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
}
.tag {
    border: 1px solid #71717A;
    color: #7a7fd8;
    border-radius: 8px;
    padding: 2px 10px;
    font-size: 0.95em;
}
</style>