import { useAuthStore } from '@/stores/authStore'
import { wsService } from '@/services/websocket'

export default function CharacterTab() {
  const { entity } = useAuthStore()

  if (!entity) {
    return <div className="text-gray-400">加载中...</div>
  }

  const handleCultivate = () => {
    wsService.send('operation', {
      action_type: 'cultivate',
      params: {},
    })
  }

  const handleMeditate = () => {
    wsService.send('operation', {
      action_type: 'meditate',
      params: {},
    })
  }

  const handleSleep = () => {
    wsService.send('operation', {
      action_type: 'sleep',
      params: {},
    })
  }

  const handleBreakthrough = () => {
    wsService.send('operation', {
      action_type: 'breakthrough',
      params: {},
    })
  }

  const realmNames: Record<string, string> = {
    mortal: '凡人',
    qi_condensation: '炼气期',
    foundation: '筑基期',
    golden_core: '金丹期',
    nascent_soul: '元婴期',
    soul_transformation: '化神期',
    void_refinement: '炼虚期',
    integration: '合体期',
    mahayana: '大乘期',
    tribulation: '渡劫期',
  }

  const qiPercent = (entity.attributes.qi / entity.attributes.max_qi) * 100
  const spPercent = (entity.attributes.spiritual_power / entity.attributes.max_spiritual_power) * 100
  const progressPercent = entity.attributes.cultivation_progress

  return (
    <div className="space-y-6">
      <div className="card">
        <h2 className="text-xl font-bold text-immortal-gold mb-4">角色信息</h2>
        
        <div className="grid grid-cols-2 gap-4">
          <div>
            <span className="text-gray-400">道号：</span>
            <span className="text-white">{entity.name}</span>
          </div>
          <div>
            <span className="text-gray-400">境界：</span>
            <span className="text-immortal-purple">{realmNames[entity.realm]}</span>
          </div>
          <div>
            <span className="text-gray-400">寿命：</span>
            <span className="text-white">
              {entity.attributes.remaining_lifespan} / {entity.attributes.max_lifespan} 年
            </span>
          </div>
          <div>
            <span className="text-gray-400">状态：</span>
            <span className="text-green-400">{entity.status}</span>
          </div>
        </div>
      </div>

      <div className="card">
        <h3 className="text-lg font-bold text-white mb-4">属性</h3>
        
        <div className="space-y-4">
          <div>
            <div className="flex justify-between text-sm mb-1">
              <span className="text-gray-400">气血</span>
              <span className="text-white">
                {entity.attributes.qi.toFixed(0)} / {entity.attributes.max_qi.toFixed(0)}
              </span>
            </div>
            <div className="h-2 bg-gray-700 rounded-full overflow-hidden">
              <div
                className="h-full bg-red-500 transition-all"
                style={{ width: `${qiPercent}%` }}
              />
            </div>
          </div>

          <div>
            <div className="flex justify-between text-sm mb-1">
              <span className="text-gray-400">灵力</span>
              <span className="text-white">
                {entity.attributes.spiritual_power.toFixed(0)} / {entity.attributes.max_spiritual_power.toFixed(0)}
              </span>
            </div>
            <div className="h-2 bg-gray-700 rounded-full overflow-hidden">
              <div
                className="h-full bg-blue-500 transition-all"
                style={{ width: `${spPercent}%` }}
              />
            </div>
          </div>

          <div>
            <div className="flex justify-between text-sm mb-1">
              <span className="text-gray-400">修为</span>
              <span className="text-white">{progressPercent.toFixed(1)}%</span>
            </div>
            <div className="h-2 bg-gray-700 rounded-full overflow-hidden">
              <div
                className="h-full bg-immortal-gold transition-all"
                style={{ width: `${progressPercent}%` }}
              />
            </div>
          </div>
        </div>

        <div className="grid grid-cols-3 gap-4 mt-6">
          <div className="text-center">
            <div className="text-2xl font-bold text-immortal-jade">
              {entity.attributes.comprehension}
            </div>
            <div className="text-sm text-gray-400">悟性</div>
          </div>
          <div className="text-center">
            <div className="text-2xl font-bold text-immortal-purple">
              {entity.attributes.constitution}
            </div>
            <div className="text-sm text-gray-400">根骨</div>
          </div>
          <div className="text-center">
            <div className="text-2xl font-bold text-immortal-gold">
              {entity.attributes.luck}
            </div>
            <div className="text-sm text-gray-400">气运</div>
          </div>
        </div>
      </div>

      <div className="card">
        <h3 className="text-lg font-bold text-white mb-4">行动</h3>
        
        <div className="grid grid-cols-2 gap-4">
          <button onClick={handleCultivate} className="btn-primary">
            修炼
          </button>
          <button onClick={handleMeditate} className="btn-secondary">
            打坐
          </button>
          <button onClick={handleSleep} className="btn-secondary">
            休息
          </button>
          <button
            onClick={handleBreakthrough}
            disabled={progressPercent < 100}
            className="btn-primary disabled:opacity-50"
          >
            突破
          </button>
        </div>
      </div>

      <div className="card">
        <h3 className="text-lg font-bold text-white mb-4">战斗属性</h3>
        
        <div className="grid grid-cols-3 gap-4">
          <div className="text-center">
            <div className="text-xl font-bold text-red-400">
              {entity.attributes.attack_power.toFixed(1)}
            </div>
            <div className="text-sm text-gray-400">攻击</div>
          </div>
          <div className="text-center">
            <div className="text-xl font-bold text-blue-400">
              {entity.attributes.defense.toFixed(1)}
            </div>
            <div className="text-sm text-gray-400">防御</div>
          </div>
          <div className="text-center">
            <div className="text-xl font-bold text-green-400">
              {entity.attributes.speed.toFixed(1)}
            </div>
            <div className="text-sm text-gray-400">速度</div>
          </div>
        </div>
      </div>
    </div>
  )
}
