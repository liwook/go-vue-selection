import { onUnmounted, ref } from 'vue'
import * as mock from './mockData'

/**
 * 大屏统一取数入口。
 * 当前直接返回本地 mock；将来接后端时，仅需替换内部数据来源
 * （轮询 / WebSocket），页面组件无需改动。
 */
export function useScreenData() {
  // 当前：直接返回 mock（整体对象，组件内用 data.value.xxx 取数）
  const data = ref(mock)

  // 预留：后端接入位（轮询 / WebSocket）
  // function startPolling() {
  //   setInterval(() => {
  //     fetch('/api/screen')
  //       .then((r) => r.json())
  //       .then((d) => (data.value = d))
  //   }, 5000)
  // }
  // function connectWS() {
  //   const ws = new WebSocket('/ws/screen')
  //   ws.onmessage = (e) => (data.value = JSON.parse(e.data))
  // }

  onUnmounted(() => {
    // 接后端时在此清理定时器 / 关闭 ws
  })

  return { data }
}
