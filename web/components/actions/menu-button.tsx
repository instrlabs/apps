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
    xSize === "sm" ? "p-2 text-sm" :
      xSize === "md" ? "p-2 gap-3" :
        xSize === "lg" ? "p-4 gap-4 text-base" :
          "";

  return (
    <button
      type={type}
      className={clsx("w-full px-2", className)}
      {...rest}
    >
      <div
        className={clsx(
          "w-full flex items-center rounded-lg text-left font-light text-white hover:bg-white/10 hover:text-white cursor-pointer transition-colors",
          sizeClass
        )}
      >
        {icon ? <span className="shrink-0">{icon}</span> : null}
        <span className="truncate">{children}</span>
      </div>
    </button>
  );
}
