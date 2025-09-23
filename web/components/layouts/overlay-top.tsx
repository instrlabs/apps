"use client";

import Image from "next/image";

import { useOverlay } from "@/hooks/useOverlay";
import ButtonIcon from "@/components/actions/button-icon";
import MenuIcon from "@/components/icons/menu";
import SearchIcon from "@/components/icons/search";
import BellIcon from "@/components/icons/bell";
import Avatar from "@/components/avatar";
import TextField from "@/components/inputs/text-field";
import {useProfile} from "@/hooks/useProfile";

export default function OverlayTop() {
  const { profile } = useProfile();
  const { openRight } = useOverlay();

  return (
    <div className="absolute top-0 left-0 right-0 ">
      <div className="h-[80px] w-full p-3 flex flex-row justify-between items-center">
        <div className="flex items-center space-x-5">
          <ButtonIcon
            xSize="lg"
            xColor="secondary"
            // onClick={() => toggleByKey("left:navigation")}
          >
            <MenuIcon className="w-6 h-6" />
          </ButtonIcon>
          <Image src="/logo.svg" alt="logo" width={40} height={40} />
        </div>
        <div className="flex items-center space-x-5">
          <div className="relative w-96 max-w-full">
            <label htmlFor="topbar-search" className="sr-only">Search</label>
            <TextField
              type="text"
              placeholder="Search..."
              className="pr-10 cursor-pointer"
              xSize="sm"
              // onFocus={() => toggleByKey("modal:search")}
              // onClick={() => toggleByKey("modal:search")}
              readOnly
            />
            <SearchIcon
              className="pointer-events-none absolute right-3 top-1/2 -translate-y-1/2 w-5 h-5"
              aria-hidden="true"
            />
          </div>
        </div>
        <div className="flex items-center space-x-5">
          <ButtonIcon
             xSize="sm"
             xColor="secondary"
             // onClick={() => toggleByKey("right:notifications")}
           >
            <BellIcon className="w-6 h-6 text-gray-200" />
          </ButtonIcon>
          <Avatar
            name={profile?.name}
            size="md"
            onClick={() => openRight("profile")}
          />
        </div>
      </div>
    </div>
  );
}
