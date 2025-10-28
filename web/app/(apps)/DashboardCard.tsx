"use client";

import React from "react";
import Link from "next/link";
import Icon from "@/components/icon";

export type DashboardCardProps = {
  title: string;
  description: string;
  href: string;
  iconName?: string;
};

export default function DashboardCard({
  title,
  description,
  href,
  iconName = "rectangle",
}: DashboardCardProps) {
  return (
    <Link
      href={href}
      className="flex flex-col items-start gap-2 rounded-lg border border-white/10 bg-white/8 p-4 transition-colors hover:bg-white/12"
    >
      <div className="relative size-10 shrink-0">
        <Icon name={iconName} size={40} />
      </div>
      <h3 className="text-lg leading-7 font-semibold text-white">{title}</h3>
      <p className="text-sm leading-5 font-normal text-white">{description}</p>
    </Link>
  );
}
