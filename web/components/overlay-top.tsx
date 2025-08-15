"use client";

import { useOverlay } from "@/hooks/useOverlay";
import ButtonIcon from "@/components/button-icon";
import MenuIcon from "@/components/icons/menu";
import SearchIcon from "@/components/icons/search";
import BellIcon from "@/components/icons/bell";
import Avatar from "@/components/avatar";
import ProfileOverlay from "@/components/reuse/profile-overlay";
import NotificationOverlay from "@/components/reuse/notification-overlay";
import NavigationOverlay from "@/components/reuse/navigation-overlay";
import clsx from "clsx";
import type { ReactNode } from "react";


type Side = "left" | "right" | "modal";

function OverlayButtonIcon(props: {
  ariaLabel?: string;
  side: Side;
  overlayKey: string;
  width?: number;
  node?: ReactNode;
  children: ReactNode;
  type?: "button" | "submit" | "reset";
}) {
  const {
    setLeftWidth,
    toggleLeft,

    setRightWidth,
    toggleRight,

    setModalNode,
    isModalOpen,
    modalActiveKey,
    openModal,
    closeModal,
    setModalActiveKey,
  } = useOverlay();

  const {
    ariaLabel,
    side,
    overlayKey,
    width,
    node,
    children,
    type = "button",
  } = props;

  function handleClick() {
    if (side === "left") {
      if (typeof width === "number") setLeftWidth(width);
      toggleLeft(overlayKey, node);
      return;
    }

    if (side === "right") {
      if (typeof width === "number") setRightWidth(width);
      toggleRight(overlayKey, node);
      return;
    }

    if (side === "modal") {
      if (isModalOpen && modalActiveKey === overlayKey) {
        setModalActiveKey(null);
        closeModal();
        return;
      }
      if (node) setModalNode(node);
      setModalActiveKey(overlayKey);
      openModal();
    }
  }

  return (
    <ButtonIcon type={type} aria-label={ariaLabel} onClick={handleClick}>
      {children}
    </ButtonIcon>
  );
}

export default function OverlayTop() {
  const {
    setRightWidth,
    setModalNode,
    isModalOpen,
    modalActiveKey,
    openModal,
    closeModal,
    toggleRight,
    setModalActiveKey,
  } = useOverlay();

  function openSearchModal() {
    const key = 'modal:search';

    if (isModalOpen && modalActiveKey === key) {
      setModalActiveKey(null);
      closeModal();
      return;
    }

    setModalNode(
      <div className="space-y-3">
        <div className="relative">
          <label htmlFor="global-search-input" className="sr-only">Search</label>
          <input
            id="global-search-input"
            type="text"
            autoFocus
            placeholder="Search..."
            className="w-full py-2 pl-3 pr-10 rounded-xl bg-white text-gray-800 placeholder:text-gray-400 border border-gray-200 focus:border-gray-400 focus:outline-none"
          />
          <SearchIcon
            className="pointer-events-none absolute right-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-500"
            aria-hidden="true"
          />
        </div>
        <div className="text-sm text-gray-500">Type to search…</div>
      </div>
    );
    setModalActiveKey(key);
    openModal();
  }

  return (
    <div className="absolute top-0 left-0 right-0 w-full p-2">
      <div className="h-[60px] w-full flex flex-row justify-between items-center px-3">
        <div className="flex items-center space-x-3">
          <OverlayButtonIcon
            side="left"
            overlayKey="left:menu"
            width={88}
            node={<NavigationOverlay />}
          >
            <MenuIcon className="w-6 h-6 text-gray-800" />
          </OverlayButtonIcon>
          <h1 className="text-xl">LOGO</h1>
        </div>
        <div className="flex items-center space-x-3 flex-1 justify-center">
          <div className="relative w-96 max-w-full">
            <label htmlFor="topbar-search" className="sr-only">Search</label>
            <input
              id="topbar-search"
              type="text"
              placeholder="Search..."
              className={clsx(
                "w-full py-3 pl-5 pr-10 rounded-full",
                "bg-white shadow-primary placeholder:text-gray-600 focus:outline-none",
                "cursor-pointer"
              )}
              onFocus={openSearchModal}
              onClick={openSearchModal}
              readOnly
            />
            <SearchIcon
              className="pointer-events-none absolute right-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-800"
              aria-hidden="true"
            />
          </div>
        </div>
        <div className="flex items-center space-x-3">
          <OverlayButtonIcon
            type="button"
            ariaLabel="Show notifications in right overlay"
            side="right"
            overlayKey="right:notifications"
            node={<NotificationOverlay />}
          >
            <BellIcon className="w-6 h-6 text-gray-800" />
          </OverlayButtonIcon>
          <button
            type="button"
            className="rounded-full focus:outline-none cursor-pointer"
            onClick={() => {
              const key = 'right:profile';

              setRightWidth(400);
              toggleRight(key, <ProfileOverlay />);
            }}
          >
            <Avatar name="Artha Dede" size={40} />
          </button>
        </div>
      </div>
    </div>
  );
}
