import { useAuthStore } from '@/stores/authStore'

export default function SocialTab() {
  const { entity } = useAuthStore()

  return (
    <div className="space-y-6">
      <div className="card">
        <h2 className="text-xl font-bold text-immortal-gold mb-4">社交</h2>
        
        <div className="text-gray-400 text-center py-8">
          社交功能开发中...
        </div>
      </div>

      <div className="card">
        <h3 className="text-lg font-bold text-white mb-4">因果业力</h3>
        
        {entity && (
          <div className="grid grid-cols-3 gap-4">
            <div className="text-center">
              <div className="text-xl font-bold text-red-400">
                {entity.karma.karma_value}
              </div>
              <div className="text-sm text-gray-400">业力</div>
            </div>
            <div className="text-center">
              <div className="text-xl font-bold text-immortal-jade">
                {entity.karma.merit}
              </div>
              <div className="text-sm text-gray-400">功德</div>
            </div>
            <div className="text-center">
              <div className="text-xl font-bold text-immortal-purple">
                {entity.karma.heavenly_mark}
              </div>
              <div className="text-sm text-gray-400">天道标记</div>
            </div>
          </div>
        )}
      </div>

      <div className="card">
        <h3 className="text-lg font-bold text-white mb-4">宗门</h3>
        
        <div className="text-gray-400 text-center py-4">
          尚未加入宗门
        </div>
        
        <button className="btn-secondary w-full">
          创建宗门
        </button>
      </div>

      <div className="card">
        <h3 className="text-lg font-bold text-white mb-4">关系</h3>
        
        <div className="text-gray-400 text-center py-4">
          暂无特殊关系
        </div>
      </div>
    </div>
  )
}
