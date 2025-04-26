import axiosInstance from "./axios"

const login = async () => {
  const res = await axiosInstance.get('/login')
  return res.data
}



export { login }