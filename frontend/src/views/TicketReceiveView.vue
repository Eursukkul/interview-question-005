<template>
  <QueueCard title="IT 05-1">
    <div class="receive-view">
      <p class="subtitle">ระบบรับบัตรคิว</p>
      <div class="ticket-visual" aria-hidden="true">
        <span class="ticket-visual__stub"></span>
        <span class="ticket-visual__code">{{ currentQueue }}</span>
      </div>

      <ErrorAlert :message="errorMessage" />

      <div class="button-stack">
        <BaseButton :loading="loading" @click="receiveQueue">รับบัตรคิว</BaseButton>
        <BaseButton variant="secondary" @click="router.push('/reset')">ล้างคิว</BaseButton>
      </div>

      <p class="helper-text">คิวปัจจุบัน: {{ currentQueue }} · Queue range: A0 - Z9</p>
    </div>
  </QueueCard>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { getCurrentQueue, getNextQueue } from '../api/queue'
import BaseButton from '../components/BaseButton.vue'
import ErrorAlert from '../components/ErrorAlert.vue'
import QueueCard from '../components/QueueCard.vue'

const router = useRouter()
const loading = ref(false)
const errorMessage = ref('')
const currentQueue = ref('00')

onMounted(async () => {
  try {
    const snapshot = await getCurrentQueue()
    currentQueue.value = snapshot.queueNumber
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'Unable to load current queue'
  }
})

async function receiveQueue() {
  loading.value = true
  errorMessage.value = ''
  try {
    const snapshot = await getNextQueue()
    const query: Record<string, string> = { queue: snapshot.queueNumber }
    if (snapshot.issuedAt) {
      query.issuedAt = snapshot.issuedAt
    }
    await router.push({ path: '/queue', query })
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'Unable to receive queue'
  } finally {
    loading.value = false
  }
}
</script>
