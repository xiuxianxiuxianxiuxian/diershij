export default function CombatTab() {
  return (
    <div className="space-y-6">
      <div className="card">
        <h2 className="text-xl font-bold text-immortal-gold mb-4">战斗</h2>
        
        <div className="text-gray-400 text-center py-8">
          战斗功能开发中...
        </div>
      </div>

      <div className="card">
        <h3 className="text-lg font-bold text-white mb-4">技能</h3>
        
        <div className="text-gray-400 text-center py-4">
          尚未学习技能
        </div>
      </div>

      <div className="card">
        <h3 className="text-lg font-bold text-white mb-4">功法</h3>
        
        <div className="text-gray-400 text-center py-4">
          尚未修炼功法
        </div>
        
        <button className="btn-secondary w-full">
          自创功法（需要10000极品灵石）
        </button>
      </div>

      <div className="card">
        <h3 className="text-lg font-bold text-white mb-4">战斗记录</h3>
        
        <div className="text-gray-400 text-center py-4">
          暂无战斗记录
        </div>
      </div>
    </div>
  )
}
