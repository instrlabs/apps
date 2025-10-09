import { NextResponse, NextRequest } from 'next/server'
import { ResponseCookie, ResponseCookies } from "next/dist/compiled/@edge-runtime/cookies";
import { info } from "@/utils/log";

export async function middleware(req: NextRequest) {
  console.log("MIDDLEWARE: ", req.url, req.method, req.headers.get("Content-Type"))
  const apiUrl = process.env.GATEWAY_URL;
  const next = NextResponse.next({ request: req });
  // next.headers.set("X-Testing", "testing")


  if (!req.nextUrl.pathname.startsWith("/login")) {
    const accessToken = req.cookies.get("AccessToken");
    const refreshToken = req.cookies.get("RefreshToken");

    if (!accessToken && !refreshToken) {
      info("redirect to /login", req);
      return NextResponse.redirect(new URL("/login", req.url));
    }

    if (!accessToken && refreshToken) {
      info("trying to refresh token", req);
      const headers = req.headers;
      const resRefresh = await fetch(`${apiUrl}/auth/refresh`, {
        method: "POST",
        headers: { "Content-Type": "application/json", ...headers }
      });

      if (resRefresh.ok) {
        info("successfully refreshed token", req);
        const reqSetCookie = new ResponseCookies(resRefresh.headers);
        const storeCookie = next.cookies;
        storeCookie.set(reqSetCookie.get("AccessToken") as ResponseCookie);
        storeCookie.set(reqSetCookie.get("RefreshToken") as ResponseCookie);
      } else {
        info("failed to refresh token", req);
        return NextResponse.redirect(new URL("/login", req.url));
      }
    }
  }

  return next;
}

export const config = {
  matcher: ['/((?!api|_next/static|_next/image|favicon.ico|\\.well-known|.*\\.png$).*)'],
}
