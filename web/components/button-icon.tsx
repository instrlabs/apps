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
        "p-2 rounded-full",
        "bg-white shadow-primary hover:bg-blue-100 focus:outline-none",
        "cursor-pointer",
        className
      )}
      {...rest}
    >
      {children}
    </button>
  );
}
