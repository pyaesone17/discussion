import { useQuery } from '@tanstack/react-query'
import { useSearchParams, Link } from 'react-router-dom'
import { getTopics, searchTopics } from '../api/topics'
import TopicList from '../components/topics/TopicList'
import Button from '../components/ui/Button'

export default function HomePage() {
  const [searchParams] = useSearchParams()
  const query = searchParams.get('q') || ''

  const { data: topics, isLoading, error, refetch } = useQuery({
    queryKey: query ? ['topics', 'search', query] : ['topics'],
    queryFn: () => (query ? searchTopics(query) : getTopics()),
    staleTime: 2 * 60 * 1000, // 2 minutes (matches API cache)
  })

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">
            {query ? `Search Results for "${query}"` : 'All Topics'}
          </h1>
          {query && (
            <Link
              to="/"
              className="text-sm text-primary-600 hover:text-primary-700 mt-2 inline-block"
            >
              ← Back to all topics
            </Link>
          )}
        </div>
        <Link to="/topics/new">
          <Button>Create New Topic</Button>
        </Link>
      </div>

      <TopicList
        topics={topics || []}
        isLoading={isLoading}
        error={error}
        refetch={refetch}
      />
    </div>
  )
}
