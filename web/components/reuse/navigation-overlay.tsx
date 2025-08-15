"use client";

import React from "react";

type NavItem = {
  key: string;
  title: string; // kept for accessibility label
  icon?: React.ReactNode;
  onClick?: () => void;
};

// Provided icons (24px) from the issue description
function AppsIcon({ className = "w-6 h-6" }: { className?: string }) {
  return (
    <svg xmlns="http://www.w3.org/2000/svg" enableBackground="new 0 0 24 24" height="24px" viewBox="0 0 24 24" width="24px" fill="#222" className={className} aria-hidden="true">
      <g><rect fill="none" height="24" width="24"/></g>
      <g><g><path d="M5,11h4c1.1,0,2-0.9,2-2V5c0-1.1-0.9-2-2-2H5C3.9,3,3,3.9,3,5v4C3,10.1,3.9,11,5,11z"/><path d="M5,21h4c1.1,0,2-0.9,2-2v-4c0-1.1-0.9-2-2-2H5c-1.1,0-2,0.9-2,2v4C3,20.1,3.9,21,5,21z"/><path d="M13,5v4c0,1.1,0.9,2,2,2h4c1.1,0,2-0.9,2-2V5c0-1.1-0.9-2-2-2h-4C13.9,3,13,3.9,13,5z"/><path d="M15,21h4c1.1,0,2-0.9,2-2v-4c0-1.1-0.9-2-2-2h-4c-1.1,0-2,0.9-2,2v4C13,20.1,13.9,21,15,21z"/></g></g>
    </svg>
  );
}

function HistoryIcon({ className = "w-6 h-6" }: { className?: string }) {
  return (
    <svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 0 24 24" width="24px" fill="#222" className={className} aria-hidden="true">
      <path d="M0 0h24v24H0V0z" fill="none"/>
      <path d="M13.26 3C8.17 2.86 4 6.95 4 12H2.21c-.45 0-.67.54-.35.85l2.79 2.8c.2.2.51.2.71 0l2.79-2.8c.31-.31.09-.85-.36-.85H6c0-3.9 3.18-7.05 7.1-7 3.72.05 6.85 3.18 6.9 6.9.05 3.91-3.1 7.1-7 7.1-1.61 0-3.1-.55-4.28-1.48-.4-.31-.96-.28-1.32.08-.42.42-.39 1.13.08 1.49C9 20.29 10.91 21 13 21c5.05 0 9.14-4.17 9-9.26-.13-4.69-4.05-8.61-8.74-8.74zm-.51 5c-.41 0-.75.34-.75.75v3.68c0 .35.19.68.49.86l3.12 1.85c.36.21.82.09 1.03-.26.21-.36.09-.82-.26-1.03l-2.88-1.71v-3.4c0-.4-.34-.74-.75-.74z"/>
    </svg>
  );
}

function StorageIcon({ className = "w-6 h-6" }: { className?: string }) {
  return (
    <svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 0 24 24" width="24px" fill="#222" className={className} aria-hidden="true">
      <path d="M0 0h24v24H0V0z" fill="none"/>
      <path d="M20.54 5.23l-1.39-1.68C18.88 3.21 18.47 3 18 3H6c-.47 0-.88.21-1.16.55L3.46 5.23C3.17 5.57 3 6.02 3 6.5V19c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V6.5c0-.48-.17-.93-.46-1.27zm-8.89 11.92L6.5 12H10v-2h4v2h3.5l-5.15 5.15c-.19.19-.51.19-.7 0zM5.12 5l.81-1h12l.94 1H5.12z"/>
    </svg>
  );
}

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
    <ul className="space-y-2 text-sm h-full">
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
  );
}
