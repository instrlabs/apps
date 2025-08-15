"use client";

import React from "react";
import clsx from "clsx";

export type MenuButtonProps = React.ButtonHTMLAttributes<HTMLButtonElement> & {
  icon?: React.ReactNode;
};

export default function MenuButton({
  icon,
  className,
  type = "button",
  children,
  ...rest
}: MenuButtonProps) {
  return (
    <button
      type={type}
      className={clsx(
        "w-full flex items-center gap-3 px-4 py-3",
        "text-sm font-medium text-gray-800",
        "bg-slate-50 hover:bg-blue-100",
        "rounded-sm cursor-pointer",
        className
      )}
      {...rest}
    >
      {icon ? <span className="shrink-0">{icon}</span> : null}
      <span className="truncate">{children}</span>
    </button>
  );
}
