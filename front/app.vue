<template>
  <NuxtLayout>
    <NuxtPage/>
    <Toaster richColors position="top-center"/>
  </NuxtLayout>
</template>

<script lang="ts" setup>
import { Toaster } from 'vue-sonner'
import "@fancyapps/ui/dist/fancybox/fancybox.css";
import '~/assets/simple-markdown.scss'
</script>

<style>
html {
  color-scheme: light;
  background: #f8fafc;
  color: #171717;
}

html.dark {
  color-scheme: dark;
  background: #0a0a0a;
  color: #f5f5f5;
}

body {
  min-height: 100vh;
  margin: 0;
  background: #f8fafc;
  color: #171717;
}

html.dark body {
  background: #0a0a0a;
  color: #f5f5f5;
}

div.grecaptcha-badge{
  display: none !important;
}
</style>
