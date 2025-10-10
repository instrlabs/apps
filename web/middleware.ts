import { NextResponse, NextRequest } from 'next/server'
import { ResponseCookie, ResponseCookies } from "next/dist/compiled/@edge-runtime/cookies";
import { error, info } from "@/utils/log";

export async function middleware(req: NextRequest) {
  const apiUrl = process.env.GATEWAY_URL;
  const headers = new Headers();
  headers.set("x-user-ip", req.headers.get("x-forwarded-for")!);
  headers.set("x-user-agent", req.headers.get("user-agent")!);
  const forwardedHost = req.headers.get("x-forwarded-host")!;
  const hostWithoutPort = forwardedHost.split(":")[0];
  const domainParts = hostWithoutPort.split(".");
  const mainDomain = domainParts.length > 2 ? domainParts.slice(-2).join(".") : hostWithoutPort;
  headers.set("x-user-host", forwardedHost);
  headers.set("x-user-origin", mainDomain);

  const accessToken = req.cookies.get("access_token");
  const refreshToken = req.cookies.get("refresh_token")
  headers.set("cookie", "access_token=" + accessToken + "; refresh_token=" + refreshToken + ";");

  const next = NextResponse.next({ headers });


  if (!req.nextUrl.pathname.startsWith("/login")) {
    if (!accessToken && !refreshToken) {
      info("redirect to /login", req);
      return NextResponse.redirect(new URL("/login", req.url));
    }

    if (!accessToken && refreshToken) {
      info("trying to refresh token", req);

      try {
        headers.set("content-type", "application/json");
        const resRefresh = await fetch(`${apiUrl}/auth/refresh`, {
          method: "POST",
          headers: headers,
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
      } catch (err) {
        error("failed to refresh token", req, err);
        return NextResponse.redirect(new URL("/login", req.url));
      }
    }
  }

  return next;
}

export const config = {
  matcher: ['/((?!api|_next/static|_next/image|favicon.ico|\\.well-known|.*\\.png$).*)'],
}
