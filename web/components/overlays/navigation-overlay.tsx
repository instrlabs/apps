"use client";

import React from "react";
import { usePathname } from "next/navigation";
import Link from "next/link";

type NavItem = {
  key: string;
  title: string;
  icon?: React.ReactNode;
};

import AppsIcon from "@/components/icons/apps";
import HistoryIcon from "@/components/icons/history";
import StorageIcon from "@/components/icons/storage";
import clsx from "clsx";

export default function NavigationOverlay({
  items = [
    { key: "apps", title: "Apps", icon: <AppsIcon />},
    { key: "histories", title: "Histories", icon: <HistoryIcon /> },
    { key: "storage", title: "Storage", icon: <StorageIcon /> },
  ],
}: {
  items?: NavItem[];
}) {
  const pathname = usePathname();

  const keyToPath: Record<string, string> = {
    apps: "/",
    histories: "/histories",
    storage: "/storage",
  };

  return (
    <div className="w-[80px] h-full">
      <div className="flex flex-col items-center space-y-5">
        {items.map((item) => {
          const target = keyToPath[item.key];
          const isActive = target ? pathname?.startsWith(target) : false;

          return (
            <div key={item.key}>
              <Link
                href={target}
                className={clsx(
                  "relative group cursor-pointer",
                  "flex items-center rounded-full",
                  "hover:bg-foreground/5 focus:outline-none",
                  isActive ? "bg-black" : "bg-white",
                )}
              >
                <span className={clsx(
                  "inline-flex items-center justify-center w-10 h-10",
                  isActive ? "text-white" : "text-gray-600",
                )}>
                  {item.icon}
                </span>
                {item.title && (
                  <span
                    className={clsx(
                      "pointer-events-none",
                      "absolute left-full top-1/2 -translate-y-1/2",
                      "ml-2 px-2 py-1 rounded-md",
                      "bg-white text-primary text-xs font-semibold opacity-0",
                      "group-hover:opacity-100 transition-opacity whitespace-nowrap shadow-lg z-10"
                    )}
                  >
                    {item.title}
                  </span>
                )}
              </Link>
            </div>
          );
        })}
      </div>
    </div>
  );
}
