"use client";

import React from "react";
import { usePathname } from "next/navigation";

type NavItem = {
  key: string;
  title: string; // kept for accessibility label
  icon?: React.ReactNode;
  onClick?: () => void;
};

import AppsIcon from "@/components/icons/apps";
import HistoryIcon from "@/components/icons/history";
import StorageIcon from "@/components/icons/storage";
import clsx from "clsx";

export default function NavigationOverlay({
  items = [
    { key: "apps", title: "Apps", icon: <AppsIcon /> },
    { key: "histories", title: "Histories", icon: <HistoryIcon /> },
    { key: "storage", title: "Storage", icon: <StorageIcon /> },
  ],
}: {
  items?: NavItem[];
}) {
  const pathname = usePathname();

  // Map known item keys to their routes. Route groups like (site) are not part of the URL.
  const keyToPath: Record<string, string> = {
    apps: "/apps",
    histories: "/histories",
    storage: "/storage",
  };

  return (
    <div className="w-full h-full bg-card shadow-primary rounded-xl py-3">
      <div className="flex flex-col items-center space-y-2">
        {items.map((item) => {
          const target = keyToPath[item.key];
          const isActive = target ? pathname?.startsWith(target) : false;

          return (
            <div key={item.key}>
              <button
                type="button"
                className={clsx(
                  "relative group flex items-center rounded-xl",
                  "hover:bg-foreground/5 focus:outline-none",
                  "cursor-pointer",
                  isActive && "bg-gray-50"
                )}
                onClick={item.onClick}
                aria-label={item.title}
                aria-current={isActive ? "page" : undefined}
              >
                <span
                  className="inline-flex items-center justify-center w-10 h-10 text-foreground"
                  aria-hidden="true"
                >
                  {item.icon}
                </span>
                {item.title && (
                  <span
                    role="tooltip"
                    className="pointer-events-none absolute left-full top-1/2 -translate-y-1/2 ml-2 px-2 py-1 rounded-md bg-primary text-primary-foreground text-xs opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap shadow-lg z-10"
                  >
                    {item.title}
                  </span>
                )}
              </button>
            </div>
          );
        })}
      </div>
    </div>
  );
}
