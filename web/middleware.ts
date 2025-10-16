import { NextResponse, NextRequest } from 'next/server'
import { ResponseCookie, ResponseCookies } from "next/dist/compiled/@edge-runtime/cookies";
import { error, info } from "@/utils/log";

const whitelistPaths = ['/login', '/register', '/forgot-password', '/reset-password', '/'];

export async function middleware(req: NextRequest) {
  const apiUrl = process.env.GATEWAY_URL;
  const headers = new Headers();
  headers.set("x-user-ip", req.headers.get("x-forwarded-for")!);
  headers.set("x-user-agent", req.headers.get("user-agent")!);
  const forwardedHost = req.headers.get("x-forwarded-host")!;
  headers.set("x-user-host", forwardedHost);
  const forwardedProto = req.headers.get("x-forwarded-proto")!;
  headers.set("x-user-origin", forwardedProto + "://" + forwardedHost);
  headers.set("cookie", req.headers.get("cookie")!);
  const next = NextResponse.next({ headers });

  let accessToken = req.cookies.get("access_token");
  let refreshToken = req.cookies.get("refresh_token");

  if (!accessToken && refreshToken) {
    info("trying to refresh token", req);

    try {
      headers.set("content-type", "application/json");
      const res = await fetch(`${apiUrl}/auth/refresh`, {
        method: "POST",
        headers: headers,
      });

      if (res.ok) {
        info("successfully refreshed token", req);
        const resSetCookie = new ResponseCookies(res.headers);
        const storeCookie = next.cookies;
        const newAccessToken = resSetCookie.get("access_token") as ResponseCookie;
        const newRefreshToken = resSetCookie.get("refresh_token") as ResponseCookie;
        storeCookie.set(newAccessToken);
        storeCookie.set(newRefreshToken);
        accessToken = newAccessToken;
        refreshToken = newRefreshToken;
      } else {
        info("failed to refresh token", req);
        return NextResponse.redirect(new URL("/login", req.url));
      }
    } catch (err) {
      error("failed to refresh token", req, err);
      return NextResponse.redirect(new URL("/login", req.url));
    }
  }

  if (!whitelistPaths.includes(req.nextUrl.pathname) && !accessToken) {
    return NextResponse.redirect(new URL("/login", req.url));
  }

  return next;
}

export const config = {
  matcher: ['/((?!api|_next/static|_next/image|favicon.ico|\\.well-known|.*\\.png$).*)'],
}
