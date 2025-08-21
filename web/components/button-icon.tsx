"use client";

import React from "react";
import clsx from "clsx";


export default function ButtonIcon({
  className,
  children,
  type = "button",
  ...rest
}: React.ButtonHTMLAttributes<HTMLButtonElement>) {
  return (
    <button
      type={type}
      className={clsx(
        "bg-[var(--btn-secondary-bg)]",
        "text-[var(--btn-secondary-text)]",
        "hover:bg-[var(--btn-secondary-hover)]",
        "active:bg-[var(--btn-secondary-active)]",
        "disabled:bg-[var(--btn-secondary-disabled)]",
        "border border-[var(--btn-border)]",
        "p-2 rounded-full cursor-pointer shadow-primary",
        className
      )}
      {...rest}
    >
      {children}
    </button>
  );
}
