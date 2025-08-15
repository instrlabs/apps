"use client";

import React from "react";

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
  return (
    <div className="w-full h-full bg-white shadow-primary rounded-xl py-3">
      <div className="flex flex-col items-center space-y-2">
        {items.map((item) => (
          <div key={item.key}>
            <button
              type="button"
              className={clsx(
                "relative group flex items-center rounded-full",
                "hover:bg-blue-100 focus:outline-none",
                "cursor-pointer",
              )}
              onClick={item.onClick}
              aria-label={item.title}
            >
              <span
                className="inline-flex items-center justify-center w-10 h-10 text-gray-700"
                aria-hidden="true"
              >
                {item.icon}
              </span>
              {item.title && (
                <span
                  role="tooltip"
                  className="pointer-events-none absolute left-full top-1/2 -translate-y-1/2 ml-2 px-2 py-1 rounded-md bg-gray-800 text-white text-xs opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap shadow-lg z-10"
                >
                  {item.title}
                </span>
              )}
            </button>
          </div>
        ))}
      </div>
    </div>
  );
}
