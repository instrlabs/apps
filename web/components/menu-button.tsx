"use client";

import React from "react";
import clsx from "clsx";

export type MenuButtonProps = React.ButtonHTMLAttributes<HTMLButtonElement> & {
  icon?: React.ReactNode;
  xSize?: "sm" | "md" | "lg";
};

export default function MenuButton({
  icon,
  className,
  type = "button",
  children,
  xSize = "md",
  ...rest
}: MenuButtonProps) {
  const sizeClass =
    xSize === "sm" ? "p-2 gap-2 text-sm" : xSize === "md" ? "p-3 gap-3" : "p-4 gap-4 text-base";

  return (
    <button
      type={type}
      className={clsx(
        "w-full flex items-center",
        sizeClass,
        "font-medium text-foreground",
        "bg-menu border border-border hover:bg-menu-hover hover:border-border-hover",
        "rounded-lg cursor-pointer",
        className
      )}
      {...rest}
    >
      {icon ? <span className="shrink-0">{icon}</span> : null}
      <span className="truncate">{children}</span>
    </button>
  );
}
