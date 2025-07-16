<template>
  <div class="render-area" v-html="renderedHtml"></div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import MarkdownIt from 'markdown-it'
import container from 'markdown-it-container'
const props = defineProps<{ modelValue: string }>()

function stripFrontmatter(src: string): string {
  // 僅當第一行是 --- 時才處理
  if (src.startsWith('---')) {
    // 找到第二個 ---
    const end = src.indexOf('---', 3)
    if (end !== -1) {
      // 跳過 frontmatter，從下一行開始
      return src.slice(end + 3).replace(/^\s*\n/, '')
    }
  }
  return src
}

const md = new MarkdownIt()
// 支援多種類型的 container 插件
const containerTypes = ['info', 'success', 'warning', 'danger', 'spoiler']

containerTypes.forEach(type => {
  md.use(container, type, {
    validate: function(params: string) {
      return new RegExp(`^${type}( [A-Za-z0-9+-]{0,50})?$`).test(params.trim())
    },
    render: function (tokens: any, idx: number) {
      const m = tokens[idx].info.trim().match(new RegExp(`^${type}( [A-Za-z0-9+-]{0,50})?`))
      if (tokens[idx].nesting === 1) {
        let label = ''
        if (m && m[1]) {
          label = m[1].trim()
        }
        
        // 根據類型設定不同的 CSS class
        let alertClass = `alert alert__${type}`
        let labelClass = `alert__${type}-label`
        
        // spoiler 特殊處理
        if (type === 'spoiler') {
          return `<details class="alert alert__spoiler"><summary class="alert__spoiler-label">${label || 'Spoiler'}</summary>\n`
        }
        
        return `<div class="${alertClass}"><h5 class="${labelClass}">${label || type.charAt(0).toUpperCase() + type.slice(1)}</h5>\n`
      } else {
        // spoiler 特殊處理
        if (type === 'spoiler') {
          return '</details>\n'
        }
        return '</div>\n'
      }
    }
  })
})
const renderedHtml = computed(() => md.render(stripFrontmatter(props.modelValue || '')))
</script>

<style scoped>
.render-area {
  width: 100%;
  min-height: 300px;
  border-radius: 8px;
  padding: 16px;
  font-size: 15px;
  line-height: 0.6;
  color: #D4D4D8;
  box-sizing: border-box;
  white-space: pre-wrap;
  word-break: break-word;
}

/* ********** render style (must be global class) ********** */
/* header */
:global(.render-area :is(h1, h2, h3, h4, h5)) {
  line-height: 1.0;
  height: auto;
  padding-top: 6px;
  padding-bottom: 6px;
  margin: 0;
}

/* list */
:global(.render-area :is(ul, ol)) {
  margin: 2px 0 2px 24px;
  padding: 0;
}
:global(.render-area li) {
  margin: 2px 0;
  padding: 0;
}

/* code */
:global(.render-area pre) {
  margin: 10px 0;
  padding: 8px 12px;
  background-color: #303036;
}

/* table */
:global(.render-area table) {
  width: 100%;
  border-collapse: collapse;
  margin: 16px 0;
  background: #23232b;
  color: #e5e5e5;
}
:global(.render-area th), :global(.render-area td) {
  border: 1px solid #444;
  padding: 8px 12px;
  text-align: left;
}
:global(.render-area th) {
  background: #353545;
  font-weight: bold;
}
:global(.render-area tr:hover) {
  background: #2d2d3a;
}

/* blockquote */
:global(.render-area blockquote) {
  margin: 8px 0;
  padding-left: 12px;
  border-left: 3px solid #b3b3b3;
  color: #666;
  background: #f7f7f7;
}

/* ********** markdown-it-container ********** */
:global(.alert) {
  line-height: 0.6;
  border-radius: 8px;
  border-left: 2px solid;
  padding: 16px 18px 12px 18px;
  margin: 12px 0;
  border-radius: 8px;
  padding: 16px 18px 12px 18px;
  margin: 10px 0px 20px 0px;
}
:global(p) {
  line-height: 1.6;
}

/* Info 樣式 */
:global(.alert__info) {
  background-color: #38bdf81a;
  color: #38bdf8;
  border-left-color: #0ea5e9;
}
:global(.alert__info-label) {
  font-weight: bold;
  font-size: 1.1em;
  margin: 0 0 8px 0;
  color: #1a71f2;
}

/* Success 樣式 */
:global(.alert__success) {
  background-color: #6db19d26;
  color: #6db19d;
  border-left-color: #55b685;
}
:global(.alert__success-label) {
  font-weight: bold;
  font-size: 1.1em;
  margin: 0 0 8px 0;
  color: #15803d;
}

/* Warning 樣式 */
:global(.alert__warning) {
  background-color: #fbbf241a;
  color: #f59e0b;
  border-left-color: #f59e0b;
}
:global(.alert__warning-label) {
  font-weight: bold;
  font-size: 1.1em;
  margin: 0 0 8px 0;
  color: #a16207;
}

/* Danger 樣式 */
:global(.alert__danger) {
  background-color: #ef444433;
  color: #f87171;
  border-left-color: #ef4444;
}
:global(.alert__danger-label) {
  font-weight: bold;
  font-size: 1.1em;
  margin: 0 0 8px 0;
  color: #dc2626;
}

/* Spoiler 樣式 */
:global(.alert__spoiler) {
  background-color: #f3f4f633;
  color: #374151;
  border-left-color: #6b7280;
}
:global(.alert__spoiler-label) {
  font-weight: bold;
  font-size: 1.1em;
  margin: 0 0 8px 0;
  color: #2b2e32;
  cursor: pointer;
}
:global(.alert__spoiler-label:hover) {
  color: #374151;
}

</style>
