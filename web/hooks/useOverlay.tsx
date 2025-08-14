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
  // new: arbitrary React content for the left overlay, and a key to trigger animations
  leftNode: React.ReactNode | null;
  leftContentKey: number;
  // right overlay content
  rightNode: React.ReactNode | null;
  rightContentKey: number;
  // modal (center) overlay content
  isModalOpen: boolean;
  modalNode: React.ReactNode | null;
  modalContentKey: number;
  // active keys for toggling by identity
  leftActiveKey?: string | null;
  rightActiveKey?: string | null;
  modalActiveKey?: string | null;
};

export type OverlayActions = {
  openLeft: () => void;
  closeLeft: () => void;
  toggleLeft: () => void;
  toggleLeftByKey: (key: string) => void;
  setLeftWidth: (px: number) => void;
  // legacy: still allow setting a type label
  setLeftContent: (c: LeftOverlayContent) => void;
  // new: set arbitrary content
  setLeftNode: (node: React.ReactNode) => void;

  openRight: () => void;
  closeRight: () => void;
  toggleRight: () => void;
  toggleRightByKey: (key: string) => void;
  setRightWidth: (px: number) => void;
  // right side content setters
  setRightNode: (node: React.ReactNode) => void;

  // modal (center) overlay actions
  openModal: () => void;
  closeModal: () => void;
  toggleModalByKey: (key: string) => void;
  setModalNode: (node: React.ReactNode) => void;
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
  const [leftContentKey, setLeftContentKey] = useState<number>(0);
  const [rightNode, setRightNodeState] = useState<React.ReactNode | null>(null);
  const [rightContentKey, setRightContentKey] = useState<number>(0);
  const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
  const [modalNode, setModalNodeState] = useState<React.ReactNode | null>(null);
  const [modalContentKey, setModalContentKey] = useState<number>(0);

  // last-shown content caches (restored on open if no new content is set)
  const [lastLeftNode, setLastLeftNode] = useState<React.ReactNode | null>(null);
  const [lastRightNode, setLastRightNode] = useState<React.ReactNode | null>(null);
  const [lastModalNode, setLastModalNode] = useState<React.ReactNode | null>(null);

  const [leftActiveKey, setLeftActiveKey] = useState<string | null>(null);
  const [rightActiveKey, setRightActiveKey] = useState<string | null>(null);
  const [modalActiveKey, setModalActiveKey] = useState<string | null>(null);

  // Clamp helper
  const clamp = (v: number, min: number, max: number) => Math.max(min, Math.min(max, v));

  // Actions
  const openLeft = useCallback(() => {
    setIsLeftOpen(true);
    // restore last content if none present
    if (leftNode == null && lastLeftNode != null) {
      setLeftNodeState(lastLeftNode);
    }
  }, [leftNode, lastLeftNode]);
  const closeLeft = useCallback(() => {
    // cleanup content before closing
    setLeftNodeState(null);
    setIsLeftOpen(false);
    setLeftActiveKey(null);
  }, []);
  const toggleLeft = useCallback(() => setIsLeftOpen(v => !v), []);
  const toggleLeftByKey = useCallback((key: string) => {
    if (isLeftOpen && leftActiveKey === key) {
      // cleanup content before closing
      setLeftNodeState(null);
      setIsLeftOpen(false);
      setLeftActiveKey(null);
    } else {
      setIsLeftOpen(true);
      setLeftActiveKey(key);
      // restore last content if none present
      if (leftNode == null && lastLeftNode != null) {
        setLeftNodeState(lastLeftNode);
      }
    }
  }, [isLeftOpen, leftActiveKey, leftNode, lastLeftNode]);
  const setLeftWidth = useCallback((px: number) => setLeftWidthState(prev => (Number.isFinite(px) ? clamp(Math.round(px), 0, 2000) : prev)), []);

  const openRight = useCallback(() => {
    setIsRightOpen(true);
    // restore last content if none present
    if (rightNode == null && lastRightNode != null) {
      setRightNodeState(lastRightNode);
    }
  }, [rightNode, lastRightNode]);
  const closeRight = useCallback(() => {
    // cleanup content before closing
    setRightNodeState(null);
    setIsRightOpen(false);
    setRightActiveKey(null);
  }, []);
  const toggleRight = useCallback(() => setIsRightOpen(v => !v), []);
  const toggleRightByKey = useCallback((key: string) => {
    if (isRightOpen && rightActiveKey === key) {
      // cleanup content before closing
      setRightNodeState(null);
      setIsRightOpen(false);
      setRightActiveKey(null);
    } else {
      setIsRightOpen(true);
      setRightActiveKey(key);
      // restore last content if none present
      if (rightNode == null && lastRightNode != null) {
        setRightNodeState(lastRightNode);
      }
    }
  }, [isRightOpen, rightActiveKey, rightNode, lastRightNode]);
  const setRightWidth = useCallback((px: number) => setRightWidthState(prev => (Number.isFinite(px) ? clamp(Math.round(px), 0, 2000) : prev)), []);
  const setLeftContent = useCallback((c: LeftOverlayContent) => setLeftContentState(c), []);

  const setLeftNode = useCallback((node: React.ReactNode) => {
    setLeftNodeState(node);
    setLastLeftNode(node);
    setLeftContentKey(k => k + 1);
  }, []);
  const setRightNode = useCallback((node: React.ReactNode) => {
    setRightNodeState(node);
    setLastRightNode(node);
    setRightContentKey(k => k + 1);
  }, []);
  // modal actions
  const openModal = useCallback(() => {
    setIsModalOpen(true);
    // restore last content if none present
    if (modalNode == null && lastModalNode != null) {
      setModalNodeState(lastModalNode);
    }
  }, [modalNode, lastModalNode]);
  const closeModal = useCallback(() => {
    // cleanup content before closing
    setModalNodeState(null);
    setIsModalOpen(false);
    setModalActiveKey(null);
  }, []);
  const toggleModalByKey = useCallback((key: string) => {
    if (isModalOpen && modalActiveKey === key) {
      // cleanup content before closing
      setModalNodeState(null);
      setIsModalOpen(false);
      setModalActiveKey(null);
    } else {
      setIsModalOpen(true);
      setModalActiveKey(key);
      // restore last content if none present
      if (modalNode == null && lastModalNode != null) {
        setModalNodeState(lastModalNode);
      }
    }
  }, [isModalOpen, modalActiveKey, modalNode, lastModalNode]);
  const setModalNode = useCallback((node: React.ReactNode) => {
    setModalNodeState(node);
    setLastModalNode(node);
    setModalContentKey(k => k + 1);
  }, []);
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
    leftContentKey,
    rightNode,
    rightContentKey,
    // modal
    isModalOpen,
    modalNode,
    modalContentKey,
    // active keys
    leftActiveKey,
    rightActiveKey,
    modalActiveKey,
    openLeft,
    closeLeft,
    toggleLeft,
    toggleLeftByKey,
    setLeftWidth,
    setLeftContent,
    setLeftNode,
    openRight,
    closeRight,
    toggleRight,
    toggleRightByKey,
    setRightWidth,
    setRightNode,
    // modal actions
    openModal,
    closeModal,
    toggleModalByKey,
    setModalNode,
  }), [isLeftOpen, isRightOpen, leftWidth, rightWidth, leftContent, leftNode, leftContentKey, rightNode, rightContentKey, isModalOpen, modalNode, modalContentKey, leftActiveKey, rightActiveKey, modalActiveKey, openLeft, closeLeft, toggleLeft, toggleLeftByKey, setLeftWidth, setLeftContent, setLeftNode, openRight, closeRight, toggleRight, toggleRightByKey, setRightWidth, setRightNode, openModal, closeModal, toggleModalByKey, setModalNode]);

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
