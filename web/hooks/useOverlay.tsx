"use client";

import React, { createContext, useCallback, useContext, useEffect, useMemo, useState } from "react";

export type OverlayState = {
  isLeftOpen: boolean;
  isRightOpen: boolean;
  leftWidth: number; // in px
  rightWidth: number; // in px
};

export type OverlayActions = {
  openLeft: () => void;
  closeLeft: () => void;
  toggleLeft: () => void;
  setLeftWidth: (px: number) => void;

  openRight: () => void;
  closeRight: () => void;
  toggleRight: () => void;
  setRightWidth: (px: number) => void;
};

export type OverlayContextType = OverlayState & OverlayActions;

const OverlayContext = createContext<OverlayContextType | undefined>(undefined);

export function OverlayProvider({
  children,
  defaultLeftOpen = true,
  defaultRightOpen = true,
  defaultLeftWidth = 300,
  defaultRightWidth = 300,
}: {
  children: React.ReactNode;
  defaultLeftOpen?: boolean;
  defaultRightOpen?: boolean;
  defaultLeftWidth?: number;
  defaultRightWidth?: number;
}) {
  const [isLeftOpen, setIsLeftOpen] = useState<boolean>(defaultLeftOpen);
  const [isRightOpen, setIsRightOpen] = useState<boolean>(defaultRightOpen);
  const [leftWidth, setLeftWidthState] = useState<number>(defaultLeftWidth);
  const [rightWidth, setRightWidthState] = useState<number>(defaultRightWidth);

  // Clamp helper
  const clamp = (v: number, min: number, max: number) => Math.max(min, Math.min(max, v));

  // Actions
  const openLeft = useCallback(() => setIsLeftOpen(true), []);
  const closeLeft = useCallback(() => setIsLeftOpen(false), []);
  const toggleLeft = useCallback(() => setIsLeftOpen(v => !v), []);
  const setLeftWidth = useCallback((px: number) => setLeftWidthState(prev => (Number.isFinite(px) ? clamp(Math.round(px), 0, 2000) : prev)), []);

  const openRight = useCallback(() => setIsRightOpen(true), []);
  const closeRight = useCallback(() => setIsRightOpen(false), []);
  const toggleRight = useCallback(() => setIsRightOpen(v => !v), []);
  const setRightWidth = useCallback((px: number) => setRightWidthState(prev => (Number.isFinite(px) ? clamp(Math.round(px), 0, 2000) : prev)), []);

  // Reflect state into CSS variables so existing components that rely on them continue to work
  useEffect(() => {
    if (typeof document === 'undefined') return;
    const root = document.documentElement;
    root.style.setProperty('--overlay-left-width', isLeftOpen ? `${leftWidth}px` : '0px');
    root.style.setProperty('--overlay-right-width', isRightOpen ? `${rightWidth}px` : '0px');
  }, [isLeftOpen, isRightOpen, leftWidth, rightWidth]);

  const value = useMemo<OverlayContextType>(() => ({
    isLeftOpen,
    isRightOpen,
    leftWidth,
    rightWidth,
    openLeft,
    closeLeft,
    toggleLeft,
    setLeftWidth,
    openRight,
    closeRight,
    toggleRight,
    setRightWidth,
  }), [isLeftOpen, isRightOpen, leftWidth, rightWidth, openLeft, closeLeft, toggleLeft, setLeftWidth, openRight, closeRight, toggleRight, setRightWidth]);

  return (
    <OverlayContext.Provider value={value}>
      {children}
    </OverlayContext.Provider>
  );
}

export function useOverlay(): OverlayContextType {
  const ctx = useContext(OverlayContext);
  if (!ctx) throw new Error('useOverlay must be used within an OverlayProvider');
  return ctx;
}
