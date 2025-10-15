"use client";

import Link from "next/link";
import clsx from "clsx";
import React from "react";

export type AppsCardProps = {
  href: string;
  title: string;
  description?: string;
  className?: string;
};

export default function AppsCard({ href, title, description, className }: AppsCardProps) {
  return (
    <Link href={href} className={clsx("group block h-full", className)}>
      <div
        className={clsx(
          "h-full p-4 rounded-lg shadow-primary bg-white/10 border border-white/10",
          "flex flex-col gap-2 transition-shadow hover:shadow-hover focus-within:shadow-focus"
        )}
      >
        <div className="w-10 h-10 rounded-md bg-white/10 border border-white/10 flex-none" aria-hidden="true" />
        <h3 className="text-lg font-semibold text-white">{title}</h3>
        <p className="text-white font-normal text-sm mt-1">{description}</p>
      </div>
    </Link>
  );
}
