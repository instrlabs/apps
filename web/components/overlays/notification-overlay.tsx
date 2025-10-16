"use client";

import React from "react";
import Text from "@/components/text";


export default function NotificationOverlay() {
  return (
    <div className="w-[400px] h-full bg-white/10 border border-white/10 rounded-lg p-4 flex flex-col gap-4">
      <div className="flex justify-between">
        <Text xSize="sm" isBold>Notifications</Text>
        <Text xSize="sm" xColor="secondary">Clear All</Text>
      </div>
      <div className="flex-1 flex flex-col items-center justify-center">
        <Text xSize="sm" xColor="secondary">Empty</Text>
      </div>
    </div>
  );
}
