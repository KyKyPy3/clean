import { create } from 'zustand'
import { persist, createJSONStorage } from 'zustand/middleware'

interface Store {
  email: string
  token: string
  setSession: (response: any) => void
  clearSession: () => void
}

const initialState: Omit<Store, 'clearSession' | 'setSession'> = {
  email: '',
  token: ''
}

export const useSessionStore = create(
  persist<Store>(
    (set) => ({
      ...initialState,
      setSession: (response) => set(response),
      clearSession: () => set(initialState)
    }),
    {
      name: 'session-store',
      storage: createJSONStorage(() => localStorage)
    }
  )
)