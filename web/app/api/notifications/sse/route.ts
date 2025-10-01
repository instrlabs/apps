import { NextRequest } from "next/server";

export const runtime = "nodejs"; // ensure Node runtime for streaming

export async function GET(req: NextRequest) {
  const targetBase = process.env.NOTIFICATION_URL;
  if (!targetBase) {
    return new Response("NOTIFICATION_URL is not configured", { status: 500 });
  }

  const targetUrl = targetBase.replace(/\/$/, "") + "/sse";

  // Forward cookies and origin to the upstream
  const cookie = req.headers.get("cookie") ?? "";
  const origin = req.headers.get("origin") ?? new URL(req.url).origin;

  let upstream: Response;
  try {
    upstream = await fetch(targetUrl, {
      method: "GET",
      headers: {
        // Accept helps some servers gate SSE
        Accept: "text/event-stream",
        Cookie: cookie,
        Origin: origin,
        // Avoid compression as it can interfere with SSE chunking
        "Accept-Encoding": "identity",
      },
      // Abort upstream if client disconnects
      signal: req.signal,
      redirect: "follow",
      cache: "no-store",
    });
  } catch (e) {
    return new Response(`Failed to connect upstream: ${e instanceof Error ? e.message : String(e)}` , { status: 502 });
  }

  if (!upstream.ok || !upstream.body) {
    const text = await upstream.text().catch(() => upstream.statusText);
    return new Response(text || "Upstream error", { status: upstream.status || 502 });
  }

  // Prepare SSE response headers
  const resHeaders = new Headers();
  resHeaders.set("Content-Type", upstream.headers.get("content-type") ?? "text/event-stream; charset=utf-8");
  resHeaders.set("Cache-Control", "no-cache, no-transform");
  resHeaders.set("Connection", "keep-alive");
  // Allow sending cookies back if upstream sets any
  const setCookie = upstream.headers.get("set-cookie");
  if (setCookie) resHeaders.append("set-cookie", setCookie);

  // Stream upstream body directly to the client
  return new Response(upstream.body, {
    status: 200,
    headers: resHeaders,
  });
}
