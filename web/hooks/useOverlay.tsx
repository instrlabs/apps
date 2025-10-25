"use client";

import React, {createContext, useCallback, useContext, useMemo, useState} from "react";

import NavigationOverlay from "@/components/widgets/navigation-overlay";
import NotificationOverlay from "@/components/widgets/notification-overlay";
import ProfileOverlay from "@/components/widgets/profile-overlay";

export type OverlayActions = {
  // left overlay state
  isLeftOpen: boolean;
  leftKey: string;
  leftNode: React.ReactNode;
  openLeft: (key: string) => void;
  closeLeft: () => void;
  // right overlay state
  isRightOpen: boolean;
  rightKey: string;
  rightNode: React.ReactNode;
  openRight: (key: string) => void;
  closeRight: () => void;
};

export type OverlaySide = "left" | "right";
export type OverlayRender = () => React.ReactNode;
export type OverlayOpts = {
  side: OverlaySide;
  render: OverlayRender;
};

const overlay = new Map<string, OverlayOpts>();

function registerOverlay(key: string, config: OverlayOpts): void {
  const [overlaySide] = config.side.split(":", 1);

  overlay.set(key, {
    side: overlaySide as OverlaySide,
    render: config.render,
  });
}

function registerOverlays() {
  // Left navigation rail
  registerOverlay("left:navigation", {
    side: "left",
    render: () => <NavigationOverlay />,
  });

  // Right notifications
  registerOverlay("right:notifications", {
    side: "right",
    render: () => <NotificationOverlay />,
  });

  // Right profile
  registerOverlay("right:profile", {
    side: "right",
    render: () => <ProfileOverlay />,
  });
}

const OverlayContext = createContext<OverlayActions | undefined>(undefined);

export function OverlayProvider({ children, defaultLeft, defaultRight }: {
  children: React.ReactNode;
  defaultLeft?: string;
  defaultRight?: string;
}) {
  registerOverlays();

  let defaultLeftOverlay;
  if (defaultLeft) defaultLeftOverlay = overlay.get(defaultLeft);
  let defaultRightOverlay;
  if (defaultRight) defaultRightOverlay = overlay.get(defaultRight);

  // left
  const [isLeftOpen, setIsLeftOpen] = useState<boolean>(!!defaultLeftOverlay);
  const [leftNode, setLeftNodeState] = useState<React.ReactNode>(!!defaultLeftOverlay ? (defaultLeftOverlay as OverlayOpts).render() : <div />);
  const [leftKey, setLeftKey] = useState<string>(defaultLeft ?? "");

  // right
  const [isRightOpen, setIsRightOpen] = useState<boolean>(!!defaultRightOverlay);
  const [rightNode, setRightNodeState] = useState<React.ReactNode>(!!defaultRightOverlay ? (defaultRightOverlay as OverlayOpts).render() : <div/>);
  const [rightKey, setRightKey] = useState<string>(defaultRight ?? "");

  const openLeft = useCallback((currentKey: string) => {
    const key = currentKey ?? "";
    const content = overlay.get(key) as OverlayOpts;
    setLeftNodeState(content.render());
    setLeftKey(key);
    setIsLeftOpen(true);
  }, []);

  const closeLeft = useCallback(() => {
    setLeftNodeState(<div/>);
    setLeftKey("");
    setIsLeftOpen(false);
  }, []);

  const openRight = useCallback((currentKey: string) => {
    const key = `right:${currentKey ?? ""}`;
    const content = overlay.get(key) as OverlayOpts;
    setRightNodeState(content.render());
    setRightKey(key);
    setIsRightOpen(true);
  }, []);

  const closeRight = useCallback(() => {
    setRightNodeState(<div/>);
    setRightKey("");
    setIsRightOpen(false);
  }, []);

  const value = useMemo<OverlayActions>(() => ({
    // left
    isLeftOpen,
    leftKey,
    leftNode,
    openLeft,
    closeLeft,
    // right
    isRightOpen,
    rightKey,
    rightNode,
    openRight,
    closeRight,
  }), [isLeftOpen, leftKey, leftNode, closeLeft, openLeft, isRightOpen, rightKey, rightNode, openRight, closeRight]);

  return (
    <OverlayContext.Provider value={value}>
      {children}
    </OverlayContext.Provider>
  );
}

export function useOverlay(): OverlayActions {
  const ctx = useContext(OverlayContext);
  if (!ctx) throw new Error('useOverlay must be used within an OverlayProvider');
  return ctx;
}
