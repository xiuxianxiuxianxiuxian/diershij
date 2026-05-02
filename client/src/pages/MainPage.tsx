import { useState, useEffect } from 'react'
import { useAuthStore } from '@/stores/authStore'
import { useGameStore } from '@/stores/gameStore'
import { wsService } from '@/services/websocket'
import Sidebar from '@/components/Sidebar'
import StatusBar from '@/components/StatusBar'
import CharacterTab from '@/components/tabs/CharacterTab'
import WorldTab from '@/components/tabs/WorldTab'
import SocialTab from '@/components/tabs/SocialTab'
import CombatTab from '@/components/tabs/CombatTab'
import SettingsTab from '@/components/tabs/SettingsTab'

type TabType = 'character' | 'world' | 'social' | 'combat' | 'settings'

export default function MainPage() {
  const [activeTab, setActiveTab] = useState<TabType>('character')
  const { entity, logout } = useAuthStore()
  const { isConnected, chatMessages } = useGameStore()

  useEffect(() => {
    return () => {
      wsService.disconnect()
    }
  }, [])

  const handleLogout = () => {
    wsService.disconnect()
    logout()
  }

  const renderTab = () => {
    switch (activeTab) {
      case 'character':
        return <CharacterTab />
      case 'world':
        return <WorldTab />
      case 'social':
        return <SocialTab />
      case 'combat':
        return <CombatTab />
      case 'settings':
        return <SettingsTab onLogout={handleLogout} />
      default:
        return <CharacterTab />
    }
  }

  const tabs: { id: TabType; label: string; icon: string }[] = [
    { id: 'character', label: '角色', icon: '👤' },
    { id: 'world', label: '世界', icon: '🌍' },
    { id: 'social', label: '社交', icon: '💬' },
    { id: 'combat', label: '战斗', icon: '⚔️' },
    { id: 'settings', label: '设置', icon: '⚙️' },
  ]

  return (
    <div className="h-screen flex flex-col bg-gray-900">
      <StatusBar
        entity={entity}
        isConnected={isConnected}
        chatCount={chatMessages.length}
      />

      <div className="flex-1 flex overflow-hidden">
        <div className="w-48 bg-gray-850 border-r border-gray-700">
          <nav className="p-2 space-y-1">
            {tabs.map((tab) => (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id)}
                className={`w-full flex items-center gap-2 px-3 py-2 rounded-lg text-left transition-colors ${
                  activeTab === tab.id
                    ? 'bg-primary-600 text-white'
                    : 'text-gray-300 hover:bg-gray-700'
                }`}
              >
                <span>{tab.icon}</span>
                <span>{tab.label}</span>
              </button>
            ))}
          </nav>
        </div>

        <div className="flex-1 flex overflow-hidden">
          <main className="flex-1 overflow-auto p-4">
            {renderTab()}
          </main>

          <Sidebar />
        </div>
      </div>
    </div>
  )
}
