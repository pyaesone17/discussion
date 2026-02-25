import { useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { createPost } from '../../api/posts'
import { useUser } from '../../context/UserContext'
import { useToast } from '../../context/ToastContext'
import Textarea from '../ui/Textarea'
import Button from '../ui/Button'

interface CreatePostFormProps {
  topicId: number
}

export default function CreatePostForm({ topicId }: CreatePostFormProps) {
  const [content, setContent] = useState('')
  const [error, setError] = useState('')
  const { currentUser } = useUser()
  const { showToast } = useToast()
  const queryClient = useQueryClient()

  const mutation = useMutation({
    mutationFn: createPost,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['topics', topicId, 'posts'] })
      showToast('Post added successfully!', 'success')
      setContent('')
      setError('')
    },
    onError: (error: Error) => {
      showToast(error.message, 'error')
    },
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    setError('')

    if (!currentUser) {
      setError('Please select a user first')
      return
    }

    if (content.length < 1) {
      setError('Post cannot be empty')
      return
    }

    if (content.length > 5000) {
      setError('Post must be less than 5000 characters')
      return
    }

    mutation.mutate({
      topic_id: topicId,
      user_id: currentUser.id,
      content: content.trim(),
    })
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div>
        <label htmlFor="content" className="block text-sm font-medium text-gray-700 mb-2">
          Add a Reply
        </label>
        <Textarea
          id="content"
          value={content}
          onChange={(e) => setContent(e.target.value)}
          placeholder="Share your thoughts..."
          rows={4}
          error={error}
          disabled={mutation.isPending}
        />
        <p className="mt-1 text-sm text-gray-500">
          {content.length}/5000 characters
        </p>
      </div>

      <Button
        type="submit"
        isLoading={mutation.isPending}
        disabled={!currentUser || content.trim().length < 1}
      >
        Post Reply
      </Button>
    </form>
  )
}
