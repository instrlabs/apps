"use client";

import React from "react";

type ButtonIconProps = React.ButtonHTMLAttributes<HTMLButtonElement> & {
  "aria-label"?: string;
};

function cx(...classes: Array<string | undefined | false | null>) {
  return classes.filter(Boolean).join(" ");
}

export default function ButtonIcon({ className, children, type = "button", ...rest }: ButtonIconProps) {
  const base = "p-2 rounded-full hover:bg-blue-50 focus:outline-none";
  return (
    <button type={type} className={cx(base, className)} {...rest}>
      {children}
    </button>
  );
}
