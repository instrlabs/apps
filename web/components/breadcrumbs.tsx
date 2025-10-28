"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import React from "react";

export default function Breadcrumbs() {
  const pathname = usePathname() || "/";
  const isHome = pathname === "/";

  if (isHome) {
    return (
      <nav aria-label="Breadcrumbs" className="flex items-center gap-2.5 p-2">
        <span className="text-base leading-6 font-semibold text-white">
          Instruction Labs
        </span>
      </nav>
    );
  }

  const paths = pathname.split("/").filter(Boolean);
  const breadcrumbs = [
    { label: "Instruction Labs", href: "/", isFirst: true },
    ...paths.map((segment, index) => {
      const currentPath = "/" + paths.slice(0, index + 1).join("/");
      return {
        label: segment.charAt(0).toUpperCase() + segment.slice(1),
        href: index === paths.length - 1 ? null : currentPath,
        isFirst: false,
      };
    }),
  ];

  return (
    <nav aria-label="Breadcrumbs" className="flex items-center gap-2.5 p-2">
      <div className="flex items-center gap-2.5">
        {breadcrumbs.map((item, index) => {
          const isLast = index === breadcrumbs.length - 1;

          return (
            <React.Fragment key={`breadcrumb-${index}`}>
              {item.href ? (
                <Link
                  href={item.href}
                  className={[
                    "leading-6 text-blue-200 transition-colors hover:text-blue-100",
                    item.isFirst
                      ? "text-base font-semibold"
                      : "text-base font-medium",
                  ]
                    .filter(Boolean)
                    .join(" ")}
                >
                  {item.label}
                </Link>
              ) : (
                <span
                  className={[
                    "leading-6 text-white",
                    item.isFirst
                      ? "text-base font-semibold"
                      : "text-base font-medium",
                  ]
                    .filter(Boolean)
                    .join(" ")}
                >
                  {item.label}
                </span>
              )}
              {!isLast && (
                <span className="text-xs leading-4 font-semibold text-white">
                  /
                </span>
              )}
            </React.Fragment>
          );
        })}
      </div>
    </nav>
  );
}
