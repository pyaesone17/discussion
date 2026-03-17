import { useQuery } from '@tanstack/react-query'
import { useParams, Link } from 'react-router-dom'
import { getTopic } from '../api/topics'
import { getTopicPosts } from '../api/posts'
import { formatRelativeTime } from '../utils/date'
import PostList from '../components/posts/PostList'
import CreatePostForm from '../components/posts/CreatePostForm'
import LoadingSpinner from '../components/ui/LoadingSpinner'
import ErrorMessage from '../components/ui/ErrorMessage'
import VoteButtons from '../components/ui/VoteButtons'
import ReportButton from '../components/ui/ReportButton'

export default function TopicDetailPage() {
  const { id } = useParams<{ id: string }>()
  const topicId = parseInt(id || '0', 10)

  // Validate topic ID early
  if (!id || isNaN(topicId) || topicId <= 0) {
    return <ErrorMessage message="Invalid topic ID" />
  }

  const {
    data: topic,
    isLoading: topicLoading,
    error: topicError,
    refetch: refetchTopic,
  } = useQuery({
    queryKey: ['topics', topicId],
    queryFn: () => getTopic(topicId),
    staleTime: 5 * 60 * 1000, // 5 minutes
  })

  const {
    data: posts,
    isLoading: postsLoading,
    error: postsError,
    refetch: refetchPosts,
  } = useQuery({
    queryKey: ['topics', topicId, 'posts'],
    queryFn: () => getTopicPosts(topicId),
    staleTime: 3 * 60 * 1000, // 3 minutes
  })

  if (topicLoading) {
    return <LoadingSpinner size="lg" />
  }

  if (topicError) {
    return <ErrorMessage message={topicError.message} retry={refetchTopic} />
  }

  if (!topic) {
    return <ErrorMessage message="Topic not found" />
  }

  return (
    <div className="space-y-6">
      {/* Breadcrumb */}
      <Link
        to="/"
        className="text-sm text-primary-600 hover:text-primary-700 inline-block"
      >
        ← Back to all topics
      </Link>

      {/* Topic Header */}
      <div className="bg-white rounded-lg shadow-sm p-6 border border-gray-200 flex gap-4">
        <VoteButtons
          votableType="topic"
          votableId={topic.id}
          score={topic.score}
          queryKeyToInvalidate={['topics', topicId]}
        />
        <div className="flex-1 min-w-0">
          <h1 className="text-3xl font-bold text-gray-900 mb-4">{topic.title}</h1>
          <div className="flex items-center gap-4 text-sm text-gray-600">
            <div className="flex items-center gap-1">
              <svg
                className="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
                />
              </svg>
              <span className="font-medium">{topic.username || 'Unknown'}</span>
            </div>
            <span className="text-gray-400">·</span>
            <div className="flex items-center gap-1">
              <svg
                className="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              <span>{formatRelativeTime(topic.created_at)}</span>
            </div>
            <div className="ml-auto">
              <ReportButton reportableType="topic" reportableId={topic.id} />
            </div>
          </div>
        </div>
      </div>

      {/* Posts Section */}
      <div className="space-y-4">
        <h2 className="text-2xl font-bold text-gray-900">
          Replies ({posts?.length || 0})
        </h2>
        <PostList
          posts={posts || []}
          isLoading={postsLoading}
          error={postsError}
          refetch={refetchPosts}
        />
      </div>

      {/* Create Post Form */}
      <div className="bg-white rounded-lg shadow-sm p-6 border border-gray-200">
        <CreatePostForm topicId={topicId} />
      </div>
    </div>
  )
}
