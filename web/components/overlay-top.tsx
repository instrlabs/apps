"use client";

import { useOverlay } from "@/hooks/useOverlay";
import ButtonIcon from "@/components/button-icon";
import MenuIcon from "@/components/icons/menu";
import SearchIcon from "@/components/icons/search";
import BellIcon from "@/components/icons/bell";
import Avatar from "@/components/avatar";
import ProfileOverlay from "@/components/reuse/profile-overlay";
import NotificationOverlay from "@/components/reuse/notification-overlay";

export default function OverlayTop() {
  const {
    isLeftOpen,
    isRightOpen,
    setLeftNode,
    toggleLeftByKey,
    leftActiveKey,
    setRightNode,
    setRightWidth,
    toggleRightByKey,
    rightActiveKey,
    setModalNode,
    toggleModalByKey,
    isModalOpen,
    modalActiveKey,
  } = useOverlay();

  function openSearchModal() {
    const key = 'modal:search';

    if (isModalOpen && modalActiveKey === key) {
      toggleModalByKey(key);
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
    toggleModalByKey(key);
  }

  return (
    <div className="absolute top-0 left-0 right-0 w-full p-2">
      <div className="h-[60px] w-full flex flex-row justify-between items-center px-3">
        <div className="flex items-center space-x-3">
          <ButtonIcon
            type="button"
            aria-label={isLeftOpen ? "Hide left" : "Show left (menu)"}
            onClick={() => {
              const key = 'left:menu';

              if (isLeftOpen && leftActiveKey === key) {
                toggleLeftByKey(key);
                return;
              }

              setLeftNode(
                <ul className="space-y-2 text-sm text-gray-700">
                  <li><button className="w-full text-left px-3 py-2 rounded hover:bg-gray-100">Dashboard</button></li>
                  <li><button className="w-full text-left px-3 py-2 rounded hover:bg-gray-100">Projects</button></li>
                  <li><button className="w-full text-left px-3 py-2 rounded hover:bg-gray-100">Teams</button></li>
                  <li><button className="w-full text-left px-3 py-2 rounded hover:bg-gray-100">Settings</button></li>
                </ul>
              );
              toggleLeftByKey(key);
            }}
          >
            <MenuIcon className="w-6 h-6 text-gray-800" />
          </ButtonIcon>
          <h1 className="text-xl">LOGO</h1>
        </div>
        <div className="flex items-center space-x-3 flex-1 justify-center">
          <div className="relative w-96 max-w-full">
            <label htmlFor="topbar-search" className="sr-only">Search</label>
            <input
              id="topbar-search"
              type="text"
              placeholder="Search..."
              className="w-full py-2 pl-5 pr-10 rounded-full bg-blue-50 hover:bg-blue-100 placeholder:text-gray-600 focus:outline-none"
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
          <ButtonIcon
            type="button"
            aria-label="Show notifications in right overlay"
            onClick={() => {
              const key = 'right:notifications';
              if (isRightOpen && rightActiveKey === key) {
                toggleRightByKey(key);
                return;
              }

              setRightNode(
                <NotificationOverlay />
              );
              toggleRightByKey(key);
            }}
          >
            <BellIcon className="w-6 h-6 text-gray-800" />
          </ButtonIcon>
          <button
            type="button"
            aria-label="Profile"
            className="rounded-full focus:outline-none"
            onClick={() => {
              const key = 'right:profile';

              if (isRightOpen && rightActiveKey === key) {
                toggleRightByKey(key);
                return;
              }

              setRightWidth(400);
              setRightNode(<ProfileOverlay />);
              toggleRightByKey(key);
            }}
          >
            <Avatar name="Artha Dede" size={40} />
          </button>
        </div>
      </div>
    </div>
  );
}
