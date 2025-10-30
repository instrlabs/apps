import { NextRequest } from "next/server";

export const runtime = "nodejs";

export async function GET(req: NextRequest) {
  const cookie = req.headers.get("cookie");

  if (!cookie?.includes("access_token=")) {
    return new Response(JSON.stringify({ error: "Unauthorized" }), {
      status: 401,
      headers: { "Content-Type": "application/json" },
    });
  }

  const notificationUrl = process.env.NOTIFICATION_URL;
  if (!notificationUrl) {
    return new Response(JSON.stringify({ error: "Service unavailable" }), {
      status: 500,
      headers: { "Content-Type": "application/json" },
    });
  }

  const targetUrl = notificationUrl.replace(/\/$/, "") + "/sse";
  const origin = req.headers.get("origin") ?? new URL(req.url).origin;

  let upstream: Response;
  try {
    upstream = await fetch(targetUrl, {
      method: "GET",
      headers: {
        Accept: "text/event-stream",
        Cookie: cookie,
        Origin: origin,
        "Accept-Encoding": "identity",
      },
      signal: req.signal,
      redirect: "follow",
      cache: "no-store",
    });
  } catch (error) {
    return new Response(JSON.stringify({ error: "Connection failed" }), {
      status: 502,
      headers: { "Content-Type": "application/json" },
    });
  }

  if (!upstream.ok || !upstream.body) {
    return new Response(JSON.stringify({ error: "Upstream error" }), {
      status: upstream.status || 502,
      headers: { "Content-Type": "application/json" },
    });
  }

  const headers = new Headers();
  headers.set(
    "Content-Type",
    upstream.headers.get("content-type") ?? "text/event-stream; charset=utf-8"
  );
  headers.set("Cache-Control", "no-cache, no-transform");
  headers.set("Connection", "keep-alive");

  const setCookie = upstream.headers.get("set-cookie");
  if (setCookie) headers.append("set-cookie", setCookie);

  return new Response(upstream.body, { status: 200, headers });
}
