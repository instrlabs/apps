import { NextResponse, NextRequest } from 'next/server'
import { ResponseCookie, ResponseCookies } from "next/dist/compiled/@edge-runtime/cookies";

export async function middleware(request: NextRequest) {
  const next = NextResponse.next();

  if (!request.nextUrl.pathname.startsWith("/login")) {
    const accessToken = request.cookies.get("AccessToken");
    const refreshToken = request.cookies.get("RefreshToken");

    if (!accessToken && !refreshToken) {
      return NextResponse.redirect(new URL("/login", request.url));
    }

    if (!accessToken && refreshToken) {
      const baseUrl = process.env.GATEWAY_URL;
      const resRefresh = await fetch(`${baseUrl}/auth/refresh`, {
        method: "POST",
      });

      if (resRefresh.ok) {
        console.log("refresh token");
        const reqSetCookie = new ResponseCookies(resRefresh.headers);
        const storeCookie = next.cookies;
        storeCookie.set(reqSetCookie.get("AccessToken") as ResponseCookie);
        storeCookie.set(reqSetCookie.get("RefreshToken") as ResponseCookie);
      } else {
        console.log("refresh token failed", resRefresh.status);
        return NextResponse.redirect(new URL("/login", request.url));
      }
    }
  }

  return next;
}

export const config = {
  matcher: ['/((?!api|_next/static|_next/image|favicon.ico|\\.well-known|.*\\.png$).*)'],
}
