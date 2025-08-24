import { cancelJob } from "@/services/jobs";

export const runtime = "nodejs";

export async function POST(_req: Request, { params }: { params: { id: string } }) {
  const job = cancelJob(params.id);
  if (!job) {
    return new Response(JSON.stringify({ ok: false, error: "Not found" }), {
      status: 404,
      headers: { "Content-Type": "application/json" },
    });
  }
  return Response.json({ ok: true, data: job });
}
