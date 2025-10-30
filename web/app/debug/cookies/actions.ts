"use server";

import { cookies } from 'next/headers';

export const setServerCookie = async (options: {
  name: string;
  value: string;
  httpOnly?: boolean;
  secure?: boolean;
  sameSite?: 'lax' | 'strict' | 'none';
  maxAge?: number;
}) => {
  const cookieStore = await cookies()
  const { name, value, httpOnly = false, secure = true, sameSite = 'lax', maxAge } = options
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