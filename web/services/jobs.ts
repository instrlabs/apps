// Simple in-memory job store and simulator for demo purposes
// NOTE: This is not persisted and will reset on server restart or deployment.

export type JobStatus = "queued" | "in_progress" | "completed" | "canceled" | "failed";

export interface Job {
  id: string;
  type: string; // e.g., "image-compress"
  filename: string;
  status: JobStatus;
  progress: number; // 0..100
  createdAt: number;
  updatedAt: number;
  error?: string;
}

// Internal store and timers
const jobs = new Map<string, Job>();
const timers = new Map<string, NodeJS.Timeout>();

function genId() {
  return Math.random().toString(36).slice(2) + Date.now().toString(36);
}

function updateJob(id: string, patch: Partial<Job>) {
  const j = jobs.get(id);
  if (!j) return;
  const next = { ...j, ...patch, updatedAt: Date.now() } as Job;
  jobs.set(id, next);
}

function clearTimer(id: string) {
  const t = timers.get(id);
  if (t) {
    clearInterval(t);
    timers.delete(id);
  }
}

export function listJobs(type?: string): Job[] {
  const arr = Array.from(jobs.values()).sort((a, b) => b.createdAt - a.createdAt);
  return type ? arr.filter(j => j.type === type) : arr;
}

export function getJob(id: string): Job | undefined {
  return jobs.get(id);
}

export function cancelJob(id: string): Job | undefined {
  const job = jobs.get(id);
  if (!job) return undefined;
  if (job.status === "completed" || job.status === "failed" || job.status === "canceled") return job;
  updateJob(id, { status: "canceled" });
  clearTimer(id);
  return jobs.get(id);
}

export interface CreateJobInput {
  type: string;
  filename: string;
}

export function createJob(input: CreateJobInput): Job {
  const id = genId();
  const now = Date.now();
  const job: Job = {
    id,
    type: input.type,
    filename: input.filename,
    status: "queued",
    progress: 0,
    createdAt: now,
    updatedAt: now,
  };
  jobs.set(id, job);

  // Simulate async processing
  // Transition to in_progress quickly, then increment progress until completion
  const startDelay = 300; // ms
  const intervalMs = 500; // ms

  const startTimer = setTimeout(() => {
    updateJob(id, { status: "in_progress", progress: 5 });

    const interval = setInterval(() => {
      const current = jobs.get(id);
      if (!current) {
        clearInterval(interval);
        return;
      }
      if (current.status === "canceled") {
        clearInterval(interval);
        return;
      }
      if (current.status === "failed" || current.status === "completed") {
        clearInterval(interval);
        return;
      }

      const increment = 10 + Math.floor(Math.random() * 15); // 10..24
      const nextProgress = Math.min(100, current.progress + increment);
      if (nextProgress >= 100) {
        updateJob(id, { progress: 100, status: "completed" });
        clearInterval(interval);
        timers.delete(id);
      } else {
        updateJob(id, { progress: nextProgress });
      }
    }, intervalMs);

    timers.set(id, interval);
  }, startDelay);

  // Track the startup timeout to be cleared if needed
  timers.set(id, startTimer as unknown as NodeJS.Timeout);

  return job;
}
