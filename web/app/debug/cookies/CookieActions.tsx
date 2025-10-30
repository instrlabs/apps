'use client'

import { useState } from 'react'
import { manageCookie } from './actions'

export default function CookieActions() {
  const [operationResult, setOperationResult] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  const handleCreateCookie = async () => {
    setLoading(true)
    try {
      const result = await manageCookie({
        key: 'test-cookies',
        value: 'test-value',
        action: 'set',
        server: true,
        client: true,
        sameSite: 'lax',
        maxAge: 3600
      })

      // Set client-side cookie
      const cookieOptions = [
        `${encodeURIComponent('test-cookies')}=${encodeURIComponent('test-value')}`
      ]
      cookieOptions.push('Path=/')
      cookieOptions.push('Secure')
      cookieOptions.push('SameSite=lax')
      document.cookie = cookieOptions.join('; ')

      setOperationResult(result.message)
    } catch (error) {
      setOperationResult(`Error: ${error instanceof Error ? error.message : 'Unknown error'}`)
    } finally {
      setLoading(false)
      setTimeout(() => setOperationResult(null), 3000)
    }
  }

  const handleDeleteCookie = async () => {
    setLoading(true)
    try {
      // Delete client-side cookie
      document.cookie = `${encodeURIComponent('test-cookies')}=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/`

      const result = await manageCookie({
        key: 'test-cookies',
        action: 'delete',
        server: true,
        client: true
      })

      setOperationResult(result.message)
    } catch (error) {
      setOperationResult(`Error: ${error instanceof Error ? error.message : 'Unknown error'}`)
    } finally {
      setLoading(false)
      setTimeout(() => setOperationResult(null), 3000)
    }
  }

  return (
    <div className="space-y-4">
      {operationResult && (
        <div className="rounded-lg border border-green-500/30 bg-green-900/20 p-3 text-green-300">
          {operationResult}
        </div>
      )}

      <div className="flex gap-3">
        <button
          onClick={handleCreateCookie}
          disabled={loading}
          className="flex-1 rounded bg-blue-600 px-6 py-3 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-60"
        >
          {loading ? 'Creating...' : 'Create test-cookies'}
        </button>
        <button
          onClick={handleDeleteCookie}
          disabled={loading}
          className="flex-1 rounded bg-red-600 px-6 py-3 text-sm font-medium text-white hover:bg-red-700 disabled:opacity-60"
        >
          {loading ? 'Deleting...' : 'Delete test-cookies'}
        </button>
      </div>
    </div>
  )
}
