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
    <div className="w-full h-full bg-white py-3">
      <ul className="flex flex-col items-center space-y-2">
        {items.map((item) => (
          <li key={item.key}>
            <button
              type="button"
              className="relative group flex items-center center p-1 rounded-full hover:bg-blue-200 focus:outline-none focus-visible:ring-2 focus-visible:ring-blue-400"
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
                  className="pointer-events-none absolute left-full top-1/2 -translate-y-1/2 ml-2 px-2 py-1 rounded-md bg-gray-900 text-white text-xs opacity-0 group-hover:opacity-100 group-focus-within:opacity-100 transition-opacity whitespace-nowrap shadow-lg z-10"
                >
                  {item.title}
                </span>
              )}
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
}
