"use client";

import React from "react";
import NavigationOverlay from "@/components/reuse/navigation-overlay";
import NotificationOverlay from "@/components/reuse/notification-overlay";
import ProfileOverlay from "@/components/reuse/profile-overlay";
import TextField from "@/components/text-field";
import SearchIcon from "@/components/icons/search";

export type OverlaySide = "left" | "right" | "modal";
export type OverlayAction = () => React.ReactNode;

export type OverlayEntry = {
  side: OverlaySide;
  width: number;
  action: OverlayAction;
};

const registry = new Map<string, OverlayEntry>();

export function getOverlayEntry(key: string): OverlayEntry | undefined {
  return registry.get(key);
}

export function resetOverlayRegistry() {
  registry.clear();
}

export function resolveOverlayNode(entry: OverlayEntry): React.ReactNode {
  return entry.action();
}

function registerOverlay(key: string, config: OverlayEntry): void {
  const [overlaySide] = config.side.split(":", 1);

  registry.set(key, {
    side: overlaySide as OverlaySide,
    width: config.width,
    action: config.action
  });
}

export function registerBuiltInOverlays() {
  // Left navigation rail
  registerOverlay("left:navigation", {
    side: "left",
    width: 80,
    action: () => <NavigationOverlay />,
  });

  // Right notifications
  registerOverlay("right:notifications", {
    side: "right",
    width: 400,
    action: () => <NotificationOverlay />,
  });

  // Right profile
  registerOverlay("right:profile", {
    side: "right",
    width: 400,
    action: () => <ProfileOverlay />,
  });

  // Modal: search
  registerOverlay("modal:search", {
    side: "modal",
    width: 400,
    action: () => (
      <div className="space-y-3">
        <div className="relative">
          <label htmlFor="global-search-input" className="sr-only">Search</label>
          <TextField
            id="global-search-input"
            type="text"
            autoFocus
            placeholder="Search..."
            className="pr-10"
            xSize="md"
          />
          <SearchIcon
            className="pointer-events-none absolute right-3 top-1/2 -translate-y-1/2 w-5 h-5 text-muted"
            aria-hidden="true"
          />
        </div>
        <div className="text-sm text-muted">Type to search…</div>
      </div>
    ),
  });
}
