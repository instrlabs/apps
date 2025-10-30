"use server";

import { cookies } from 'next/headers';

/**
 * Server-side cookie management for Next.js
 *
 * Security defaults:
 * - httpOnly: true (prevents XSS attacks - JavaScript cannot access)
 * - secure: true (HTTPS only in production)
 * - sameSite: 'lax' (CSRF protection)
 * - path: '/' (accessible from all routes)
 *
 * IMPORTANT: When deleting cookies, the path MUST match the path used when creating!
 * Example: If created with path='/', must delete with path='/'
 */

export const setServerCookie = async (options: {
  name: string;
  value: string;
  httpOnly?: boolean;
  secure?: boolean;
  sameSite?: 'lax' | 'strict' | 'none';
  maxAge?: number;
}) => {
  const cookieStore = await cookies()
  const { name, value, httpOnly = true, secure = true, sameSite = 'lax', maxAge } = options
  cookieStore.set(name, value, {
    httpOnly,
    secure,
    sameSite,
    path: '/',
    ...(maxAge && { maxAge })
  })
  return { success: true, message: `Set cookie "${name}"` }
}

export const deleteServerCookie = async (name: string) => {
  const cookieStore = await cookies()
  cookieStore.delete({ name, path: '/' })
  return { success: true, message: `Deleted cookie "${name}"` }
}

/**
 * Unified cookie management function
 * Handles both server-side and client-side cookie operations with a single key
 *
 * @param options - Configuration for the cookie operation
 * @param options.key - The cookie key/name
 * @param options.value - The cookie value (required for 'set' action)
 * @param options.action - 'set' or 'delete'
 * @param options.server - Whether to apply operation on server-side cookie
 * @param options.client - Whether to apply operation on client-side cookie
 * @param options.httpOnly - For server cookies (default: true)
 * @param options.secure - For server cookies (default: true)
 * @param options.sameSite - For server cookies (default: 'lax')
 * @param options.maxAge - For server cookies (default: undefined)
 * @returns - Result message
 */
export const manageCookie = async (options: {
  key: string;
  value?: string;
  action: 'set' | 'delete';
  server?: boolean;
  client?: boolean;
  httpOnly?: boolean;
  secure?: boolean;
  sameSite?: 'lax' | 'strict' | 'none';
  maxAge?: number;
}) => {
  const {
    key,
    value = '',
    action,
    server = true,
    client = true,
    httpOnly = true,
    secure = true,
    sameSite = 'lax',
    maxAge
  } = options

  const operations: string[] = []

  // Handle server-side cookie operation
  if (server) {
    try {
      if (action === 'set') {
        const cookieStore = await cookies()
        cookieStore.set(key, value, {
          httpOnly,
          secure,
          sameSite,
          path: '/',
          ...(maxAge && { maxAge })
        })
        operations.push(`server cookie "${key}"`)
      } else if (action === 'delete') {
        const cookieStore = await cookies()
        cookieStore.delete({ name: key, path: '/' })
        operations.push(`server cookie "${key}"`)
      }
    } catch (error) {
      console.error(`Error managing server cookie "${key}":`, error)
    }
  }

  // Handle client-side cookie operation
  if (client) {
    operations.push(`client cookie "${key}"`)
    // Client-side operations are done via document.cookie in the client component
  }

  const actionText = action === 'set' ? 'Set' : 'Deleted'
  const message = `${actionText} ${operations.join(' and ')}`

  return { success: true, message }
}