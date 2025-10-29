import { NextResponse, NextRequest } from "next/server";
import { ResponseCookie, ResponseCookies } from "next/dist/compiled/@edge-runtime/cookies";
import { error, info } from "@/utils/log";

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

  const accessToken = req.cookies.get("access_token");
  const refreshToken = req.cookies.get("refresh_token");

  if (!accessToken && refreshToken) {
    info("Attempting automatic token refresh", req);

    try {
      headers.set("content-type", "application/json");
      const res = await fetch(`${apiUrl}/auth/refresh`, { method: "POST", headers: headers });

      if (res.ok) {
        info("Successfully refreshed access token", req);
        const resSetCookie = new ResponseCookies(res.headers);
        const storeCookie = next.cookies;

        const newAccessToken = resSetCookie.get("access_token") as ResponseCookie;
        const newRefreshToken = resSetCookie.get("refresh_token") as ResponseCookie;

        if (newAccessToken && newRefreshToken) {
          storeCookie.set(newAccessToken);
          storeCookie.set(newRefreshToken);
        } else {
          error("Failed to retrieve new access and refresh tokens", req);
          return redirectToLogin(req);
        }
      } else {
        const errorText = await res.text();
        error("Token refresh failed with status " + res.status, req, new Error(errorText));
        return redirectToLogin(req);
      }
    } catch (err) {
      error("Network error during token refresh", req, err);
      return redirectToLogin(req);
    }
  }

  function redirectToLogin(request: NextRequest) {
    const redirectResponse = NextResponse.redirect(new URL("/login", request.url));
    redirectResponse.cookies.delete("access_token");
    redirectResponse.cookies.delete("refresh_token");
    info("Redirected to login due to authentication failure", request);
    return redirectResponse;
  }

  return next;
}

export const config = {
  matcher: ["/((?!api|_next/static|_next/image|favicon.ico|\\.well-known|.*\\.png$).*)"],
};
