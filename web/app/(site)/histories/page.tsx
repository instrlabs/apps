"use client";

import { useEffect, useMemo, useState } from "react";
import { useSearchParams } from "next/navigation";
import { MOCK_MODE } from "@/constants/mock";

type Job = {
  id: string;
  type: string;
  filename: string;
  status: "queued" | "in_progress" | "completed" | "canceled" | "failed";
  progress: number;
  createdAt: number;
  updatedAt: number;
  error?: string;
};

const TOOL_OPTIONS = [
  { label: "All tools", value: "" },
  { label: "Image Compress", value: "image-compress" },
];

const MOCK_JOBS: Job[] = [
  {
    id: "mock-1",
    type: "image-compress",
    filename: "holiday-photo.jpg",
    status: "in_progress",
    progress: 42,
    createdAt: Date.now() - 1000 * 60 * 3,
    updatedAt: Date.now() - 1000 * 15,
  },
  {
    id: "mock-2",
    type: "image-compress",
    filename: "screenshot.png",
    status: "queued",
    progress: 0,
    createdAt: Date.now() - 1000 * 60 * 5,
    updatedAt: Date.now() - 1000 * 45,
  },
  {
    id: "mock-3",
    type: "image-compress",
    filename: "banner.webp",
    status: "completed",
    progress: 100,
    createdAt: Date.now() - 1000 * 60 * 60 * 2,
    updatedAt: Date.now() - 1000 * 60 * 60 * 2 + 5000,
  },
  {
    id: "mock-4",
    type: "image-compress",
    filename: "bad-image.bmp",
    status: "failed",
    progress: 73,
    createdAt: Date.now() - 1000 * 60 * 60 * 5,
    updatedAt: Date.now() - 1000 * 60 * 60 * 5 + 8000,
    error: "Network error while processing",
  },
  {
    id: "mock-5",
    type: "image-compress",
    filename: "draft.png",
    status: "canceled",
    progress: 25,
    createdAt: Date.now() - 1000 * 60 * 50,
    updatedAt: Date.now() - 1000 * 60 * 50 + 6000,
  },
];

export default function HistoriesPage() {
  const searchParams = useSearchParams();
  const highlightedJobId = searchParams.get("jobId");

  const [jobs, setJobs] = useState<Job[]>([]);
  const [filter, setFilter] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  const inProgressIds = useMemo(() => new Set(jobs.filter(j => j.status === "queued" || j.status === "in_progress").map(j => j.id)), [jobs]);

  async function fetchJobs(signal?: AbortSignal) {
    if (MOCK_MODE) {
      setLoading(true);
      // Simulate latency
      setTimeout(() => {
        const data = filter ? MOCK_JOBS.filter(j => j.type === filter) : MOCK_JOBS;
        setJobs(data);
        setError(null);
        setLoading(false);
      }, 300);
      return;
    }

    setLoading(true);
    setError(null);
    try {
      const url = new URL("/api/jobs", window.location.origin);
      if (filter) url.searchParams.set("type", filter);
      const res = await fetch(url.toString(), { signal });
      const json = await res.json();
      if (!res.ok || !json?.ok) throw new Error(json?.error || "Failed to load jobs");
      setJobs(json.data as Job[]);
    } catch (e) {
      const isAbort = e instanceof DOMException && e.name === "AbortError";
      if (!isAbort) {
        const message = e instanceof Error ? e.message : "Failed to load jobs";
        setError(message);
      }
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    const controller = new AbortController();
    fetchJobs(controller.signal);
    return () => controller.abort();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [filter]);

  useEffect(() => {
    if (MOCK_MODE) return; // no polling in mock mode
    // Poll every 2s while there are jobs in progress
    if (inProgressIds.size === 0) return;
    const id = setInterval(() => {
      fetchJobs();
    }, 2000);
    return () => clearInterval(id);
  }, [inProgressIds]);

  async function cancel(id: string) {
    if (MOCK_MODE) {
      // Update local state to reflect cancellation
      setJobs(prev => prev.map(j => j.id === id ? { ...j, status: "canceled" as const } : j));
      return;
    }
    try {
      const res = await fetch(`/api/jobs/${encodeURIComponent(id)}/cancel`, { method: "POST" });
      const json = await res.json();
      if (!res.ok || !json?.ok) throw new Error(json?.error || "Failed to cancel job");
      // Refresh list
      fetchJobs();
    } catch (e) {
      const message = e instanceof Error ? e.message : "Failed to cancel job";
      alert(message);
    }
  }

  return (
    <div className="container mx-auto p-4">
      <div className="mb-6 flex items-end justify-between gap-4 flex-wrap">
        <div>
          <h1 className="text-2xl font-semibold">Histories</h1>
          <p className="text-sm text-muted">Track your processing jobs. Filter by tool type and cancel if needed.</p>
        </div>
        <div>
          <label className="block text-sm font-medium mb-1">Filter by tool</label>
          <select
            value={filter}
            onChange={(e) => setFilter(e.target.value)}
            className="border border-border rounded-md px-2 py-1 bg-background"
          >
            {TOOL_OPTIONS.map(opt => (
              <option key={opt.value} value={opt.value}>{opt.label}</option>
            ))}
          </select>
        </div>
      </div>

      {loading && jobs.length === 0 && (
        <div className="text-sm text-muted">Loading...</div>
      )}
      {error && (
        <div className="text-sm text-red-600">{error}</div>
      )}

      <div className="space-y-3">
        {jobs.map(job => {
          const canCancel = job.status === "queued" || job.status === "in_progress";
          const highlighted = job.id === highlightedJobId;
          return (
            <div
              key={job.id}
              className={`border border-border rounded-lg p-4 ${highlighted ? "ring-2 ring-primary" : ""}`}
            >
              <div className="flex items-center justify-between gap-4 flex-wrap">
                <div>
                  <div className="text-sm text-muted">{new Date(job.createdAt).toLocaleString()}</div>
                  <div className="font-medium mt-0.5">{job.filename} <span className="text-muted text-sm">({job.type})</span></div>
                </div>
                <div className="flex items-center gap-2">
                  <StatusBadge status={job.status} />
                  {canCancel && (
                    <button
                      onClick={() => cancel(job.id)}
                      className="px-3 py-1 rounded-md border border-border hover:bg-muted"
                    >
                      Cancel
                    </button>
                  )}
                </div>
              </div>

              <div className="mt-3">
                <div className="h-2 w-full bg-muted rounded">
                  <div
                    className={`h-2 rounded ${job.status === "canceled" ? "bg-gray-400" : job.status === "failed" ? "bg-red-500" : "bg-primary"}`}
                    style={{ width: `${job.progress}%` }}
                  />
                </div>
                <div className="mt-1 text-xs text-muted">{job.progress}%</div>
              </div>

              {job.error && (
                <div className="mt-2 text-sm text-red-600">{job.error}</div>
              )}
            </div>
          );
        })}
      </div>

      {jobs.length === 0 && !loading && (
        <div className="text-sm text-muted">No jobs yet. Try starting one from the Apps &gt; Image Compress tool.</div>
      )}
    </div>
  );
}

function StatusBadge({ status }: { status: Job["status"] }) {
  const label =
    status === "queued" ? "Queued" :
    status === "in_progress" ? "In Progress" :
    status === "completed" ? "Completed" :
    status === "canceled" ? "Canceled" :
    "Failed";
  const cls =
    status === "completed" ? "bg-green-100 text-green-800" :
    status === "failed" ? "bg-red-100 text-red-700" :
    status === "canceled" ? "bg-gray-100 text-gray-700" :
    "bg-blue-100 text-blue-800";
  return <span className={`text-xs px-2 py-0.5 rounded ${cls}`}>{label}</span>;
}
