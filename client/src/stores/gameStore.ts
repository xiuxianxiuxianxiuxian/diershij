import { create } from 'zustand'
import type { ChatMessage, Region } from '@/types'

interface GameState {
  currentRegion: Region | null
  nearbyEntities: string[]
  chatMessages: ChatMessage[]
  worldTime: number
  isConnected: boolean
  setCurrentRegion: (region: Region | null) => void
  setNearbyEntities: (entities: string[]) => void
  addChatMessage: (message: ChatMessage) => void
  setWorldTime: (time: number) => void
  setConnected: (connected: boolean) => void
}

export const useGameStore = create<GameState>((set) => ({
  currentRegion: null,
  nearbyEntities: [],
  chatMessages: [],
  worldTime: 0,
  isConnected: false,
  setCurrentRegion: (region) => set({ currentRegion: region }),
  setNearbyEntities: (entities) => set({ nearbyEntities: entities }),
  addChatMessage: (message) =>
    set((state) => ({
      chatMessages: [...state.chatMessages.slice(-100), message],
    })),
  setWorldTime: (time) => set({ worldTime: time }),
  setConnected: (connected) => set({ isConnected: connected }),
}))
