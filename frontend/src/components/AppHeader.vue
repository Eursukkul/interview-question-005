<template>
  <header class="app-header">
    <h1>{{ title }}</h1>
  </header>
</template>

<script setup lang="ts">
defineProps<{
  title: string
}>()
</script>
