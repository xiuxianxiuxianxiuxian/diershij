import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { register as registerApi } from '@/services/api'

export default function RegisterPage() {
  const navigate = useNavigate()
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')

    if (password !== confirmPassword) {
      setError('两次输入的密令不一致')
      return
    }

    if (password.length < 6) {
      setError('密令至少需要6个字符')
      return
    }

    setLoading(true)

    try {
      const result = await registerApi(username, password)
      
      if (result.success) {
        navigate('/login')
      } else {
        setError(result.message || '注册失败')
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
          <h1 className="text-3xl font-bold text-immortal-gold mb-2">开辟道途</h1>
          <p className="text-gray-400">踏入修仙之路，开启你的传奇</p>
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
              minLength={2}
              maxLength={20}
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
              minLength={6}
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-300 mb-1">
              确认密令
            </label>
            <input
              type="password"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              className="input-field"
              placeholder="请再次输入密令"
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
            {loading ? '开辟中...' : '开辟道途'}
          </button>
        </form>

        <div className="mt-6 text-center text-gray-400 text-sm">
          已有道途？{' '}
          <Link to="/login" className="text-primary-400 hover:text-primary-300">
            返回登录
          </Link>
        </div>
      </div>
    </div>
  )
}
