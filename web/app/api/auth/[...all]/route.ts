import {NextRequest, NextResponse} from "next/server";

export const dynamic = "force-dynamic";

function buildTarget(req: NextRequest): string {
  const baseUrl = process.env.API_URL as string;
  const parts = req.nextUrl.pathname.split("/").slice(3);
  const subPath = parts.join("/");
  const qs = req.nextUrl.search || "";
  return `${baseUrl}/auth/${subPath}${qs}`;
}

async function forward(req: NextRequest) {
  const target = buildTarget(req);

  let body: string | undefined;
  if (req.method !== "GET" && req.method !== "HEAD") {
    try {
      const json = await req.json();
      body = JSON.stringify(json);
    } catch {
      body = undefined;
    }
  }

  const headers: HeadersInit = {};
  const cookie = req.headers.get("cookie");
  if (cookie) headers["Cookie"] = cookie;
  if (body) headers["Content-Type"] = "application/json";

  const res = await fetch(target, {
    method: req.method,
    credentials: "include",
    headers,
    body,
  });

  const resJson = await res.json().catch(() => ({}));
  const nextRes = NextResponse.json(resJson, { status: res.status });

  const setCookie = res.headers.get("Set-Cookie");
  if (setCookie) {
    nextRes.headers.set("Set-Cookie", setCookie);
  }

  return nextRes;
}

export async function GET(req: NextRequest) {
  return forward(req);
}

export async function POST(req: NextRequest) {
  return forward(req);
}

export async function PUT(req: NextRequest) {
  return forward(req);
}

export async function PATCH(req: NextRequest) {
  return forward(req);
}

export async function DELETE(req: NextRequest) {
  return forward(req);
}

export async function OPTIONS(req: NextRequest) {
  return forward(req);
}

export async function HEAD(req: NextRequest) {
  return forward(req);
}
