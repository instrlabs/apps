import { NextRequest } from "next/server";

function formatLine({ time, ip, host, path, message }: {
  time: string; ip: string; host: string; path: string; message: string }) {
  return `[instrlabs-web]: time="${time}" ip="${ip}" host="${host}" path="${path}" message="${message}"`;
}

function extractMeta(req: NextRequest) {
  const time = new Date().toUTCString();
  const ip = req.headers.get("x-forwarded-for")!;
  const host = req.headers.get("x-forwarded-host")!;
  const path = req.nextUrl?.pathname || "-";
  return { time, ip, host, path };
}

export function info(message: string, req: NextRequest) {
  const meta = extractMeta(req);
  console.log(formatLine({ ...meta, message }));
}

export function warn(message: string, req: NextRequest) {
  const meta = extractMeta(req);
  console.warn(formatLine({ ...meta, message }));
}

export function error(message: string, req: NextRequest, error?: unknown) {
  const meta = extractMeta(req);
  if (error) console.error(formatLine({ ...meta, message }) + (" error=" + JSON.stringify(serializeError(error))));
  else console.error(formatLine({ ...meta, message }));
}

function serializeError(err: unknown) {
  if (err instanceof Error) {
    return {
      name: err.name,
      message: err.message,
      stack: err.stack
    };
  }
  return err;
}
