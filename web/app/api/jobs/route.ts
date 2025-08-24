import { NextRequest } from "next/server";
import { createJob, listJobs } from "@/services/jobs";

export const runtime = "nodejs";

export async function GET(req: NextRequest) {
  const { searchParams } = new URL(req.url);
  const type = searchParams.get("type") || undefined;
  const data = listJobs(type);
  return Response.json({ ok: true, data });
}

export async function POST(req: NextRequest) {
  // Expect multipart/form-data with fields: file, type
  try {
    const form = await req.formData();
    const file = form.get("file");
    const type = (form.get("type") || "") as string;

    if (!type) {
      return new Response(JSON.stringify({ ok: false, error: "Missing type" }), {
        status: 400,
        headers: { "Content-Type": "application/json" },
      });
    }

    if (!(file instanceof File)) {
      return new Response(JSON.stringify({ ok: false, error: "Missing file" }), {
        status: 400,
        headers: { "Content-Type": "application/json" },
      });
    }

    const filename = (file as File).name || "upload";
    const job = createJob({ type, filename });
    return Response.json({ ok: true, data: job });
  } catch (e) {
    const message = e instanceof Error ? e.message : "Failed to create job";
    return new Response(JSON.stringify({ ok: false, error: message }), {
      status: 500,
      headers: { "Content-Type": "application/json" },
    });
  }
}
