import axios from 'axios'
// 创建axios实例
const service = axios.create({
  // baseURL: 'http://host.docker.internal:4396', // api的base_url
  baseURL: 'http://localhost:4396', // api的base_url
  timeout: 20000, // 请求超时时间
})
export default service