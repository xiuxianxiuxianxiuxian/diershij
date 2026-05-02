import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { Entity } from '@/types'

interface AuthState {
  token: string | null
  entity: Entity | null
  isAuthenticated: boolean
  setAuth: (token: string, entity: Entity) => void
  logout: () => void
  updateEntity: (entity: Partial<Entity>) => void
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      token: null,
      entity: null,
      isAuthenticated: false,
      setAuth: (token, entity) =>
        set({
          token,
          entity,
          isAuthenticated: true,
        }),
      logout: () =>
        set({
          token: null,
          entity: null,
          isAuthenticated: false,
        }),
      updateEntity: (updates) =>
        set((state) => ({
          entity: state.entity ? { ...state.entity, ...updates } : null,
        })),
    }),
    {
      name: 'auth-storage',
    }
  )
)
