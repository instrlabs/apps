"use client";

import Image from "next/image";

import { useOverlay } from "@/hooks/useOverlay";
import ButtonIcon from "@/components/button-icon";
import MenuIcon from "@/components/icons/menu";
import SearchIcon from "@/components/icons/search";
import BellIcon from "@/components/icons/bell";
import Avatar from "@/components/avatar";
import TextField from "@/components/text-field";
import {useProfile} from "@/hooks/useProfile";

export default function OverlayTop() {
  const { toggleByKey } = useOverlay();
  const { profile } = useProfile();

  return (
    <div className="absolute top-0 left-0 right-0 w-full p-2">
      <div className="h-[60px] w-full flex flex-row justify-between items-center px-3">
        <div className="flex items-center space-x-5">
          <ButtonIcon type="button" xSize="lg" xColor="secondary" onClick={() => toggleByKey("left:navigation")}>
            <MenuIcon className="w-6 h-6" />
          </ButtonIcon>
          <Image src="/logo.svg" alt="logo" width={40} height={40} />
        </div>
        <div className="flex items-center space-x-5 flex-1 justify-center">
          <div className="relative w-96 max-w-full">
            <label htmlFor="topbar-search" className="sr-only">Search</label>
            <TextField
              id="topbar-search"
              type="text"
              placeholder="Search..."
              className="pr-10 cursor-pointer"
              xIsRounded
              xSize="md"
              onFocus={() => toggleByKey("modal:search")}
              onClick={() => toggleByKey("modal:search")}
              readOnly
            />
            <SearchIcon
              className="pointer-events-none absolute right-3 top-1/2 -translate-y-1/2 w-5 h-5"
              aria-hidden="true"
            />
          </div>
        </div>
        <div className="flex items-center space-x-3">
          <ButtonIcon
             type="button"
             xSize="lg"
             xColor="secondary"
             onClick={() => toggleByKey("right:notifications")}
           >
            <BellIcon className="w-6 h-6 text-gray-800" />
          </ButtonIcon>
          <ButtonIcon
            type="button"
            xSize="lg"
            xColor="secondary"
            aria-label="Open profile overlay"
            className="p-0 rounded-full bg-transparent shadow-none hover:bg-transparent"
            onClick={() => toggleByKey("right:profile")}
          >
            <Avatar name={profile?.name} size={40} />
          </ButtonIcon>
        </div>
      </div>
    </div>
  );
}
