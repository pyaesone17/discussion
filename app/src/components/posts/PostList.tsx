import type { Post } from '../../types'
import PostCard from './PostCard'
import LoadingSpinner from '../ui/LoadingSpinner'
import ErrorMessage from '../ui/ErrorMessage'

interface PostListProps {
  posts: Post[]
  isLoading: boolean
  error: Error | null
  refetch?: () => void
}

export default function PostList({ posts, isLoading, error, refetch }: PostListProps) {
  if (isLoading) {
    return <LoadingSpinner size="lg" />
  }

  if (error) {
    return <ErrorMessage message={error.message} retry={refetch} />
  }

  if (posts.length === 0) {
    return (
      <div className="text-center py-8 text-gray-600">
        <p>No posts yet. Be the first to reply!</p>
      </div>
    )
  }

  return (
    <div className="space-y-4">
      {posts.map((post) => (
        <PostCard key={post.id} post={post} />
      ))}
    </div>
  )
}
