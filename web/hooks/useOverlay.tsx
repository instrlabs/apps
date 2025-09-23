"use client";

import React, {createContext, useCallback, useContext, useMemo, useRef, useState} from "react";

import NavigationOverlay from "@/components/overlays/navigation-overlay";
import NotificationOverlay from "@/components/overlays/notification-overlay";
import ProfileOverlay from "@/components/overlays/profile-overlay";

export type OverlayActions = {
  // left overlay state
  isLeftOpen: boolean;
  leftKey: string;
  leftNode: React.ReactNode;
  leftWidth: number;
  openLeft: (key: string) => void;
  closeLeft: () => void;
  // right overlay state
  isRightOpen: boolean;
  rightKey: string;
  rightNode: React.ReactNode;
  rightWidth: number;
  openRight: (key: string) => void;
  closeRight: () => void;
};

export type OverlaySide = "left" | "right";
export type OverlayRender = () => React.ReactNode;
export type OverlayOpts = {
  side: OverlaySide;
  width: number;
  render: OverlayRender;
};

const overlay = new Map<string, OverlayOpts>();

function registerOverlay(key: string, config: OverlayOpts): void {
  const [overlaySide] = config.side.split(":", 1);

  overlay.set(key, {
    side: overlaySide as OverlaySide,
    width: config.width,
    render: config.render,
  });
}

function registerOverlays() {
  // Left navigation rail
  registerOverlay("left:navigation", {
    side: "left",
    width: 80,
    render: () => <NavigationOverlay />,
  });

  // Right notifications
  registerOverlay("right:notifications", {
    side: "right",
    width: 400,
    render: () => <NotificationOverlay />,
  });

  // Right profile
  registerOverlay("right:profile", {
    side: "right",
    width: 400,
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
  const [leftNode, setLeftNodeState] = useState<React.ReactNode>(!!defaultLeftOverlay ? defaultLeftOverlay.render() : <div />);
  const [leftKey, setLeftKey] = useState<string>(defaultLeft ?? "");
  const [leftWidth, setLeftWidth] = useState<number>(!!defaultLeftOverlay ? defaultLeftOverlay.width : 0);

  // right
  const [isRightOpen, setIsRightOpen] = useState<boolean>(!!defaultRightOverlay);
  const [rightNode, setRightNodeState] = useState<React.ReactNode>(!!defaultRightOverlay ? defaultRightOverlay.render() : <div/>);
  const [rightKey, setRightKey] = useState<string>(defaultRight ?? "");
  const [rightWidth, setRightWidth] = useState<number>(!!defaultRightOverlay ? defaultRightOverlay.width : 0);

  const openLeft = useCallback((currentKey: string) => {
    const key = currentKey ?? "";
    const content = overlay.get(key) as OverlayOpts;
    setLeftNodeState(content.render());
    setLeftWidth(content.width);
    setLeftKey(key);
    setIsLeftOpen(true);
  }, []);

  const closeLeft = useCallback(() => {
    setLeftNodeState(<div/>);
    setLeftWidth(0);
    setLeftKey("");
    setIsLeftOpen(false);
  }, []);

  const openRight = useCallback((currentKey: string) => {
    const key = `right:${currentKey ?? ""}`;
    const content = overlay.get(key) as OverlayOpts;
    setRightNodeState(content.render());
    setRightWidth(content.width);
    setRightKey(key);
    setIsRightOpen(true);
  }, []);

  const closeRight = useCallback(() => {
    setRightNodeState(<div/>);
    setRightWidth(0);
    setRightKey("");
    setIsRightOpen(false);
  }, []);

  const value = useMemo<OverlayActions>(() => ({
    // left
    isLeftOpen,
    leftKey,
    leftNode,
    leftWidth,
    openLeft,
    closeLeft,
    // right
    isRightOpen,
    rightKey,
    rightNode,
    rightWidth,
    openRight,
    closeRight,
  }), [isLeftOpen, leftKey, leftNode, leftWidth, closeLeft, openLeft, isRightOpen, rightKey, rightNode, rightWidth, openRight, closeRight]);

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
