import { apiClient } from './client'
import type { Post, CreatePostRequest } from '../types'

export const getTopicPosts = async (topicId: number): Promise<Post[]> => {
  const { data } = await apiClient.get(`/topics/${topicId}/posts`)
  return data
}

export const createPost = async (request: CreatePostRequest): Promise<Post> => {
  const { data } = await apiClient.post('/posts', request)
  return data
}
