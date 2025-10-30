'use client'

import { useState, useEffect } from 'react'
import { setServerCookie, deleteServerCookie } from './actions'

interface CookieInfo {
  name: string
  value: string
  path?: string
  domain?: string
  expires?: string
  maxAge?: number
  secure?: boolean
  httpOnly?: boolean
  sameSite?: 'lax' | 'strict' | 'none'
}

export default function CookieActions() {
  const [serverCookies, setServerCookies] = useState<CookieInfo[]>([])
  const [clientCookies, setClientCookies] = useState<CookieInfo[]>([])
  const [operationResult, setOperationResult] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  
  // Helper to parse document.cookie
  const parseDocumentCookie = (cookieString: string): CookieInfo[] => {
    const cookies = cookieString.split(';').filter(Boolean)
    return cookies.map(cookie => {
      const [name, ...parts] = cookie.trim().split('=')
      const value = parts.join('=')

      return {
        name: decodeURIComponent(name),
        value: decodeURIComponent(value),
        path: '/',
        domain: '',
        secure: true,
        httpOnly: false,
        sameSite: 'lax' as const
      }
    }).filter(cookie => cookie.name && cookie.value)
  }

  // Read cookies on client side
  useEffect(() => {
    try {
      const documentCookies = parseDocumentCookie(document.cookie)
      setClientCookies(documentCookies)
    } catch (error) {
      console.error('Error reading client cookies:', error)
    }

    const interval = setInterval(() => {
      try {
        const documentCookies = parseDocumentCookie(document.cookie)
        setClientCookies(documentCookies)
      } catch (error) {
        console.error('Error reading client cookies:', error)
      }
    }, 1000)

    return () => clearInterval(interval)
  }, [])

  // Note: Server cookies cannot be read directly from client component
  // useEffect for server cookies would need to be in a server component

  const handleSetCookie = async (type: 'server' | 'client', options: {
    name: string
    value: string
    httpOnly?: boolean
    secure?: boolean
    sameSite?: 'lax' | 'strict' | 'none'
    maxAge?: number
  }) => {
    setLoading(true)
    try {
      if (type === 'server') {
        const result = await setServerCookie(options)
        setOperationResult(result.message)
        // Note: Cannot refresh server cookies from client component
        // serverCookies state will show "No server cookies found"
      } else {
        // Client-side cookie
        const cookieOptions = [
          `${encodeURIComponent(options.name)}=${encodeURIComponent(options.value)}`
        ]

        if (options.maxAge) {
          cookieOptions.push(`Max-Age=${options.maxAge}`)
        }
        cookieOptions.push('Path=/')
        if (options.secure) {
          cookieOptions.push('Secure')
        }
        if (options.sameSite) {
          cookieOptions.push(`SameSite=${options.sameSite}`)
        }

        document.cookie = cookieOptions.join('; ')
        setOperationResult(`Set client cookie "${options.name}"`)
      }
    } catch (error) {
      setOperationResult(`Error: ${error instanceof Error ? error.message : 'Unknown error'}`)
    } finally {
      setLoading(false)
      setTimeout(() => setOperationResult(null), 3000)
    }
  }

  const handleDeleteCookie = async (type: 'server' | 'client', name: string) => {
    setLoading(true)
    try {
      if (type === 'server') {
        await deleteServerCookie(name)
        setOperationResult(`Deleted server cookie "${name}"`)
      } else {
        document.cookie = `${encodeURIComponent(name)}=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/`
        setOperationResult(`Deleted client cookie "${name}"`)
      }
    } catch (error) {
      setOperationResult(`Error: ${error instanceof Error ? error.message : 'Unknown error'}`)
    } finally {
      setLoading(false)
      setTimeout(() => setOperationResult(null), 3000)
    }
  }

  const clearAllCookies = async () => {
    setLoading(true)
    try {
      // Delete server cookies
      serverCookies.forEach(cookie => {
        deleteServerCookie(cookie.name)
      })

      // Delete client cookies
      clientCookies.forEach(cookie => {
        document.cookie = `${encodeURIComponent(cookie.name)}=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/`
      })

      setOperationResult('Cleared all cookies')
      setServerCookies([])
      setClientCookies([])
    } catch (error) {
      setOperationResult(`Error: ${error instanceof Error ? error.message : 'Unknown error'}`)
    } finally {
      setLoading(false)
      setTimeout(() => setOperationResult(null), 3000)
    }
  }

  const formatCookieValue = (value: string, isSensitive = false) => {
    if (isSensitive) {
      return '*'.repeat(Math.min(value.length, 8))
    }
    return value
  }

  return (
    <div className="space-y-6">
      {/* Operation Results */}
      {operationResult && (
        <div className={`bg-green-900/20 border border-green-500/30 rounded-lg p-3 text-green-300`}>
          {operationResult}
        </div>
      )}

      {/* Cookie Creation Section */}
      <div className="bg-white/5 border border-white/20 rounded-lg p-4">
        <h2 className="text-xl font-semibold mb-4">Create Cookies</h2>

        <div className="space-y-4">
          <div className="grid md:grid-cols-2 gap-6">
            {/* Server Cookie Creation */}
            <div className="space-y-3">
              <h3 className="font-medium text-blue-300">Server-side Cookie</h3>
              <div className="space-y-2">
                <input
                  type="text"
                  placeholder="Cookie name"
                  className="w-full bg-gray-800 border border-gray-600 rounded px-3 py-2 text-sm"
                  id="server-name"
                />
                <input
                  type="text"
                  placeholder="Cookie value"
                  className="w-full bg-gray-800 border border-gray-600 rounded px-3 py-2 text-sm"
                  id="server-value"
                />
                <div className="grid grid-cols-2 gap-2">
                  <label className="flex items-center text-sm">
                    <input
                      type="checkbox"
                      className="mr-2"
                      id="server-http-only"
                    />
                    HttpOnly
                  </label>
                  <label className="flex items-center text-sm">
                    <input
                      type="checkbox"
                      className="mr-2"
                      id="server-secure"
                      defaultChecked
                    />
                    Secure
                  </label>
                </div>
                <select
                  className="w-full bg-gray-800 border border-gray-600 rounded px-3 py-2 text-sm"
                  id="server-same-site"
                  defaultValue="lax"
                >
                  <option value="lax">SameSite: Lax</option>
                  <option value="strict">SameSite: Strict</option>
                  <option value="none">SameSite: None</option>
                </select>
                <div className="flex gap-2">
                  <button
                    onClick={() => {
                      const name = (document.getElementById('server-name') as HTMLInputElement)?.value || 'server-cookie'
                      const value = (document.getElementById('server-value') as HTMLInputElement)?.value || 'server-value'
                      const httpOnly = (document.getElementById('server-http-only') as HTMLInputElement)?.checked || false
                      const secure = (document.getElementById('server-secure') as HTMLInputElement)?.checked || false
                      const sameSite = (document.getElementById('server-same-site') as HTMLSelectElement)?.value as 'lax' | 'strict' | 'none' || 'lax'
                      const maxAge = 3600 // 1 hour

                      handleSetCookie('server', { name, value, httpOnly, secure, sameSite, maxAge })
                    }}
                    disabled={loading}
                    className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded text-sm flex-1"
                  >
                    {loading ? 'Creating...' : 'Create Server Cookie'}
                  </button>
                </div>
              </div>
            </div>

            {/* Client Cookie Creation */}
            <div className="space-y-3">
              <h3 className="font-medium text-green-300">Client-side Cookie</h3>
              <div className="space-y-2">
                <input
                  type="text"
                  placeholder="Cookie name"
                  className="w-full bg-gray-800 border border-gray-600 rounded px-3 py-2 text-sm"
                  id="client-name"
                />
                <input
                  type="text"
                  placeholder="Cookie value"
                  className="w-full bg-gray-800 border border-gray-600 rounded px-3 py-2 text-sm"
                  id="client-value"
                />
                <div className="grid grid-cols-2 gap-2">
                  <label className="flex items-center text-sm">
                    <input
                      type="checkbox"
                      className="mr-2"
                      id="client-http-only"
                    />
                    HttpOnly
                  </label>
                  <label className="flex items-center text-sm">
                    <input
                      type="checkbox"
                      className="mr-2"
                      id="client-secure"
                      defaultChecked
                    />
                    Secure
                  </label>
                </div>
                <select
                  className="w-full bg-gray-800 border border-gray-600 rounded px-3 py-2 text-sm"
                  id="client-same-site"
                  defaultValue="lax"
                >
                  <option value="lax">SameSite: Lax</option>
                  <option value="strict">SameSite: Strict</option>
                  <option value="none">SameSite: None</option>
                </select>
                <div className="flex gap-2">
                  <button
                    onClick={() => {
                      const name = (document.getElementById('client-name') as HTMLInputElement)?.value || 'client-cookie'
                      const value = (document.getElementById('client-value') as HTMLInputElement)?.value || 'client-value'
                      const httpOnly = (document.getElementById('client-http-only') as HTMLInputElement)?.checked || false
                      const secure = (document.getElementById('client-secure') as HTMLInputElement)?.checked || false
                      const sameSite = (document.getElementById('client-same-site') as HTMLSelectElement)?.value as 'lax' | 'strict' | 'none' || 'lax'
                      const maxAge = 3600 // 1 hour

                      handleSetCookie('client', { name, value, httpOnly, secure, sameSite, maxAge })
                    }}
                    disabled={loading}
                    className="bg-green-600 hover:bg-green-700 text-white px-4 py-2 rounded text-sm flex-1"
                  >
                    {loading ? 'Creating...' : 'Create Client Cookie'}
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Current Cookies Display */}
      <div className="grid md:grid-cols-2 gap-6">
        {/* Server Cookies */}
        <div className="bg-white/5 border border-white/20 rounded-lg p-4">
          <div className="flex justify-between items-center mb-4">
            <h3 className="font-medium text-blue-300">Server-side Cookies</h3>
            <button
              onClick={() => clearAllCookies()}
              disabled={loading || serverCookies.length === 0}
              className="bg-red-600 hover:bg-red-700 text-white px-3 py-1 rounded text-sm"
            >
              Clear All
            </button>
          </div>

          {serverCookies.length === 0 ? (
            <p className="text-gray-400 text-sm">No server cookies found</p>
          ) : (
            <div className="space-y-2 max-h-64 overflow-y-auto">
              {serverCookies.map((cookie, index) => (
                <div key={index} className="bg-gray-800/50 rounded p-3 text-sm">
                  <div className="flex justify-between items-start mb-2">
                    <span className="font-medium text-white">{cookie.name}</span>
                    <div className="flex gap-1">
                      <button
                        onClick={() => handleDeleteCookie('server', cookie.name)}
                        disabled={loading}
                        className="bg-red-600 hover:bg-red-700 text-white px-2 py-1 rounded text-xs"
                      >
                        Delete
                      </button>
                    </div>
                  </div>
                  <div className="text-gray-300 mb-2">
                    Value: <span className="text-green-400">{formatCookieValue(cookie.value, cookie.httpOnly)}</span>
                  </div>
                  <div className="text-xs text-gray-400 space-y-1">
                    <div>Path: {cookie.path}</div>
                    {cookie.maxAge && <div>Max-Age: {cookie.maxAge}s</div>}
                    <div className="flex gap-2">
                      {cookie.httpOnly && <span className="bg-purple-600 text-white px-1 rounded">HttpOnly</span>}
                      {cookie.secure && <span className="bg-yellow-600 text-white px-1 rounded">Secure</span>}
                      {cookie.sameSite && <span className="bg-blue-600 text-white px-1 rounded">SameSite: {cookie.sameSite}</span>}
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Client Cookies */}
        <div className="bg-white/5 border border-white/20 rounded-lg p-4">
          <h3 className="font-medium text-green-300 mb-4">Client-side Cookies</h3>

          {clientCookies.length === 0 ? (
            <p className="text-gray-400 text-sm">No client cookies found</p>
          ) : (
            <div className="space-y-2 max-h-64 overflow-y-auto">
              {clientCookies.map((cookie, index) => (
                <div key={index} className="bg-gray-800/50 rounded p-3 text-sm">
                  <div className="flex justify-between items-start mb-2">
                    <span className="font-medium text-white">{cookie.name}</span>
                    <div className="flex gap-1">
                      <button
                        onClick={() => handleDeleteCookie('client', cookie.name)}
                        disabled={loading}
                        className="bg-red-600 hover:bg-red-700 text-white px-2 py-1 rounded text-xs"
                      >
                        Delete
                      </button>
                    </div>
                  </div>
                  <div className="text-gray-300 mb-2">
                    Value: <span className="text-green-400">{cookie.value}</span>
                  </div>
                  <div className="text-xs text-gray-400 space-y-1">
                    <div>Accessible via JavaScript: Yes</div>
                    <div>Secure: {cookie.secure ? 'Yes' : 'No'}</div>
                    <div>SameSite: {cookie.sameSite}</div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>

      {/* Cookie Security Tips */}
      <div className="bg-orange-900/20 border border-orange-500/30 rounded-lg p-4">
        <h3 className="text-lg font-semibold mb-3 text-orange-300">Security Best Practices</h3>
        <div className="grid md:grid-cols-2 gap-4 text-sm text-gray-300">
          <div>
            <h4 className="font-medium mb-2">Server Cookies</h4>
            <ul className="space-y-1">
              <li>• Always use Secure flag for production</li>
              <li>• Use HttpOnly for session tokens</li>
              <li>• Set appropriate SameSite values</li>
              <li>• Set reasonable expiration times</li>
            </ul>
          </div>
          <div>
            <h4 className="font-medium mb-2">Client Cookies</h4>
            <ul className="space-y-1">
              <li>• Accessible via JavaScript (no HttpOnly)</li>
              <li>• Subject to XSS vulnerabilities</li>
              <li>• Limited to 4KB per cookie</li>
              <li>• Browser storage limitations apply</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  )
}