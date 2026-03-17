import { apiClient } from './client'

export const castVote = async (
  votableType: 'topic' | 'post',
  votableId: number,
  userId: number,
  value: 1 | -1
): Promise<void> => {
  const resource = votableType === 'topic' ? 'topics' : 'posts'
  await apiClient.post(`/${resource}/${votableId}/vote`, {
    user_id: userId,
    value,
  })
}

export const removeVote = async (
  votableType: 'topic' | 'post',
  votableId: number,
  userId: number
): Promise<void> => {
  const resource = votableType === 'topic' ? 'topics' : 'posts'
  await apiClient.delete(`/${resource}/${votableId}/vote`, {
    data: { user_id: userId },
  })
}
