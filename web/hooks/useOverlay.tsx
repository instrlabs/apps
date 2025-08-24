"use client";

import React, { createContext, useCallback, useContext, useEffect, useMemo, useState } from "react";import { getOverlayEntry, resolveOverlayNode } from "@/hooks/overlayRegistry";

export type OverlayActions = {
  // actions
  toggleByKey: (key: string) => void;
  closeAll: () => void;
  // left overlay state
  isLeftOpen: boolean;
  leftNode: React.ReactNode;
  leftContentKey: string;
  leftWidth?: number;
  // right overlay state
  isRightOpen: boolean;
  rightNode: React.ReactNode;
  rightContentKey: string;
  rightWidth?: number;
  // modal state
  isModalOpen: boolean;
  modalNode: React.ReactNode;
  modalContentKey: string;
};

const OverlayContext = createContext<OverlayActions | undefined>(undefined);

export function OverlayProvider({ children }: { children: React.ReactNode }) {
  // left
  const [isLeftOpen, setIsLeftOpen] = useState<boolean>(false);
  const [leftNode, setLeftNodeState] = useState<React.ReactNode>(<div />);
  const [leftKey, setLeftKey] = useState<string>("");
  const [leftWidth, setLeftWidth] = useState<number | undefined>(undefined);

  // right
  const [isRightOpen, setIsRightOpen] = useState<boolean>(false);
  const [rightNode, setRightNodeState] = useState<React.ReactNode>(<div />);
  const [rightActiveKey, setRightActiveKey] = useState<string>("");
  const [rightWidth, setRightWidth] = useState<number | undefined>(undefined);

  // modal
  const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
  const [modalNode, setModalNodeState] = useState<React.ReactNode>(<div />);
  const [modalActiveKey, setModalActiveKey] = useState<string>("");

  const toggleLeft = useCallback((currentKey: string | null, nextNode?: React.ReactNode, width?: number) => {
    const key = currentKey ?? "";
    if (isLeftOpen && leftKey === key) {
      setIsLeftOpen(false);
      setLeftKey("");
      return;
    }

    if (nextNode) setLeftNodeState(nextNode);
    setLeftKey(key);
    setLeftWidth(width);
    setIsLeftOpen(true);
  }, [isLeftOpen, leftKey]);

  const toggleRight = useCallback((currentKey: string | null, nextNode?: React.ReactNode, width?: number) => {
    const key = currentKey ?? "";
    if (isRightOpen && rightActiveKey === key) {
      setIsRightOpen(false);
      setRightActiveKey("");
      return;
    }
    if (nextNode) setRightNodeState(nextNode);
    setRightActiveKey(key);
    setRightWidth(width);
    setIsRightOpen(true);
  }, [isRightOpen, rightActiveKey]);

  const toggleModal = useCallback((currentKey: string | null, nextNode?: React.ReactNode) => {
    const key = currentKey ?? "";
    if (isModalOpen && modalActiveKey === key) {
      setIsModalOpen(false);
      setModalActiveKey("");
      return;
    }

    if (nextNode) setModalNodeState(nextNode);
    setModalActiveKey(key);
    setIsModalOpen(true);
  }, [isModalOpen, modalActiveKey]);

  const closeAll = useCallback(() => {
    setIsLeftOpen(false);
    setIsRightOpen(false);
    setIsModalOpen(false);
  }, []);

  const toggleByKey = useCallback((key: string) => {
    const entry = getOverlayEntry(key);
    if (!entry) return;

    const node = resolveOverlayNode(entry);

    if (entry.side === 'left') {
      toggleLeft(key, node, entry.width);
      return;
    }

    if (entry.side === 'right') {
      toggleRight(key, node, entry.width);
      return;
    }

    if(entry.side === 'modal') {
      toggleModal(key, node);
    }
  }, [toggleLeft, toggleRight, toggleModal]);

  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        if (isModalOpen || isLeftOpen || isRightOpen) {
          e.stopPropagation();
          closeAll();
        }
      }
    };
    window.addEventListener('keydown', onKey);
    return () => window.removeEventListener('keydown', onKey);
  }, [isModalOpen, isLeftOpen, isRightOpen, closeAll]);

  const value = useMemo<OverlayActions>(() => ({
    // actions
    toggleByKey,
    closeAll,
    // left
    isLeftOpen,
    leftNode,
    leftContentKey: leftKey,
    leftWidth,
    // right
    isRightOpen: isRightOpen,
    rightNode,
    rightContentKey: rightActiveKey,
    rightWidth,
    // modal
    isModalOpen,
    modalNode,
    modalContentKey: modalActiveKey,
  }), [toggleByKey, closeAll, isLeftOpen, leftNode, leftKey, leftWidth, isRightOpen, rightNode, rightActiveKey, rightWidth, isModalOpen, modalNode, modalActiveKey]);

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
