import { NextRequest } from "next/server";
import { getJob } from "@/services/jobs";

export const runtime = "nodejs";

export async function GET(_req: NextRequest, { params }: { params: { id: string } }) {
  const job = getJob(params.id);
  if (!job) {
    return new Response(JSON.stringify({ ok: false, error: "Not found" }), {
      status: 404,
      headers: { "Content-Type": "application/json" },
    });
  }
  return Response.json({ ok: true, data: job });
}
