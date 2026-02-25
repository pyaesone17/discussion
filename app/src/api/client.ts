import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api/v1'

export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Response interceptor for consistent error handling
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    // Extract error message from API response {error: "message"}
    const message = error.response?.data?.error || error.message || 'An error occurred'
    return Promise.reject(new Error(message))
  }
)
