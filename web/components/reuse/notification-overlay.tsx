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
      return "bg-blue-100 text-blue-800 border-blue-200";
    case "On Processing":
      return "bg-amber-100 text-amber-800 border-amber-200";
    case "Finished":
      return "bg-green-100 text-green-800 border-green-200";
    case "Failed":
      return "bg-red-100 text-red-800 border-red-200";
    default:
      return "bg-green-100 text-green-800 border-green-200";
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
        icon: <BellIcon className="w-6 h-6 text-blue-600" aria-hidden="true" />,
      },
      {
        id: "2",
        title: "Build succeeded on main",
        createdAt: new Date(now - 2 * 60 * 60 * 1000), // 2h ago
        status: "Finished",
        icon: <BellIcon className="w-6 h-6 text-green-600" aria-hidden="true" />,
      },
      {
        id: "3",
        title: "You have 3 unread messages",
        createdAt: new Date(now - 26 * 60 * 60 * 1000), // 26h ago
        status: "On Processing",
        icon: <BellIcon className="w-6 h-6 text-indigo-600" aria-hidden="true" />,
      },
      {
        id: "4",
        title: "Deployment failed: api-service",
        createdAt: new Date(now - 3 * 24 * 60 * 60 * 1000), // 3d ago
        status: "Failed",
        icon: <BellIcon className="w-6 h-6 text-red-600" aria-hidden="true" />,
      },
      {
        id: "5",
        title: "Weekly summary is ready",
        createdAt: new Date(now - 10 * 24 * 60 * 60 * 1000), // 10d ago
        status: "Finished",
        icon: <BellIcon className="w-6 h-6 text-gray-600" aria-hidden="true" />,
      },
    ];
  }, []);

  const data = items && items.length > 0 ? items : fallbackItems;

  return (
    <section className="h-full flex flex-col bg-blue-50" aria-labelledby="notifications-title">
      <header className="sticky top-0 z-10 backdrop-blur px-5 py-5 border-b border-gray-200">
        <h2 id="notifications-title" className="text-base font-semibold text-gray-900">
          Notifications
        </h2>
      </header>

      <div className="flex-1 overflow-auto p-2">
        <ul role="list" className="space-y-1">
          {data.map((n) => (
            <li key={n.id} role="listitem">
              <div className="flex items-center gap-3 rounded-xl px-3 py-2 hover:bg-blue-200 transition">
                <div className="shrink-0 w-10 h-10 rounded-full bg-gray-100 flex items-center justify-center">
                  {n.icon ?? <BellIcon className="w-6 h-6 text-gray-700" aria-hidden="true" />}
                </div>
                <div className="min-w-0 flex-1">
                  <div className="flex flex-col">
                    <p className="truncate text-sm text-gray-900">{n.title}</p>
                    <span className="text-xs text-gray-500 whitespace-nowrap">{timeAgo(n.createdAt)}</span>
                  </div>
                </div>
                <div className="ml-2">
                  <span
                    className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium border ${statusBadgeClass(n.status)}${n.status === "On Processing" ? " animate-pulse" : ""}`}
                    aria-label={`status: ${n.status}`}
                  >
                    {n.status}
                  </span>
                </div>
              </div>
            </li>
          ))}
        </ul>
      </div>
    </section>
  );
}
