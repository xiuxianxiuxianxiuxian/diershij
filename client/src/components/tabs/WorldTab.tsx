import { useGameStore } from '@/stores/gameStore'
import { wsService } from '@/services/websocket'

export default function WorldTab() {
  const { currentRegion, worldTime } = useGameStore()

  const handleExplore = () => {
    wsService.send('operation', {
      action_type: 'explore',
      params: {},
    })
  }

  const handleGather = () => {
    wsService.send('operation', {
      action_type: 'gather',
      params: {},
    })
  }

  const formatWorldTime = (time: number) => {
    const date = new Date(time * 1000)
    return date.toLocaleString('zh-CN')
  }

  return (
    <div className="space-y-6">
      <div className="card">
        <h2 className="text-xl font-bold text-immortal-gold mb-4">世界信息</h2>
        
        <div className="grid grid-cols-2 gap-4">
          <div>
            <span className="text-gray-400">世界时间：</span>
            <span className="text-white">{formatWorldTime(worldTime)}</span>
          </div>
        </div>
      </div>

      <div className="card">
        <h3 className="text-lg font-bold text-white mb-4">当前位置</h3>
        
        {currentRegion ? (
          <div className="space-y-4">
            <div>
              <span className="text-gray-400">区域：</span>
              <span className="text-white text-lg">{currentRegion.name}</span>
            </div>
            
            <p className="text-gray-300">{currentRegion.description}</p>
            
            <div className="grid grid-cols-3 gap-4">
              <div className="text-center">
                <div className="text-xl font-bold text-immortal-jade">
                  {currentRegion.spiritual_density.toFixed(0)}%
                </div>
                <div className="text-sm text-gray-400">灵气浓度</div>
              </div>
              <div className="text-center">
                <div className="text-xl font-bold text-immortal-purple">
                  {currentRegion.spiritual_tier}
                </div>
                <div className="text-sm text-gray-400">灵气品阶</div>
              </div>
              <div className="text-center">
                <div className="text-xl font-bold text-red-400">
                  {currentRegion.danger_level}
                </div>
                <div className="text-sm text-gray-400">危险等级</div>
              </div>
            </div>

            {currentRegion.resources.length > 0 && (
              <div>
                <h4 className="text-sm font-medium text-gray-400 mb-2">资源</h4>
                <div className="flex flex-wrap gap-2">
                  {currentRegion.resources.map((res) => (
                    <span
                      key={res.id}
                      className="px-2 py-1 bg-gray-700 rounded text-sm text-gray-300"
                    >
                      {res.name} x{res.quantity}
                    </span>
                  ))}
                </div>
              </div>
            )}
          </div>
        ) : (
          <div className="text-gray-400">加载区域信息...</div>
        )}
      </div>

      <div className="card">
        <h3 className="text-lg font-bold text-white mb-4">探索</h3>
        
        <div className="grid grid-cols-2 gap-4">
          <button onClick={handleExplore} className="btn-primary">
            探索区域
          </button>
          <button onClick={handleGather} className="btn-secondary">
            采集资源
          </button>
        </div>
      </div>

      <div className="card">
        <h3 className="text-lg font-bold text-white mb-4">世界地图</h3>
        
        <div className="text-gray-400 text-center py-8">
          世界地图功能开发中...
        </div>
      </div>
    </div>
  )
}
