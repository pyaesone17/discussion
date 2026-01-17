import { useQuery } from '@tanstack/react-query'
import { getUsers } from '../../api/users'
import { useUser } from '../../context/UserContext'
import type { User } from '../../types'
import Modal from '../ui/Modal'
import LoadingSpinner from '../ui/LoadingSpinner'
import ErrorMessage from '../ui/ErrorMessage'

export default function UserSelectModal() {
  const { currentUser, setCurrentUser, isUserSelected } = useUser()
  const { data: users, isLoading, error, refetch } = useQuery({
    queryKey: ['users'],
    queryFn: getUsers,
    enabled: !isUserSelected,
  })

  const handleSelectUser = (user: User) => {
    setCurrentUser(user)
  }

  return (
    <Modal
      isOpen={!isUserSelected}
      onClose={() => {}}
      title="Select Your User"
    >
      <div className="space-y-4">
        <p className="text-gray-600 text-sm">
          Choose a user to continue. This selection will be remembered for your next visit.
        </p>

        {isLoading && <LoadingSpinner size="md" />}

        {error && <ErrorMessage message={error.message} retry={refetch} />}

        {users && (
          <div className="grid grid-cols-1 gap-3">
            {users.map((user) => (
              <button
                key={user.id}
                onClick={() => handleSelectUser(user)}
                className="flex items-center gap-3 p-4 border border-gray-200 rounded-lg hover:border-primary-500 hover:bg-primary-50 transition-all duration-200 text-left"
              >
                <div className="flex-shrink-0 w-10 h-10 bg-primary-100 rounded-full flex items-center justify-center">
                  <svg
                    className="w-6 h-6 text-primary-600"
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
                </div>
                <div className="flex-1">
                  <p className="font-medium text-gray-900">{user.username}</p>
                  <p className="text-sm text-gray-500">{user.email}</p>
                </div>
              </button>
            ))}
          </div>
        )}
      </div>
    </Modal>
  )
}
