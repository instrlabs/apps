"use client";

import { useOverlay } from "@/hooks/useOverlay";
import ButtonIcon from "@/components/button-icon";
import MenuIcon from "@/components/icons/menu";
import SearchIcon from "@/components/icons/search";
import BellIcon from "@/components/icons/bell";
import Avatar from "@/components/avatar";

export default function OverlayTop() {
  const { isLeftOpen, toggleLeft, isRightOpen, toggleRight } = useOverlay();
  return (
    <div className="absolute top-0 left-0 right-0 w-full p-3">
      <div className="h-[60px] rounded-full bg-neutral-50 flex flex-row justify-between items-center px-3">
        <div className="flex items-center space-x-3">
          <ButtonIcon
            type="button"
            aria-label={isLeftOpen ? "Hide left menu" : "Show left menu"}
            onClick={toggleLeft}
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
            aria-label={isRightOpen ? "Hide right overlay" : "Show right overlay"}
            onClick={toggleRight}
          >
            <BellIcon className="w-6 h-6 text-gray-800" />
          </ButtonIcon>
          <button
            type="button"
            aria-label="Profile"
            className="rounded-full focus:outline-none"
            onClick={() => {}}
          >
            <Avatar name="Artha Dede" size={40} />
          </button>
        </div>
      </div>
    </div>
  );
}
