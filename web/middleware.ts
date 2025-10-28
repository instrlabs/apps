import { NextResponse, NextRequest } from 'next/server'
import { ResponseCookie, ResponseCookies } from "next/dist/compiled/@edge-runtime/cookies";
import { error, info } from "@/utils/log";

const whitelistPaths = ['/login', '/register', '/forgot-password', '/reset-password', '/'];

export async function middleware(req: NextRequest) {
  const startTime = Date.now();
  const apiUrl = process.env.GATEWAY_URL;

  // Set up headers for gateway
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

  info(`Middleware processing request: ${req.method} ${req.nextUrl.pathname}`, req);
  info(`Auth state - Access token: ${accessToken ? "present" : "missing"}, Refresh token: ${refreshToken ? "present" : "missing"}`, req);

  if (!accessToken && refreshToken) {
    info("Attempting automatic token refresh", req);

    try {
      headers.set("content-type", "application/json");
      const refreshStartTime = Date.now();

      const res = await fetch(`${apiUrl}/auth/refresh`, {
        method: "POST",
        headers: headers,
      });

      const refreshDuration = Date.now() - refreshStartTime;
      info(`Token refresh request completed in ${refreshDuration}ms`, req);

      if (res.ok) {
        info("Successfully refreshed access token", req);
        const resSetCookie = new ResponseCookies(res.headers);
        const storeCookie = next.cookies;

        const newAccessToken = resSetCookie.get("access_token") as ResponseCookie;
        const newRefreshToken = resSetCookie.get("refresh_token") as ResponseCookie;

        if (newAccessToken && newRefreshToken) {
          storeCookie.set(newAccessToken);
          storeCookie.set(newRefreshToken);
          accessToken = newAccessToken;
          refreshToken = newRefreshToken;
          info("New tokens set successfully", req);
        } else {
          error("Invalid token response from server", req, new Error("Missing tokens in response"));
          return redirectToLogin(req);
        }
      } else {
        const errorText = await res.text();
        error(`Token refresh failed with status ${res.status}`, req, new Error(errorText));
        return redirectToLogin(req);
      }
    } catch (err) {
      error("Network error during token refresh", req, err);
      return redirectToLogin(req);
    }
  }

  // Helper function to redirect to login and clear cookies
  function redirectToLogin(request: NextRequest) {
    const nextReset = NextResponse.redirect(new URL("/login", request.url));
    nextReset.cookies.delete("access_token");
    nextReset.cookies.delete("refresh_token");
    info("Redirected to login due to authentication failure", request);
    return nextReset;
  }

  if (!whitelistPaths.includes(req.nextUrl.pathname) && !accessToken) {
    info("Access denied - no valid token and path requires authentication", req);
    return NextResponse.redirect(new URL("/login", req.url));
  }

  const totalDuration = Date.now() - startTime;
  info(`Middleware completed in ${totalDuration}ms`, req);

  return next;
}

export const config = {
  matcher: ['/((?!api|_next/static|_next/image|favicon.ico|\\.well-known|.*\\.png$).*)'],
}
