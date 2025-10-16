"use client";

import React from "react";


export default function NotificationOverlay() {
  return (
    <div className="h-full w-full bg-card shadow-primary rounded-lg">
      <div className="flex flex-col">
        <header className="sticky top-0 z-10 px-5 py-5 bg-primary shadow-lg rounded-t-xl">
          <h2 className="text-lg font-bold text-primary-foreground">
            Notifications
          </h2>
        </header>


      </div>
    </div>
  );
}
