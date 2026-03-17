import type { Post } from '../../types'
import { formatRelativeTime } from '../../utils/date'
import VoteButtons from '../ui/VoteButtons'
import ReportButton from '../ui/ReportButton'

interface PostCardProps {
  post: Post
}

export default function PostCard({ post }: PostCardProps) {
  return (
    <div className="bg-white rounded-lg shadow-sm p-6 border border-gray-200 flex gap-4">
      <VoteButtons
        votableType="post"
        votableId={post.id}
        score={post.score}
        queryKeyToInvalidate={['topics', post.topic_id, 'posts']}
      />
      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-2 mb-3">
          <svg
            className="w-5 h-5 text-gray-500"
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
          <span className="font-medium text-gray-900">{post.username || 'Unknown'}</span>
          <span className="text-gray-400">·</span>
          <span className="text-sm text-gray-500">{formatRelativeTime(post.created_at)}</span>
        </div>
        <p className="text-gray-800 whitespace-pre-wrap leading-relaxed">{post.content}</p>
      </div>
      <div className="flex-shrink-0 self-start">
        <ReportButton reportableType="post" reportableId={post.id} />
      </div>
    </div>
  )
}
