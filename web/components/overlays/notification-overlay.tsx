"use client";

import React, { useMemo } from "react";
import BellIcon from "@/components/icons/bell";

export type NotificationStatus = "On Queue" | "On Processing" | "Finished" | "Failed";

export type NotificationItem = {
  id: string;
  title: string;
  createdAt: Date | string | number;
  status: NotificationStatus;
  icon?: React.ReactNode;
};

function toDate(input: Date | string | number): Date {
  if (input instanceof Date) return input;
  const d = new Date(input);
  return isNaN(d.getTime()) ? new Date() : d;
}

function timeAgo(input: Date | string | number): string {
  const now = new Date().getTime();
  const t = toDate(input).getTime();
  const diff = Math.max(0, Math.floor((now - t) / 1000));
  if (diff < 60) return `${diff}s ago`;
  const m = Math.floor(diff / 60);
  if (m < 60) return `${m}m ago`;
  const h = Math.floor(m / 60);
  if (h < 24) return `${h}h ago`;
  const d = Math.floor(h / 24);
  if (d < 30) return `${d}d ago`;
  const mo = Math.floor(d / 30);
  if (mo < 12) return `${mo}mo ago`;
  const y = Math.floor(mo / 12);
  return `${y}y ago`;
}

function statusBadgeClass(status: NotificationStatus): string {
  switch (status) {
    case "On Queue":
    case "On Processing":
      return "bg-primary/10 text-primary border border-primary/20";
    case "Finished":
      return "bg-foreground/10 text-foreground border border-border";
    case "Failed":
      return "bg-foreground/10 text-foreground border border-border";
    default:
      return "bg-foreground/10 text-foreground border border-border";
  }
}

export default function NotificationOverlay({
  title = "Notifications",
  items,
}: {
  title?: string;
  items?: NotificationItem[];
}) {
  const fallbackItems = useMemo<NotificationItem[]>(() => {
    const now = Date.now();
    return [
      {
        id: "1",
        title: "New project assigned: Phoenix",
        createdAt: new Date(now - 5 * 60 * 1000), // 5m ago
        status: "On Queue",
        icon: <BellIcon className="w-6 h-6 text-primary" aria-hidden="true" />,
      },
      {
        id: "2",
        title: "Build succeeded on main",
        createdAt: new Date(now - 2 * 60 * 60 * 1000), // 2h ago
        status: "Finished",
        icon: <BellIcon className="w-6 h-6 text-primary" aria-hidden="true" />,
      },
      {
        id: "3",
        title: "You have 3 unread messages",
        createdAt: new Date(now - 26 * 60 * 60 * 1000), // 26h ago
        status: "On Processing",
        icon: <BellIcon className="w-6 h-6 text-primary" aria-hidden="true" />,
      },
      {
        id: "4",
        title: "Deployment failed: api-service",
        createdAt: new Date(now - 3 * 24 * 60 * 60 * 1000), // 3d ago
        status: "Failed",
        icon: <BellIcon className="w-6 h-6 text-primary" aria-hidden="true" />,
      },
      {
        id: "5",
        title: "Weekly summary is ready",
        createdAt: new Date(now - 10 * 24 * 60 * 60 * 1000), // 10d ago
        status: "Finished",
        icon: <BellIcon className="w-6 h-6 text-primary" aria-hidden="true" />,
      },
    ];
  }, []);

  const data = items && items.length > 0 ? items : fallbackItems;

  return (
    <div className="h-full w-full bg-card shadow-primary rounded-lg">
      <div className="flex flex-col">
        <header className="sticky top-0 z-10 px-5 py-5 bg-primary shadow-lg rounded-t-xl">
          <h2 className="text-lg font-bold text-primary-foreground">
            Notifications
          </h2>
        </header>

        <div className="flex-1 px-2 py-4 space-y-1">
            {data.map((n) => (
              <div key={n.id} className="flex items-center gap-3 rounded-lg px-3 py-2 hover:bg-foreground/10 transition">
                <div className="shrink-0 w-10 h-10 rounded-full bg-border flex items-center justify-center">
                  {n.icon ?? <BellIcon className="w-6 h-6 text-foreground" aria-hidden="true"/>}
                </div>
                <div className="min-w-0 flex-1">
                  <div className="flex flex-col">
                    <p className="truncate text-sm text-foreground">{n.title}</p>
                    <span className="text-xs text-muted whitespace-nowrap">{timeAgo(n.createdAt)}</span>
                  </div>
                </div>
                <div className="ml-2">
                  <span
                    className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${statusBadgeClass(n.status)}${n.status === "On Processing" ? " animate-pulse" : ""}`}
                    aria-label={`status: ${n.status}`}
                  >
                    {n.status}
                  </span>
                </div>
              </div>
            ))}
        </div>
      </div>
    </div>
  );
}
