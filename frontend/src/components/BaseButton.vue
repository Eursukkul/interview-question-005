<template>
  <button class="base-button" :class="`base-button--${variant}`" :disabled="disabled || loading">
    <span v-if="loading" class="base-button__spinner" aria-hidden="true"></span>
    <span><slot /></span>
  </button>
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    variant?: 'primary' | 'secondary' | 'danger'
    loading?: boolean
    disabled?: boolean
  }>(),
  {
    variant: 'primary',
    loading: false,
    disabled: false
  }
)
</script>
