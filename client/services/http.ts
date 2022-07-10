import axios from 'axios'
import auth from './auth'

axios.defaults.headers.post['content-type'] = 'application/json'
axios.defaults.headers.get['content-type'] = 'application/json'

export const http = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_BASE_URL,
})

http.interceptors.request.use(
  async (config) => {
    const token = auth.getAccessToken()
    if (token) {
      config.headers = {
        Authorization: `Bearer ${token}`,
        Accept: 'application/json',
        'Content-Type': 'application/json',
      }
    }

    return config
  },
  (error) => {
    Promise.reject(error)
  }
)

http.interceptors.response.use(
  (response) => {
    return response
  },
  async function (error) {
    if (error.response.status === 403) {
      auth.logout()
      window.location.href = '/login'
    }

    return Promise.reject(error)
  }
)
