"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import React from "react";

export default function Breadcrumbs() {
  const pathname = (usePathname() || "/");

  const isHome = pathname === "/";
  const paths = pathname.split("/");

  return (
    <nav
      aria-label="Breadcrumbs"
      className="pointer-events-auto flex items-center gap-2 px-2 py-1"
    >
      {isHome ? (
        <span className="text-muted text-sm font-normal">Instruction Labs</span>
      ) : (
        <div className="flex items-center gap-2">
          <Link
            href="/"
            className="text-muted hover:text-foreground text-sm font-normal transition-colors"
          >
            Home
          </Link>
          <span className="text-muted text-[10px] leading-[20px]">/</span>
          <span className="text-foreground text-sm font-normal capitalize">{paths[1]}</span>
          <span className="text-muted text-[10px] leading-[20px]">/</span>
          <span className="text-foreground text-sm font-normal capitalize">{paths[2]}</span>
        </div>
      )}
    </nav>
  );
}
