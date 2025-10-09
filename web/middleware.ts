import { NextResponse, NextRequest } from 'next/server'
import { ResponseCookie, ResponseCookies } from "next/dist/compiled/@edge-runtime/cookies";

export async function middleware(req: NextRequest) {
  const apiUrl = process.env.GATEWAY_URL;
  const next = NextResponse.next({ request: req });

  const time = new Date().toUTCString();
  const ip = req.headers.get("X-Forwarded-For");
  const path = req.nextUrl.pathname;

  if (!req.nextUrl.pathname.startsWith("/login")) {
    const accessToken = req.cookies.get("AccessToken");
    const refreshToken = req.cookies.get("RefreshToken");

    if (!accessToken && !refreshToken) {
      console.log(`[instrlabs-web]: time="${time}" ip="${ip}" path="${path}" message="redirect to /login"`);
      return NextResponse.redirect(new URL("/login", req.url));
    }

    if (!accessToken && refreshToken) {
      console.log(`[instrlabs-web]: time="${time}" ip="${ip}" path="${path}" message="trying to refresh token"`);
      const headers = req.headers;
      const resRefresh = await fetch(`${apiUrl}/auth/refresh`, {
        method: "POST",
        headers: { "Content-Type": "application/json", ...headers }
      });

      if (resRefresh.ok) {
        console.log(`[instrlabs-web]: time="${time}" ip="${ip}" path="${path}" message="trying to refresh token"`);
        const reqSetCookie = new ResponseCookies(resRefresh.headers);
        const storeCookie = next.cookies;
        storeCookie.set(reqSetCookie.get("AccessToken") as ResponseCookie);
        storeCookie.set(reqSetCookie.get("RefreshToken") as ResponseCookie);
      } else {
        console.log(`[instrlabs-web]: time="${time}" ip="${ip}" path="${path}" message="failed to refresh token"`);
        return NextResponse.redirect(new URL("/login", req.url));
      }
    }
  }

  return next;
}

export const config = {
  matcher: ['/((?!api|_next/static|_next/image|favicon.ico|\\.well-known|.*\\.png$).*)'],
}
