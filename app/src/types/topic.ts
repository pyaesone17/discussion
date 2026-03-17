export interface Topic {
  id: number
  title: string
  user_id: number
  username?: string
  score: number
  created_at: string
  updated_at: string
}

export interface CreateTopicRequest {
  title: string
  user_id: number
}

export interface UpdateTopicRequest {
  title?: string
}
