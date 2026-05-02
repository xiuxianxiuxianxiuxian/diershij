import type { Entity } from '@/types'

interface StatusBarProps {
  entity: Entity | null
  isConnected: boolean
  chatCount: number
}

export default function StatusBar({ entity, isConnected, chatCount }: StatusBarProps) {
  const realmNames: Record<string, string> = {
    mortal: '凡人',
    qi_condensation: '炼气',
    foundation: '筑基',
    golden_core: '金丹',
    nascent_soul: '元婴',
    soul_transformation: '化神',
    void_refinement: '炼虚',
    integration: '合体',
    mahayana: '大乘',
    tribulation: '渡劫',
  }

  const getRealmClass = (realm: string) => {
    const classes: Record<string, string> = {
      mortal: 'realm-mortal',
      qi_condensation: 'realm-qi_condensation',
      foundation: 'realm-foundation',
      golden_core: 'realm-golden_core',
      nascent_soul: 'realm-nascent_soul',
      soul_transformation: 'realm-soul_transformation',
    }
    return classes[realm] || 'realm-mortal'
  }

  return (
    <div className="status-bar">
      <div className="flex items-center gap-4">
        <span className="text-immortal-gold font-bold">修仙世界</span>
        
        {entity && (
          <>
            <span className="text-gray-400">|</span>
            <span className="text-white">{entity.name}</span>
            <span className={`realm-badge ${getRealmClass(entity.realm)}`}>
              {realmNames[entity.realm] || entity.realm}
            </span>
          </>
        )}
      </div>

      <div className="flex items-center gap-4">
        <div className="flex items-center gap-2">
          <span
            className={`w-2 h-2 rounded-full ${
              isConnected ? 'bg-green-500' : 'bg-red-500'
            }`}
          />
          <span className="text-sm text-gray-400">
            {isConnected ? '已连接' : '断开'}
          </span>
        </div>

        <span className="text-sm text-gray-400">
          💬 {chatCount}
        </span>
      </div>
    </div>
  )
}
