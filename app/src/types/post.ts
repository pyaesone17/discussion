export interface Post {
  id: number
  topic_id: number
  user_id: number
  username?: string
  content: string
  score: number
  created_at: string
  updated_at: string
}

export interface CreatePostRequest {
  topic_id: number
  user_id: number
  content: string
}

export interface UpdatePostRequest {
  content?: string
}
