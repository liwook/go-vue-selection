<template>
  <div class="logo" v-if="settings.logoShow">
    <img :src="settings.logo" alt="logo" />
    <p>{{ settings.title }}</p>
  </div>
</template>

<script setup lang="ts">
import settings from '@/settings'
</script>

<style scoped lang="scss">
.logo {
  width: 100%;
  height: $base-logo-height;
  display: flex;
  align-items: center;
  padding: 0 10px;
  color: #fff;

  img {
    width: 40px;
    height: 40px;
    margin-right: 10px;
  }

  p {
    font-size: 18px;
    white-space: nowrap;
  }
}
</style>
