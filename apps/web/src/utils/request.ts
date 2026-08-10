import axios from 'axios'
import { ElMessage } from 'element-plus'
import 'element-plus/es/components/message/style/css'

//创建axios实例
const request = axios.create({
  baseURL: '/api',
  timeout: 5000,
})
//请求拦截器
request.interceptors.request.use()

//响应拦截器
request.interceptors.response.use(
  (response) => {
    // 成功回调
    /* 判断服务返回的 code（按业务域分段）
       200 -> 请求成功
       99xxxx -> 通用业务错误码段
         990001 -> 请求参数错误
         990002 -> 服务繁忙，服务内部错误
       110xxx -> 用户/权限模块错误码段
         110001 -> 用户名已存在，用于创建用户
         110002 -> 用户名不存在，用于登录
         110003 -> 用户名或密码错误，用于登录
         110004 -> 无效的 Token
         110005 -> 无权访问，需要登录
       120xxx -> 商品分类模块错误码段
         120001 -> 该节点下有子节点，不可以删除
    */
    // HTTP 层系统错误：系统崩溃/未捕获异常返回 500，body 仍带 code
    if (response.status === 500) {
      ElMessage({ type: 'error', message: '服务器内部错误，请稍后重试' })
      return Promise.reject(new Error('服务器内部错误'))
    }
    const code = response.data.code
    if (code !== 200) {
      // 提示错误信息
      ElMessage({
        type: 'error',
        message: response.data.message,
      })
      // 抛出错误
      return Promise.reject(new Error(response.data.message))
    }

    // 返回数据
    return response.data
  },
  (error) => {
    // 失败回调：处理 http 网络错误（真实后端业务错误已在成功回调按 code 处理）
    let msg = '网络异常，请稍后重试'
    if (error.response) {
      // 有 HTTP 响应但状态码非 2xx（防御性，真实后端一般返回 200）
      switch (error.response.status) {
        case 401:
          msg = '未授权，请重新登录'
          break
        case 403:
          msg = '拒绝访问'
          break
        case 404:
          msg = '请求地址不存在'
          break
        case 500:
          msg = '服务器内部错误'
          break
        default:
          msg = `请求错误(${error.response.status})`
      }
    } else if (error.code === 'ECONNABORTED') {
      msg = '请求超时，请稍后重试' // 命中你配的 timeout: 5000
    } else if (!window.navigator.onLine) {
      msg = '网络已断开，请检查网络连接'
    }
    ElMessage({ type: 'error', message: msg })
    return Promise.reject(error)
  },
)
export default request
