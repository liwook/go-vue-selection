<template>
  <svg :style="{ width, height }">
    <use :href="symbolId" :fill="iconColor" />
  </svg>
</template>

<script setup lang="ts">
import { computed } from 'vue'

defineOptions({ name: 'SvgIcon' })

interface Props {
  name: string
  iconColor?: string
  width?: string
  height?: string
}

const props = withDefaults(defineProps<Props>(), {
  iconColor: '#333',
  width: '16px',
  height: '16px',
})

const symbolId = computed(() => `#icon-${props.name}`)
</script>
