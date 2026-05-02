import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { useAuthStore } from '@/stores/authStore'
import { login as loginApi } from '@/services/api'
import { wsService } from '@/services/websocket'

export default function LoginPage() {
  const navigate = useNavigate()
  const { setAuth } = useAuthStore()
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    try {
      const result = await loginApi(username, password)
      
      if (result.success && result.token && result.entity) {
        setAuth(result.token, result.entity as any)
        wsService.connect(result.token)
        navigate('/')
      } else {
        setError(result.message || '登录失败')
      }
    } catch (err) {
      setError('网络错误，请检查服务器连接')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-b from-gray-900 to-gray-800">
      <div className="card w-full max-w-md">
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-immortal-gold mb-2">修仙世界</h1>
          <p className="text-gray-400">踏入仙途，证道长生</p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-300 mb-1">
              道号
            </label>
            <input
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className="input-field"
              placeholder="请输入道号"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-300 mb-1">
              密令
            </label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="input-field"
              placeholder="请输入密令"
              required
            />
          </div>

          {error && (
            <div className="text-red-400 text-sm text-center">{error}</div>
          )}

          <button
            type="submit"
            disabled={loading}
            className="btn-primary w-full disabled:opacity-50"
          >
            {loading ? '登录中...' : '入世'}
          </button>
        </form>

        <div className="mt-6 text-center text-gray-400 text-sm">
          初入仙途？{' '}
          <Link to="/register" className="text-primary-400 hover:text-primary-300">
            开辟道途
          </Link>
        </div>
      </div>
    </div>
  )
}
