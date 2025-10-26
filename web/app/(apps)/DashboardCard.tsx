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
  iconName = "circle",
}: DashboardCardProps) {
  return (
    <Link
      href={href}
      className="flex flex-col gap-2 items-start p-4 bg-white/8 border border-white/10 rounded-lg transition-colors hover:bg-white/12"
    >
      <div className="relative shrink-0 size-10">
        <Icon name={iconName} size={40} />
      </div>
      <h3 className="text-lg leading-7 font-semibold text-white">
        {title}
      </h3>
      <p className="text-sm leading-5 font-normal text-white">
        {description}
      </p>
    </Link>
  );
}
