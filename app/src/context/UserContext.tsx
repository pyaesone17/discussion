import { createContext, useContext, useState, useEffect, useCallback, useMemo, ReactNode } from 'react'
import type { User } from '../types'
import { STORAGE_KEYS } from '../utils/constants'

interface UserContextType {
  currentUser: User | null
  setCurrentUser: (user: User | null) => void
  isUserSelected: boolean
}

const UserContext = createContext<UserContextType | undefined>(undefined)

export function UserProvider({ children }: { children: ReactNode }) {
  const [currentUser, setCurrentUserState] = useState<User | null>(null)

  // Load user from localStorage on mount
  useEffect(() => {
    const stored = localStorage.getItem(STORAGE_KEYS.CURRENT_USER)
    if (stored) {
      try {
        setCurrentUserState(JSON.parse(stored))
      } catch (error) {
        console.error('Failed to parse stored user:', error)
        localStorage.removeItem(STORAGE_KEYS.CURRENT_USER)
      }
    }
  }, [])

  // Save user to localStorage whenever it changes
  const setCurrentUser = useCallback((user: User | null) => {
    setCurrentUserState(user)
    if (user) {
      localStorage.setItem(STORAGE_KEYS.CURRENT_USER, JSON.stringify(user))
    } else {
      localStorage.removeItem(STORAGE_KEYS.CURRENT_USER)
    }
  }, [])

  const value = useMemo(
    () => ({
      currentUser,
      setCurrentUser,
      isUserSelected: currentUser !== null,
    }),
    [currentUser, setCurrentUser]
  )

  return <UserContext.Provider value={value}>{children}</UserContext.Provider>
}

export function useUser() {
  const context = useContext(UserContext)
  if (!context) {
    throw new Error('useUser must be used within UserProvider')
  }
  return context
}
