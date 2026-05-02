interface SettingsTabProps {
  onLogout: () => void
}

export default function SettingsTab({ onLogout }: SettingsTabProps) {
  return (
    <div className="space-y-6">
      <div className="card">
        <h2 className="text-xl font-bold text-immortal-gold mb-4">设置</h2>
        
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-300 mb-2">
              服务器地址
            </label>
            <input
              type="text"
              defaultValue="http://localhost:8080"
              className="input-field"
              disabled
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-300 mb-2">
              WebSocket地址
            </label>
            <input
              type="text"
              defaultValue="ws://localhost:8080/ws"
              className="input-field"
              disabled
            />
          </div>
        </div>
      </div>

      <div className="card">
        <h3 className="text-lg font-bold text-white mb-4">显示设置</h3>
        
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <span className="text-gray-300">显示聊天时间戳</span>
            <input type="checkbox" defaultChecked className="w-4 h-4" />
          </div>
          
          <div className="flex items-center justify-between">
            <span className="text-gray-300">显示系统消息</span>
            <input type="checkbox" defaultChecked className="w-4 h-4" />
          </div>
        </div>
      </div>

      <div className="card">
        <h3 className="text-lg font-bold text-white mb-4">关于</h3>
        
        <div className="space-y-2 text-gray-400">
          <p>修仙世界 v0.1.0</p>
          <p>一个完全自主演化的修仙MUD世界</p>
          <p className="text-sm mt-4">
            NPC与现实玩家共享完全一致的操作接口，
            通过混合AI系统自主决策，实现真正的自治世界。
          </p>
        </div>
      </div>

      <div className="card">
        <button
          onClick={onLogout}
          className="w-full py-2 bg-red-600 hover:bg-red-700 text-white rounded-lg transition-colors"
        >
          退出登录
        </button>
      </div>
    </div>
  )
}
