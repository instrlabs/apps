import { NextResponse, NextRequest } from 'next/server'
import { ResponseCookie, ResponseCookies } from "next/dist/compiled/@edge-runtime/cookies";
import { info } from "@/utils/log";
import { revalidatePath, revalidateTag } from "next/cache";

export async function middleware(req: NextRequest) {
  const apiUrl = process.env.GATEWAY_URL;
  const next = NextResponse.next({ request: req });

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
      headers.set("Content-Type", "application/json");
      headers.set("X-User-Agent", req.headers.get("user-agent")!);
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
    }
  }

  return next;
}

export const config = {
  matcher: ['/((?!api|_next/static|_next/image|favicon.ico|\\.well-known|.*\\.png$).*)'],
}
