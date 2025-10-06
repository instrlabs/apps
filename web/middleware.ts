import { NextResponse, NextRequest } from 'next/server'
import { ResponseCookie, ResponseCookies } from "next/dist/compiled/@edge-runtime/cookies";
import { store } from "next/dist/build/output/store";

export async function middleware(request: NextRequest) {
  const next = NextResponse.next();

  if (!request.nextUrl.pathname.startsWith("/login")) {
    const accessToken = request.cookies.get("AccessToken");
    const refreshToken = request.cookies.get("RefreshToken");

    if (!accessToken && !refreshToken) {
      return NextResponse.redirect(new URL("/login", request.url));
    }

    if (!accessToken && refreshToken) {
      console.log("RefreshToken: Trying to refresh");
      const baseUrl = process.env.GATEWAY_URL;
      const resRefresh = await fetch(`${baseUrl}/auth/refresh`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "Origin": request.headers.get("origin") ?? "",
          "Cookie": request.headers.get("cookie") ?? ""
        }
      });

      if (resRefresh.ok) {
        console.log("RefreshToken: Success");
        const reqSetCookie = new ResponseCookies(resRefresh.headers);
        const storeCookie = next.cookies;
        storeCookie.set(reqSetCookie.get("AccessToken") as ResponseCookie);
        storeCookie.set(reqSetCookie.get("RefreshToken") as ResponseCookie);
      } else {
        console.log("RefreshToken: Failed");
        return NextResponse.redirect(new URL("/login", request.url));
      }
    }
  }

  return next;
}

export const config = {
  matcher: ['/((?!api|_next/static|_next/image|favicon.ico|\\.well-known|.*\\.png$).*)'],
}
