<template>
  <div v-if="message" class="error-alert" role="alert">
    {{ message }}
  </div>
</template>

<script setup lang="ts">
defineProps<{
  message: string
}>()
</script>
