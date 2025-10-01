"use client"

import React from "react";

export default function Loading() {
  const containerClasses = "flex items-center justify-center";

  return (
    <div className={`${containerClasses}`} role="status" aria-live="polite" aria-busy="true">
      <div className="flex items-center gap-3 p-4 rounded-lg">
        <span className="inline-block animate-spin rounded-full border-2 border-gray-300 border-t-blue-500 size-10" />
      </div>
    </div>
  );
}
