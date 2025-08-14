"use client";

import { useOverlay } from "@/hooks/useOverlay";
import ButtonIcon from "@/components/button-icon";
import MenuIcon from "@/components/icons/menu";
import SearchIcon from "@/components/icons/search";
import BellIcon from "@/components/icons/bell";
import Avatar from "@/components/avatar";

export default function OverlayTop() {
  const {
    isLeftOpen,
    setLeftTitle,
    setLeftNode,
    openLeft,
    setRightTitle,
    setRightNode,
    openRight,
    // modal actions
    setModalTitle,
    setModalNode,
    openModal,
  } = useOverlay();

  function openSearchModal() {
    setModalTitle('Search');
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
    openModal();
  }

  return (
    <div className="absolute top-0 left-0 right-0 w-full p-2">
      <div className="h-[60px] w-full rounded-full bg-neutral-50 flex flex-row justify-between items-center px-3">
        <div className="flex items-center space-x-3">
          <ButtonIcon
            type="button"
            aria-label={isLeftOpen ? "Hide left" : "Show left (menu)"}
            onClick={() => {
              setLeftTitle('Menu');
              setLeftNode(
                <ul className="space-y-2 text-sm text-gray-700">
                  <li><button className="w-full text-left px-3 py-2 rounded hover:bg-gray-100">Dashboard</button></li>
                  <li><button className="w-full text-left px-3 py-2 rounded hover:bg-gray-100">Projects</button></li>
                  <li><button className="w-full text-left px-3 py-2 rounded hover:bg-gray-100">Teams</button></li>
                  <li><button className="w-full text-left px-3 py-2 rounded hover:bg-gray-100">Settings</button></li>
                </ul>
              );
              openLeft();
            }}
          >
            <MenuIcon className="w-6 h-6 text-gray-800" />
          </ButtonIcon>
          <h1 className="text-xl">LOGO</h1>
        </div>
        <div className="flex items-center space-x-3 flex-1 justify-center">
          <div className="relative w-72 max-w-full">
            <label htmlFor="topbar-search" className="sr-only">Search</label>
            <input
              id="topbar-search"
              type="text"
              placeholder="Search..."
              className="w-full py-1.5 pl-3 pr-10 rounded-full bg-white text-gray-800 placeholder:text-gray-400 hover:border-gray-400 focus:outline-none"
              onFocus={openSearchModal}
              onClick={openSearchModal}
              readOnly
            />
            <SearchIcon
              className="pointer-events-none absolute right-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-500"
              aria-hidden="true"
            />
          </div>
        </div>
        <div className="flex items-center space-x-3">
          <ButtonIcon
            type="button"
            aria-label="Show notifications in right overlay"
            onClick={() => {
              setRightTitle('Notification Overlay');
              setRightNode(
                <div className="text-sm text-gray-700">Notification Overlay</div>
              );
              openRight();
            }}
          >
            <BellIcon className="w-6 h-6 text-gray-800" />
          </ButtonIcon>
          <button
            type="button"
            aria-label="Profile"
            className="rounded-full focus:outline-none"
            onClick={() => {
              setRightTitle('Profile Overlay');
              setRightNode(
                <div className="text-sm text-gray-700">Profile Overlay</div>
              );
              openRight();
            }}
          >
            <Avatar name="Artha Dede" size={40} />
          </button>
        </div>
      </div>
    </div>
  );
}
