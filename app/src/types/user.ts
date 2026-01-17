export interface User {
  id: number
  username: string
  email: string
  created_at: string
  updated_at: string
}

export interface CreateUserRequest {
  username: string
  email: string
}

export interface UpdateUserRequest {
  username?: string
  email?: string
}
