"use client";

import React, { createContext, useCallback, useContext, useEffect, useMemo, useState } from "react";

export type LeftOverlayContent = "menu" | "notifications"; // kept for backward compatibility

export type OverlayState = {
  isLeftOpen: boolean;
  isRightOpen: boolean;
  leftWidth: number; // in px
  rightWidth: number; // in px
  // legacy flag of what type of content we intend to show (optional)
  leftContent: LeftOverlayContent;
  // new: arbitrary React content and title for the left overlay, and a key to trigger animations
  leftNode: React.ReactNode | null;
  leftTitle: string;
  leftContentKey: number;
  // right overlay content
  rightNode: React.ReactNode | null;
  rightTitle: string;
  rightContentKey: number;
  // modal (center) overlay content
  isModalOpen: boolean;
  modalNode: React.ReactNode | null;
  modalTitle: string;
  modalContentKey: number;
};

export type OverlayActions = {
  openLeft: () => void;
  closeLeft: () => void;
  toggleLeft: () => void;
  setLeftWidth: (px: number) => void;
  // legacy: still allow setting a type label
  setLeftContent: (c: LeftOverlayContent) => void;
  // new: set arbitrary content and title
  setLeftNode: (node: React.ReactNode) => void;
  setLeftTitle: (title: string) => void;

  openRight: () => void;
  closeRight: () => void;
  toggleRight: () => void;
  setRightWidth: (px: number) => void;
  // right side content setters
  setRightNode: (node: React.ReactNode) => void;
  setRightTitle: (title: string) => void;

  // modal (center) overlay actions
  openModal: () => void;
  closeModal: () => void;
  setModalNode: (node: React.ReactNode) => void;
  setModalTitle: (title: string) => void;
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
  const [leftContent, setLeftContentState] = useState<LeftOverlayContent>("menu");
  const [leftNode, setLeftNodeState] = useState<React.ReactNode | null>(null);
  const [leftTitle, setLeftTitleState] = useState<string>("");
  const [leftContentKey, setLeftContentKey] = useState<number>(0);
  const [rightNode, setRightNodeState] = useState<React.ReactNode | null>(null);
  const [rightTitle, setRightTitleState] = useState<string>("");
  const [rightContentKey, setRightContentKey] = useState<number>(0);
  // modal state
  const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
  const [modalNode, setModalNodeState] = useState<React.ReactNode | null>(null);
  const [modalTitle, setModalTitleState] = useState<string>("");
  const [modalContentKey, setModalContentKey] = useState<number>(0);

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
  const setLeftContent = useCallback((c: LeftOverlayContent) => setLeftContentState(c), []);

  const setLeftNode = useCallback((node: React.ReactNode) => {
    setLeftNodeState(node);
    setLeftContentKey(k => k + 1);
  }, []);
  const setLeftTitle = useCallback((title: string) => setLeftTitleState(title), []);

  const setRightNode = useCallback((node: React.ReactNode) => {
    setRightNodeState(node);
    setRightContentKey(k => k + 1);
  }, []);
  const setRightTitle = useCallback((title: string) => setRightTitleState(title), []);

  // modal actions
  const openModal = useCallback(() => setIsModalOpen(true), []);
  const closeModal = useCallback(() => setIsModalOpen(false), []);
  const setModalNode = useCallback((node: React.ReactNode) => {
    setModalNodeState(node);
    setModalContentKey(k => k + 1);
  }, []);
  const setModalTitle = useCallback((title: string) => setModalTitleState(title), []);

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
    leftContent,
    leftNode,
    leftTitle,
    leftContentKey,
    rightNode,
    rightTitle,
    rightContentKey,
    // modal
    isModalOpen,
    modalNode,
    modalTitle,
    modalContentKey,
    openLeft,
    closeLeft,
    toggleLeft,
    setLeftWidth,
    setLeftContent,
    setLeftNode,
    setLeftTitle,
    openRight,
    closeRight,
    toggleRight,
    setRightWidth,
    setRightNode,
    setRightTitle,
    // modal actions
    openModal,
    closeModal,
    setModalNode,
    setModalTitle,
  }), [isLeftOpen, isRightOpen, leftWidth, rightWidth, leftContent, leftNode, leftTitle, leftContentKey, rightNode, rightTitle, rightContentKey, isModalOpen, modalNode, modalTitle, modalContentKey, openLeft, closeLeft, toggleLeft, setLeftWidth, setLeftContent, setLeftNode, setLeftTitle, openRight, closeRight, toggleRight, setRightWidth, setRightNode, setRightTitle, openModal, closeModal, setModalNode, setModalTitle]);

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
