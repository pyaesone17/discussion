import { apiClient } from './client'

export const createReport = async (
  reportableType: 'topic' | 'post',
  reportableId: number,
  userId: number,
  reason: string
): Promise<void> => {
  const resource = reportableType === 'topic' ? 'topics' : 'posts'
  await apiClient.post(`/${resource}/${reportableId}/report`, {
    user_id: userId,
    reason,
  })
}
