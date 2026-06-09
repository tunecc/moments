<template>
  <meting-js v-if="id && server && type && api" :server="server" :type="type" :id="id" :api="api"/>
</template>

<script setup lang="ts">
import type {MetingJSDTO} from "@/types";

const props = defineProps<MetingJSDTO>()

// 仅在音乐组件实际渲染时按需加载 APlayer/Meting 资源,
// 避免在没有音乐的页面上无条件加载这些脚本。
useHead({
  link: [
    {key: "aplayer-css", href: "/css/APlayer.min.css", rel: "stylesheet"},
  ],
  script: [
    {key: "aplayer-js", src: "/js/APlayer.min.js", type: "text/javascript", defer: true},
    {key: "meting-js", src: "/js/Meting.min.js", type: "text/javascript", defer: true},
  ],
})
</script>

<style scoped>

</style>
