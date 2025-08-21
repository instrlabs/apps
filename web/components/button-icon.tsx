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
        "p-2 rounded-xl",
        "bg-card shadow-primary hover:bg-foreground/5 focus:outline-none",
        "cursor-pointer",
        // Ensure solid icon color in light mode, and readable color in dark mode
        "text-gray-950 [data-theme=\"dark\"]:text-foreground",
        className
      )}
      {...rest}
    >
      {children}
    </button>
  );
}
