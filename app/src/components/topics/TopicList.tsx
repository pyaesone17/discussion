import type { Topic } from '../../types'
import TopicCard from './TopicCard'
import LoadingSpinner from '../ui/LoadingSpinner'
import ErrorMessage from '../ui/ErrorMessage'

interface TopicListProps {
  topics: Topic[]
  isLoading: boolean
  error: Error | null
  refetch?: () => void
}

export default function TopicList({ topics, isLoading, error, refetch }: TopicListProps) {
  if (isLoading) {
    return <LoadingSpinner size="lg" />
  }

  if (error) {
    return <ErrorMessage message={error.message} retry={refetch} />
  }

  if (topics.length === 0) {
    return (
      <div className="text-center py-12">
        <svg
          className="w-16 h-16 text-gray-300 mx-auto mb-4"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth="2"
            d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4"
          />
        </svg>
        <h3 className="text-lg font-medium text-gray-900 mb-2">No topics found</h3>
        <p className="text-gray-600">Be the first to create a topic!</p>
      </div>
    )
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      {topics.map((topic) => (
        <TopicCard key={topic.id} topic={topic} />
      ))}
    </div>
  )
}
