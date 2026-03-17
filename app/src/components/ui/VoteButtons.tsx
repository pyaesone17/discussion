import { useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { castVote, removeVote } from '../../api/votes'
import { useUser } from '../../context/UserContext'
import { useToast } from '../../context/ToastContext'

interface VoteButtonsProps {
  votableType: 'topic' | 'post'
  votableId: number
  score: number
  queryKeyToInvalidate: unknown[]
}

export default function VoteButtons({
  votableType,
  votableId,
  score,
  queryKeyToInvalidate,
}: VoteButtonsProps) {
  const { currentUser } = useUser()
  const { showToast } = useToast()
  const queryClient = useQueryClient()
  const [userVote, setUserVote] = useState<1 | -1 | null>(null)

  const voteMutation = useMutation({
    mutationFn: async ({ value }: { value: 1 | -1 }) => {
      if (!currentUser) throw new Error('No user selected')

      if (userVote === value) {
        await removeVote(votableType, votableId, currentUser.id)
        return null
      } else {
        await castVote(votableType, votableId, currentUser.id, value)
        return value
      }
    },
    onSuccess: (newVote) => {
      setUserVote(newVote)
      queryClient.invalidateQueries({ queryKey: queryKeyToInvalidate })
    },
    onError: (error: Error) => {
      showToast(error.message, 'error')
    },
  })

  const handleVote = (value: 1 | -1) => {
    if (!currentUser) {
      showToast('Please select a user first', 'error')
      return
    }
    voteMutation.mutate({ value })
  }

  return (
    <div className="flex flex-col items-center gap-1">
      <button
        onClick={(e) => {
          e.preventDefault()
          e.stopPropagation()
          handleVote(1)
        }}
        disabled={voteMutation.isPending}
        className={`p-1 rounded transition-colors ${
          userVote === 1
            ? 'text-primary-600'
            : 'text-gray-400 hover:text-primary-600'
        }`}
        aria-label="Upvote"
      >
        <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
          <path
            fillRule="evenodd"
            d="M14.707 12.707a1 1 0 01-1.414 0L10 9.414l-3.293 3.293a1 1 0 01-1.414-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 010 1.414z"
            clipRule="evenodd"
          />
        </svg>
      </button>
      <span
        className={`text-sm font-semibold tabular-nums ${
          userVote === 1
            ? 'text-primary-600'
            : userVote === -1
              ? 'text-red-500'
              : 'text-gray-700'
        }`}
      >
        {score}
      </span>
      <button
        onClick={(e) => {
          e.preventDefault()
          e.stopPropagation()
          handleVote(-1)
        }}
        disabled={voteMutation.isPending}
        className={`p-1 rounded transition-colors ${
          userVote === -1
            ? 'text-red-500'
            : 'text-gray-400 hover:text-red-500'
        }`}
        aria-label="Downvote"
      >
        <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
          <path
            fillRule="evenodd"
            d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z"
            clipRule="evenodd"
          />
        </svg>
      </button>
    </div>
  )
}
