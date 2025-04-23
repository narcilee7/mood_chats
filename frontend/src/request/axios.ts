import axios from "axios";

const instance = axios.create({
  baseURL: "http://localhost:8081/api",
});


instance.interceptors.request.use(
  async (config) => {
    // todo
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);


instance.interceptors.response.use(
  (response) => {
    return response;
  },
  async (error) => {
    return Promise.reject(error);
  }
);

export default instance;