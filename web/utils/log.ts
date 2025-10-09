import { NextRequest } from "next/server";

// Formats and prints a standardized log line for instrlabs-web
// Example: [instrlabs-web]: time="..." ip="..." host="..." path="..." message="..."
function formatLine({ time, ip, host, path, message }: { time: string; ip: string; host: string; path: string; message: string }) {
  return `[instrlabs-web]: time="${time}" ip="${ip}" host="${host}" path="${path}" message="${message}"`;
}

function extractMeta(req?: NextRequest) {
  const time = new Date().toUTCString();
  if (!req) {
    return { time, ip: "unknown", host: "unknown", path: "-" };
  }

  // IP: prefer X-Forwarded-For, else req.ip (not always available in edge runtime)
  const fwdFor = req.headers.get("x-forwarded-for") || req.headers.get("X-Forwarded-For") || "";
  const ip = fwdFor.split(",")[0].trim() || (req as any).ip || "unknown";

  // Host: prefer X-Forwarded-Host, fallback to Host
  const host =
    req.headers.get("x-forwarded-host") ||
    req.headers.get("X-Forwarded-Host") ||
    req.headers.get("host") ||
    req.headers.get("Host") ||
    "unknown";

  const path = req.nextUrl?.pathname || "-";

  return { time, ip, host, path };
}

export function webLog(message: string, req?: NextRequest) {
  const meta = extractMeta(req);
  // eslint-disable-next-line no-console
  console.log(formatLine({ ...meta, message }));
}

export function webWarn(message: string, req?: NextRequest) {
  const meta = extractMeta(req);
  // eslint-disable-next-line no-console
  console.warn(formatLine({ ...meta, message }));
}

export function webError(message: string, req?: NextRequest, error?: unknown) {
  const meta = extractMeta(req);
  // eslint-disable-next-line no-console
  if (error) {
    console.error(formatLine({ ...meta, message }) + (" error=" + JSON.stringify(serializeError(error))));
  } else {
    console.error(formatLine({ ...meta, message }));
  }
}

function serializeError(err: unknown) {
  if (err instanceof Error) {
    return { name: err.name, message: err.message, stack: err.stack };
  }
  return err;
}
