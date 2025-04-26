import axiosInstance from './axios'

// 创建空会话
export const createEmptySession = async (data: any) => {
  const res = await axiosInstance.post('/create-empty-session', data)
  return res.data
}

// 创建会话并发送消息
export const createSessionWithMessage = async (data: any) => {
  const res = await axiosInstance.post('/create-session-with-message', data)
  return res.data
}

// 获取会话列表
export const getSessionList = async () => {
  const res = await axiosInstance.get('/get-session-list')
  return res.data
}