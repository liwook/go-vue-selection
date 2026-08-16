<template>
  <div class="liquid-ball" :style="{ width: size + 'px', height: size + 'px' }">
    <svg viewBox="0 0 100 100" class="liquid-ball__svg">
      <defs>
        <!-- 圆形裁剪，把波浪裁成球内区域 -->
        <clipPath id="liquid-clip">
          <circle cx="50" cy="50" r="48" />
        </clipPath>
        <linearGradient id="liquid-fill" x1="0" y1="0" x2="0" y2="1">
          <stop offset="0%" :stop-color="color" stop-opacity="0.9" />
          <stop offset="100%" :stop-color="color" stop-opacity="0.5" />
        </linearGradient>
      </defs>

      <!-- 球底（未填充部分） -->
      <circle cx="50" cy="50" r="48" class="liquid-ball__bg" />

      <!-- 水位 + 两层波浪，用 clip 裁出圆形 -->
      <g clip-path="url(#liquid-clip)">
        <g :style="{ transform: `translateY(${(1 - clamped) * 100}%)` }" class="liquid-ball__wave-group">
          <path
            class="liquid-ball__wave liquid-ball__wave--back" :fill="color" fill-opacity="0.5"
            d="M0 50 Q 12.5 35 25 50 T 50 50 T 75 50 T 100 50 T 125 50 T 150 50 V100 H0 Z"
          />
          <path
            class="liquid-ball__wave liquid-ball__wave--front" :fill="color"
            d="M0 50 Q 12.5 35 25 50 T 50 50 T 75 50 T 100 50 T 125 50 T 150 50 V100 H0 Z"
          />
        </g>
      </g>

      <!-- 球边框 -->
      <circle cx="50" cy="50" r="48" class="liquid-ball__border" />
    </svg>

    <!-- 居中百分比 -->
    <span class="liquid-ball__label" :style="{ color: textColor }">
      {{ percent }}%
    </span>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    /** 水位比例，0~1 */
    value: number
    /** 直径像素 */
    size?: number
    /** 主题色 */
    color?: string
    /** 文字颜色 */
    textColor?: string
  }>(),
  { size: 140, color: '#4fc3f7', textColor: '#ffffff' },
)

const clamped = computed(() => Math.min(1, Math.max(0, props.value)))
const percent = computed(() => Math.round(clamped.value * 100))
</script>

<style scoped lang="scss">
.liquid-ball {
  position: relative;
  display: inline-block;
}

.liquid-ball__svg {
  width: 100%;
  height: 100%;
  display: block;
}

.liquid-ball__bg {
  fill: rgba(79, 195, 247, 0.08);
}

.liquid-ball__border {
  fill: none;
  stroke: rgba(79, 195, 247, 0.6);
  stroke-width: 2;
}

// 水位整体上下平移由 transform 控制（GPU 合成，不占主线程）
.liquid-ball__wave-group {
  transition: transform 0.6s ease-out;
  will-change: transform;
}

// 两层波浪做水平往返平移，形成流动观感
.liquid-ball__wave {
  animation: liquid-wave 3.5s linear infinite;
}

.liquid-ball__wave--back {
  animation-duration: 4.5s;
  animation-direction: reverse;
  opacity: 0.6;
}

@keyframes liquid-wave {
  from {
    transform: translateX(0);
  }
  to {
    // 一个完整波长为 50（viewBox 单位），平移 -50 即可无缝循环
    transform: translateX(-50px);
  }
}

.liquid-ball__label {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 26px;
  font-weight: 700;
  text-shadow: 0 0 8px rgba(0, 0, 0, 0.4);
}

// 尊重无障碍：减少动态偏好时停掉动画
@media (prefers-reduced-motion: reduce) {
  .liquid-ball__wave {
    animation: none;
  }
}
</style>
