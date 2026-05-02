import { useState } from 'react'
import { useGameStore } from '@/stores/gameStore'
import { wsService } from '@/services/websocket'

export default function Sidebar() {
  const { chatMessages, isConnected } = useGameStore()
  const [input, setInput] = useState('')
  const [channel, setChannel] = useState('world')

  const handleSend = () => {
    if (input.trim() && isConnected) {
      wsService.send('chat', {
        content: input.trim(),
        channel,
      })
      setInput('')
    }
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleSend()
    }
  }

  const formatTime = (timestamp: number) => {
    return new Date(timestamp / 1000000).toLocaleTimeString('zh-CN', {
      hour: '2-digit',
      minute: '2-digit',
    })
  }

  return (
    <div className="w-80 bg-gray-850 border-l border-gray-700 flex flex-col">
      <div className="p-2 border-b border-gray-700">
        <div className="flex gap-2">
          <button
            onClick={() => setChannel('world')}
            className={`px-3 py-1 rounded text-sm ${
              channel === 'world'
                ? 'bg-primary-600 text-white'
                : 'bg-gray-700 text-gray-300'
            }`}
          >
            世界
          </button>
          <button
            onClick={() => setChannel('region')}
            className={`px-3 py-1 rounded text-sm ${
              channel === 'region'
                ? 'bg-primary-600 text-white'
                : 'bg-gray-700 text-gray-300'
            }`}
          >
            区域
          </button>
          <button
            onClick={() => setChannel('sect')}
            className={`px-3 py-1 rounded text-sm ${
              channel === 'sect'
                ? 'bg-primary-600 text-white'
                : 'bg-gray-700 text-gray-300'
            }`}
          >
            宗门
          </button>
        </div>
      </div>

      <div className="flex-1 overflow-y-auto p-2 space-y-2">
        {chatMessages.map((msg) => (
          <div
            key={msg.id}
            className="text-sm"
          >
            <span className="text-gray-500 text-xs">
              {formatTime(msg.timestamp)}
            </span>
            <span className="text-primary-400 ml-2">
              {msg.sender_name}:
            </span>
            <span className="text-gray-300 ml-1">{msg.content}</span>
          </div>
        ))}
        
        {chatMessages.length === 0 && (
          <div className="text-gray-500 text-sm text-center py-4">
            暂无消息
          </div>
        )}
      </div>

      <div className="p-2 border-t border-gray-700">
        <div className="flex gap-2">
          <input
            type="text"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            onKeyDown={handleKeyDown}
            placeholder={isConnected ? '输入消息...' : '未连接'}
            disabled={!isConnected}
            className="input-field flex-1"
          />
          <button
            onClick={handleSend}
            disabled={!isConnected || !input.trim()}
            className="btn-primary px-4 disabled:opacity-50"
          >
            发送
          </button>
        </div>
      </div>
    </div>
  )
}
