"use client";

import React from "react";

export default function NotificationOverlay() {
  return (
    <div className="w-full md:w-[400px] h-full bg-white/10 border border-white/10 rounded-lg p-4 flex flex-col gap-4">
      <div className="flex justify-between">
        <span className="text-sm font-semibold">Notifications</span>
        <button type="button" className="text-sm text-white/50 hover:text-white transition-colors">Clear All</button>
      </div>
      <div className="flex-1 flex flex-col items-center justify-center">
        <span className="text-sm text-white/50">Empty</span>
      </div>
    </div>
  );
}
