<template>
  <router-view v-slot="{ Component, route }">
    <transition name="fade" mode="out-in">
      <component :is="Component" :key="route.path" />
    </transition>
  </router-view>
</template>

<script setup lang="ts"></script>

<style>
/* 切页过渡动画：旧页淡出 + 轻微缩放，新页淡入 + 轻微缩放
   注意：必须【非 scoped】！scoped 会加上 [data-v-xxx] 属性选择器，
   而 transition 把 fade-enter-from 加到子组件根元素上，子组件根没有
   data-v 属性 → 选择器匹配不到 → 动画失效。
   另：必须用 CSS transition（不是 animation/@keyframes），否则异步组件
   + addRoute 二次导航的竞态会残留 opacity:0 导致白屏。 */
.fade-enter-active,
.fade-leave-active {
  transition:
    opacity 0.3s ease,
    transform 0.3s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: scale(0.98);
}
</style>
