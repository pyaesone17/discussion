import { apiClient } from './client'
import type { Topic, CreateTopicRequest } from '../types'

export const getTopics = async (): Promise<Topic[]> => {
  const { data } = await apiClient.get('/topics')
  return data
}

export const searchTopics = async (query: string): Promise<Topic[]> => {
  const { data } = await apiClient.get('/topics/search', { params: { q: query } })
  return data
}

export const getTopic = async (id: number): Promise<Topic> => {
  const { data } = await apiClient.get(`/topics/${id}`)
  return data
}

export const createTopic = async (request: CreateTopicRequest): Promise<Topic> => {
  const { data } = await apiClient.post('/topics', request)
  return data
}
