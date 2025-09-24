"use client";

import Image from "next/image";

import { useOverlay } from "@/hooks/useOverlay";
import Avatar from "@/components/avatar";
import {useProfile} from "@/hooks/useProfile";

export default function OverlayTop() {
  const { profile } = useProfile();
  const { openRight } = useOverlay();

  return (
    <div className="absolute top-0 left-0 right-0 ">
      <div className="h-[80px] w-full p-3 flex flex-row justify-between items-center">
        <div className="flex items-center space-x-5">
          <Image src="/logo.svg" alt="logo" width={40} height={40} />
        </div>
        <div className="flex items-center space-x-5">
          <Avatar
            xsize="sm"
            name={profile?.name}
            onClick={() => openRight("profile")}
          />
        </div>
      </div>
    </div>
  );
}
