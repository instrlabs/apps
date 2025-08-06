import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';
import { ROUTES } from './constants/routes';

export function middleware(request: NextRequest) {
  const path = request.nextUrl.pathname;
  
  const isHomePage = path === ROUTES.HOME;
  
  const token = request.cookies.get('authToken')?.value;
  const isAuthenticated = !!token;
  
  if (isHomePage && !isAuthenticated) {
    return NextResponse.redirect(new URL(ROUTES.LOGIN, request.url));
  }
  
  return NextResponse.next();
}

export const config = {
  matcher: [
  '/((?!api).*)',
  '/((?!_next/static).*)',
  '/((?!_next/image).*)',
  '/((?!favicon.ico).*)',
  '/((?!login).*)',
  '/((?!register).*)',
  '/((?!forgot-password).*)',
  '/((?!reset-password).*)',
  '/((?!google/callback).*)',
  ],
};