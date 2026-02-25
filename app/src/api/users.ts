import { apiClient } from './client'
import type { User } from '../types'

export const getUsers = async (): Promise<User[]> => {
  const { data } = await apiClient.get('/users')
  return data
}
