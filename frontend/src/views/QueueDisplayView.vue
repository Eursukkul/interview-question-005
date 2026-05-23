<template>
  <QueueCard title="IT 05-2">
    <div class="display-view">
      <p class="queue-label">หมายเลขคิว</p>

      <LoadingSpinner v-if="loading" />
      <div v-else class="queue-number">{{ queueNumber }}</div>

      <p v-if="!loading && issuedAtText" class="queue-timestamp">{{ issuedAtText }}</p>

      <ErrorAlert :message="errorMessage" />
      <BaseButton variant="primary" @click="router.push('/')">กลับไปหน้ารับบัตรคิว</BaseButton>
    </div>
  </QueueCard>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getCurrentQueue } from '../api/queue'
import BaseButton from '../components/BaseButton.vue'
import ErrorAlert from '../components/ErrorAlert.vue'
import LoadingSpinner from '../components/LoadingSpinner.vue'
import QueueCard from '../components/QueueCard.vue'

const route = useRoute()
const router = useRouter()
const loading = ref(false)
const errorMessage = ref('')
const queueNumber = ref('00')
const issuedAtText = ref('')

onMounted(async () => {
  const queryQueue = route.query.queue
  const queryIssuedAt = route.query.issuedAt

  if (typeof queryQueue === 'string' && queryQueue.length > 0) {
    queueNumber.value = queryQueue
    if (typeof queryIssuedAt === 'string' && queryIssuedAt.length > 0) {
      issuedAtText.value = formatIssuedAt(queryIssuedAt)
    }
    return
  }

  loading.value = true
  try {
    const snapshot = await getCurrentQueue()
    queueNumber.value = snapshot.queueNumber
    issuedAtText.value = snapshot.issuedAt ? formatIssuedAt(snapshot.issuedAt) : ''
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'Unable to load current queue'
  } finally {
    loading.value = false
  }
})

function formatIssuedAt(iso: string): string {
  const date = new Date(iso)
  if (Number.isNaN(date.getTime())) {
    return ''
  }
  const day = String(date.getDate()).padStart(2, '0')
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const year = date.getFullYear()
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  return `วันที่ : ${day}/${month}/${year} เวลา ${hours}:${minutes} น.`
}
</script>
