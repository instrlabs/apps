"use client";

import { useOverlay } from "@/hooks/useOverlay";
import ButtonIcon from "@/components/button-icon";
import MenuIcon from "@/components/icons/menu";
import SearchIcon from "@/components/icons/search";
import BellIcon from "@/components/icons/bell";

export default function OverlayTop() {
  const { isLeftOpen, toggleLeft } = useOverlay();
  return (
    <div className="absolute top-0 left-0 right-0 w-full p-3">
      <div className="h-[60px] rounded-2xl bg-gray-200 flex flex-row justify-between items-center px-3">
        <div className="flex items-center space-x-3">
          <ButtonIcon
            type="button"
            aria-label={isLeftOpen ? "Hide left menu" : "Show left menu"}
            onClick={toggleLeft}
          >
            <MenuIcon className="w-6 h-6 text-gray-800" />
          </ButtonIcon>
          <h1>LOGO</h1>
        </div>
        <div className="flex items-center space-x-3 flex-1 justify-center">
          <div className="relative w-72 max-w-full">
            <label htmlFor="topbar-search" className="sr-only">Search</label>
            <input
              id="topbar-search"
              type="text"
              placeholder="Search..."
              className="w-full py-2 pl-3 pr-10 rounded-full bg-white text-gray-800 placeholder:text-gray-400 border border-gray-300 hover:border-gray-400 focus:outline-none"
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
            aria-label="Notifications"
          >
            <BellIcon className="w-6 h-6 text-gray-800" />
          </ButtonIcon>
          <button
            type="button"
            aria-label="Profile"
            className="rounded-full focus:outline-none"
          >
            <span className="inline-flex items-center justify-center w-8 h-8 rounded-full bg-gray-300 text-gray-800 text-sm font-semibold">U</span>
          </button>
        </div>
      </div>
    </div>
  );
}
