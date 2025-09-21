"use client"

import React from "react";

type LoadingProps = {
  message?: string;
  fullScreen?: boolean;
  size?: number; // px
  className?: string;
};

export default function Loading({ message = "Loading...", fullScreen = false, size = 32, className = "" }: LoadingProps) {
  const spinnerSize = Math.max(16, size);
  const containerClasses = fullScreen
    ? "fixed inset-0 z-50 flex items-center justify-center bg-white/70 backdrop-blur-sm"
    : "flex items-center justify-center";

  return (
    <div className={`${containerClasses} ${className}`.trim()} role="status" aria-live="polite" aria-busy="true">
      <div className="flex items-center gap-3 p-4 rounded-lg">
        <span
          className="inline-block animate-spin rounded-full border-2 border-gray-300 border-t-blue-500"
          style={{ width: spinnerSize, height: spinnerSize }}
        />
        {message ? <span className="text-sm text-gray-700">{message}</span> : null}
      </div>
    </div>
  );
}
