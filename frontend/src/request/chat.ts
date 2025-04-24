import axiosInstance from './axios'

// 创建会话
export const createSession = async (data: any) => {
  const res = await axiosInstance.post('/session', data)
  return res.data
}

