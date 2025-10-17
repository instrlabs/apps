"use client";

import React from "react";

export default function OverlayBody({ children }: {
  children: React.ReactNode;
}) {
  return (
    <div className="h-full flex flex-col transition-width duration-300 ease-in-out">
      {children}
    </div>
  );
}
