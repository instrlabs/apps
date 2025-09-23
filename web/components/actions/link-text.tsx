import React from "react";
import Link from "next/link";
import clsx from "clsx";

type LinkTextProps = {
  href: string;
  children: React.ReactNode;
  className?: string;
};

export default function LinkText({ href, children, className }: LinkTextProps) {
  return (
    <Link href={href} className={clsx("text-blue-400", className)}>
      {children}
    </Link>
  );
}
