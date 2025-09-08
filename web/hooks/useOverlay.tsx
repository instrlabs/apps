"use client";

import React, {createContext, useCallback, useContext, useEffect, useMemo, useState} from "react";

import NavigationOverlay from "@/components/overlays/navigation-overlay";
import NotificationOverlay from "@/components/overlays/notification-overlay";
import ProfileOverlay from "@/components/overlays/profile-overlay";
import SearchOverlay from "@/components/overlays/search-overlay";

export type OverlayActions = {
  // actions
  toggleByKey: (key: string) => void;
  closeAll: () => void;
  // left overlay state
  isLeftOpen: boolean;
  leftKey: string;
  leftNode: React.ReactNode;
  leftWidth: number;
  // right overlay state
  isRightOpen: boolean;
  rightKey: string;
  rightNode: React.ReactNode;
  rightWidth: number;
  // modal state
  isModalOpen: boolean;
  modalKey: string;
  modalNode: React.ReactNode;
  modalWidth: number;
};

export type OverlaySide = "left" | "right" | "modal";
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

  // Modal: search
  registerOverlay("modal:search", {
    side: "modal",
    width: 800,
    render: () => <SearchOverlay />,
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

  // modal
  const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
  const [modalNode, setModalNodeState] = useState<React.ReactNode>(<div />);
  const [modalKey, setModalKey] = useState<string>("");
  const [modalWidth, setModalWidth] = useState<number>(0);

  const toggleLeft = useCallback((currentKey: string | null, params: OverlayOpts) => {
    const key = currentKey ?? "";
    if (isLeftOpen && leftKey === key) {
      setIsLeftOpen(false);
      setLeftKey("");
      return;
    }

    setLeftNodeState(params.render());
    setLeftKey(key);
    setLeftWidth(params.width);
    setIsLeftOpen(true);
  }, [isLeftOpen, leftKey]);

  const toggleRight = useCallback((currentKey: string | null, params: OverlayOpts) => {
    const key = currentKey ?? "";
    if (isRightOpen && rightKey === key) {
      setIsRightOpen(false);
      setRightKey("");
      return;
    }

    setRightNodeState(params.render());
    setRightKey(key);
    setRightWidth(params.width);
    setIsRightOpen(true);
  }, [isRightOpen, rightKey]);

  const toggleModal = useCallback((currentKey: string | null, params: OverlayOpts) => {
    const key = currentKey ?? "";
    if (isModalOpen && modalKey === key) {
      setIsModalOpen(false);
      setModalKey("");
      return;
    }

    setModalNodeState(params.render());
    setModalKey(key);
    setModalWidth(params.width);
    setIsModalOpen(true);
  }, [isModalOpen, modalKey]);

  const closeAll = useCallback(() => {
    setIsLeftOpen(false);
    setIsRightOpen(false);
    setIsModalOpen(false);
  }, []);

  const toggleByKey = useCallback((key: string) => {
    const entry = overlay.get(key);
    if (!entry) return;

    switch (entry.side) {
      case 'left':
        toggleLeft(key, entry);
        return;
      case 'right':
        toggleRight(key, entry);
        return;
      case 'modal':
        toggleModal(key, entry);
        return;
      default:
        return;
    }
  }, [toggleLeft, toggleRight, toggleModal]);

  const value = useMemo<OverlayActions>(() => ({
    // actions
    toggleByKey,
    closeAll,
    // left
    isLeftOpen,
    leftNode,
    leftKey,
    leftWidth,
    // right
    isRightOpen,
    rightNode,
    rightKey,
    rightWidth,
    // modal
    isModalOpen,
    modalNode,
    modalKey,
    modalWidth,
  }), [toggleByKey, closeAll, isLeftOpen, leftNode, leftKey, leftWidth, isRightOpen, rightNode, rightKey, rightWidth, isModalOpen, modalNode, modalKey, modalWidth]);

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
