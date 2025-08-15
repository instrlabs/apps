"use client";

import React, { createContext, useCallback, useContext, useEffect, useMemo, useState } from "react";

export type OverlayState = {
  isLeftOpen: boolean;
  leftWidth: number;
  leftNode: React.ReactNode | null;
  leftContentKey: number;
  leftActiveKey?: string | null;

  isRightOpen: boolean;
  rightContentKey: number;
  rightWidth: number;
  rightNode: React.ReactNode | null;

  isModalOpen: boolean;
  modalNode: React.ReactNode | null;
  modalContentKey: number;
  modalActiveKey?: string | null;
};

export type OverlayActions = {
  openLeft: () => void;
  closeLeft: () => void;
  setLeftWidth: (px: number) => void;
  setLeftNode: (node: React.ReactNode) => void;
  setLeftActiveKey: (key: string | null) => void;

  toggleRight: (currentKey: string | null, nextNode?: React.ReactNode) => void;
  setRightWidth: (px: number) => void;
  setRightNode: (node: React.ReactNode) => void;

  openModal: () => void;
  closeModal: () => void;
  setModalNode: (node: React.ReactNode) => void;
  setModalActiveKey: (key: string | null) => void;
};

export type OverlayContextType = OverlayState & OverlayActions;

const OverlayContext = createContext<OverlayContextType | undefined>(undefined);

export function OverlayProvider({
  children,
  defaultLeftOpen = false,
  defaultRightOpen = false,
  defaultLeftWidth = 0,
  defaultRightWidth = 0,
}: {
  children: React.ReactNode;
  defaultLeftOpen?: boolean;
  defaultRightOpen?: boolean;
  defaultLeftWidth?: number;
  defaultRightWidth?: number;
}) {
  const [isLeftOpen, setIsLeftOpen] = useState<boolean>(defaultLeftOpen);
  const [leftWidth, setLeftWidthState] = useState<number>(defaultLeftWidth);
  const [leftNode, setLeftNodeState] = useState<React.ReactNode | null>(null);
  const [leftContentKey, setLeftContentKey] = useState<number>(0);
  const [lastLeftNode, setLastLeftNode] = useState<React.ReactNode | null>(null);
  const [leftActiveKey, setLeftActiveKey] = useState<string | null>(null);

  const [isRightOpen, setIsRightOpen] = useState<boolean>(defaultRightOpen);
  const [rightWidth, setRightWidthState] = useState<number>(defaultRightWidth);
  const [rightNode, setRightNodeState] = useState<React.ReactNode | null>(null);
  const [rightContentKey, setRightContentKey] = useState<number>(0);
  const [lastRightNode, setLastRightNode] = useState<React.ReactNode | null>(null);
  const [rightActiveKey, setRightActiveKey] = useState<string | null>(null);

  const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
  const [modalNode, setModalNodeState] = useState<React.ReactNode | null>(null);
  const [modalContentKey, setModalContentKey] = useState<number>(0);
  const [lastModalNode, setLastModalNode] = useState<React.ReactNode | null>(null);
  const [modalActiveKey, setModalActiveKey] = useState<string | null>(null);

  const openLeft = useCallback(() => {
    setIsLeftOpen(true);
    // restore last content if none present
    if (leftNode == null && lastLeftNode != null) {
      setLeftNodeState(lastLeftNode);
    }
  }, [leftNode, lastLeftNode]);
  const closeLeft = useCallback(() => {
    setLeftNodeState(null);
    setIsLeftOpen(false);
    setLeftActiveKey(null);
  }, []);
  const setLeftWidth = useCallback((px: number) => {
    setLeftWidthState(prev => (Number.isFinite(px) ? Math.round(px) : prev));
  }, []);


  const toggleRight = useCallback((currentKey: string | null, nextNode?: React.ReactNode) => {
    const lastKey = rightActiveKey;

    if (
      isRightOpen &&
      lastKey != null &&
      currentKey === lastKey
    ) {
      setIsRightOpen(false);
      return;
    }

    // Switching to a different key (or first time)
    if (nextNode !== undefined) {
      setRightNodeState(nextNode);
      setLastRightNode(nextNode);
      setRightContentKey(k => k + 1);
    } else {
      // If no node provided and none currently, try to restore the last one
      if (rightNode == null && lastRightNode != null) {
        setRightNodeState(lastRightNode);
      }
    }

    setRightActiveKey(currentKey);

    if (!isRightOpen) {
      // Case 3: was closed -> open it
      setIsRightOpen(true);
    }
    // Case 2: already open -> remain open
  }, [isRightOpen, rightNode, lastRightNode, rightActiveKey]);
  const setRightWidth = useCallback((px: number) => {
    setRightWidthState(prev => (Number.isFinite(px) ? Math.round(px) : prev))
  }, []);

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
  const setModalNode = useCallback((node: React.ReactNode) => {
    setModalNodeState(node);
    setLastModalNode(node);
    setModalContentKey(k => k + 1);
  }, []);

  const value = useMemo<OverlayContextType>(() => ({
    // Left overlay
    isLeftOpen,
    leftWidth,
    leftNode,
    leftContentKey,
    leftActiveKey,
    openLeft,
    closeLeft,
    setLeftWidth,
    setLeftNode,
    setLeftActiveKey,

    // Right overlay
    isRightOpen,
    rightWidth,
    rightNode,
    rightContentKey,
    toggleRight,
    setRightWidth,
    setRightNode,

    // Modal overlay
    isModalOpen,
    modalNode,
    modalContentKey,
    modalActiveKey,
    openModal,
    closeModal,
    setModalNode,
    setModalActiveKey,
  }), [
    // Left overlay
    isLeftOpen,
    leftWidth,
    leftNode,
    leftContentKey,
    leftActiveKey,
    openLeft,
    closeLeft,
    setLeftWidth,
    setLeftNode,
    setLeftActiveKey,

    // Right overlay
    isRightOpen,
    rightWidth,
    rightNode,
    rightContentKey,
    toggleRight,
    setRightWidth,
    setRightNode,

    // Modal overlay
    isModalOpen,
    modalNode,
    modalContentKey,
    modalActiveKey,
    openModal,
    closeModal,
    setModalNode,
    setModalActiveKey
  ]);

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
