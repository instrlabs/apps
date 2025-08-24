"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";
import { MOCK_MODE } from "@/constants/mock";

export default function ImageCompressPage() {
  const router = useRouter();
  const [file, setFile] = useState<File | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    if (!file) return;

    setSubmitting(true);
    try {
      if (MOCK_MODE) {
        // Skip API call; just fake a job id and redirect
        const jobId = `mock-${Math.random().toString(36).slice(2)}`;
        router.push(`/histories?jobId=${encodeURIComponent(jobId)}`);
        return;
      }

      const form = new FormData();
      form.set("type", "image-compress");
      form.set("file", file);

      const res = await fetch("/api/jobs", {
        method: "POST",
        body: form,
      });
      const json = await res.json();
      if (!res.ok || !json?.ok) {
        throw new Error(json?.error || "Failed to start job");
      }
      const jobId = json.data.id as string;
      // Redirect to histories page showing the job in progress
      router.push(`/histories?jobId=${encodeURIComponent(jobId)}`);
    } catch (err) {
      const message = err instanceof Error ? err.message : "Something went wrong";
      setError(message);
      setSubmitting(false);
    }
  }

  return (
    <div className="container mx-auto p-4">
      <div className="mb-6">
        <h1 className="text-2xl font-semibold">Image Compress</h1>
        <p className="text-sm text-muted">Upload an image and start compression. You&apos;ll be redirected to Histories to track progress.</p>
      </div>

      <form onSubmit={onSubmit} className="max-w-xl space-y-4">
        <div>
          <label className="block text-sm font-medium mb-1">Select image</label>
          <input
            type="file"
            accept="image/*"
            onChange={(e) => setFile(e.target.files?.[0] || null)}
            className="block w-full"
          />
        </div>

        {error && (
          <div className="text-red-600 text-sm">{error}</div>
        )}

        <button
          type="submit"
          disabled={!file || submitting}
          className="inline-flex items-center px-4 py-2 rounded-md bg-primary text-white disabled:opacity-50"
        >
          {submitting ? "Starting..." : "Start compress"}
        </button>
      </form>
    </div>
  );
}
