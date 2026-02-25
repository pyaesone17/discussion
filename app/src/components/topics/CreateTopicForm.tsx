import { useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useNavigate } from 'react-router-dom'
import { createTopic } from '../../api/topics'
import { useUser } from '../../context/UserContext'
import { useToast } from '../../context/ToastContext'
import Input from '../ui/Input'
import Button from '../ui/Button'

export default function CreateTopicForm() {
  const [title, setTitle] = useState('')
  const [error, setError] = useState('')
  const { currentUser } = useUser()
  const { showToast } = useToast()
  const navigate = useNavigate()
  const queryClient = useQueryClient()

  const mutation = useMutation({
    mutationFn: createTopic,
    onSuccess: (newTopic) => {
      queryClient.invalidateQueries({ queryKey: ['topics'] })
      showToast('Topic created successfully!', 'success')
      navigate(`/topics/${newTopic.id}`)
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

    if (title.length < 5) {
      setError('Title must be at least 5 characters')
      return
    }

    if (title.length > 200) {
      setError('Title must be less than 200 characters')
      return
    }

    mutation.mutate({
      title: title.trim(),
      user_id: currentUser.id,
    })
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div>
        <label htmlFor="title" className="block text-sm font-medium text-gray-700 mb-2">
          Topic Title
        </label>
        <Input
          id="title"
          type="text"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          placeholder="Enter a descriptive title (5-200 characters)"
          error={error}
          disabled={mutation.isPending}
          autoFocus
        />
        <p className="mt-1 text-sm text-gray-500">
          {title.length}/200 characters
        </p>
      </div>

      <div className="flex gap-3">
        <Button
          type="submit"
          isLoading={mutation.isPending}
          disabled={!currentUser || title.trim().length < 5}
        >
          Create Topic
        </Button>
        <Button
          type="button"
          variant="secondary"
          onClick={() => navigate('/')}
          disabled={mutation.isPending}
        >
          Cancel
        </Button>
      </div>
    </form>
  )
}
