import { NextResponse, NextRequest } from "next/server";

export async function proxy(req: NextRequest) {
  const forwardedFor = req.headers.get("x-forwarded-for")!;
  const forwardedProto = req.headers.get("x-forwarded-proto")!;
  const forwardedHost = req.headers.get("x-forwarded-host")!;
  const userAgent = req.headers.get("user-agent")!;
  const cookie = req.headers.get("cookie")!;

  const headers = new Headers();
  headers.set("x-user-ip", forwardedFor);
  headers.set("x-user-agent", userAgent);
  headers.set("x-user-host", forwardedHost);
  headers.set("x-user-origin", forwardedProto + "://" + forwardedHost);
  headers.set("cookie", cookie);
  const next = NextResponse.next({ headers });

  const accessToken = req.cookies.get("access_token");
  const refreshToken = req.cookies.get("refresh_token");

  // TODO: check that page is whitelisted or not

  if (!accessToken && refreshToken) {
    // TODO: do refresh token by auth.ts:refresh 401 call auth.ts:logout
  }

  return next;
}

export const config = {
  matcher: ["/((?!api|_next/static|_next/image|favicon.ico|\\.well-known|.*\\.png$).*)"],
};
