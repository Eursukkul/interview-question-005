<template>
  <main class="page-shell">
    <section class="queue-card">
      <AppHeader :title="title" />
      <div class="queue-card__body">
        <slot />
      </div>
    </section>
  </main>
</template>

<script setup lang="ts">
import AppHeader from './AppHeader.vue'

defineProps<{
  title: string
}>()
</script>
