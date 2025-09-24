"use client";

import React from "react";


export default function InlineSpinner() {
  return (
    <div className="flex items-center justify-center size-4 overflow-hidden">
      <span className="inline-block aspect-square size-4 animate-spin rounded-full border-2 border-gray-500 border-t-black" />
    </div>
  );
}
