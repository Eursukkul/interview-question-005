<template>
  <QueueCard title="IT 05-3">
    <div class="reset-view">
      <BaseButton variant="primary" :loading="loading" @click="handleReset">ล้างคิว</BaseButton>

      <p class="queue-label">หมายเลขคิวปัจจุบัน</p>
      <div class="queue-number">{{ currentQueue }}</div>

      <ErrorAlert :message="errorMessage" />
      <div v-if="successMessage" class="success-alert" role="status">{{ successMessage }}</div>

      <BaseButton variant="primary" @click="router.push('/')">กลับไปหน้ารับบัตรคิว</BaseButton>
    </div>
  </QueueCard>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { getCurrentQueue, resetQueue } from '../api/queue'
import BaseButton from '../components/BaseButton.vue'
import ErrorAlert from '../components/ErrorAlert.vue'
import QueueCard from '../components/QueueCard.vue'

const router = useRouter()
const currentQueue = ref('00')
const loading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

onMounted(async () => {
  try {
    const snapshot = await getCurrentQueue()
    currentQueue.value = snapshot.queueNumber
  } catch {
    currentQueue.value = '00'
  }
})

async function handleReset() {
  loading.value = true
  errorMessage.value = ''
  successMessage.value = ''
  try {
    const snapshot = await resetQueue()
    currentQueue.value = snapshot.queueNumber
    successMessage.value = 'ล้างคิวเรียบร้อยแล้ว'
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'Unable to reset queue'
  } finally {
    loading.value = false
  }
}
</script>
