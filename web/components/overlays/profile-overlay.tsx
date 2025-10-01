"use client";

import React from "react";
import { redirect, RedirectType } from "next/navigation";

import { useProfile } from "@/hooks/useProfile";
import MenuButton from "@/components/actions/menu-button";
import { logout } from "@/services/auth";

export default function ProfileOverlay() {
  const { profile } = useProfile();

  return (
    <div className="h-full w-[300px] p-4 pl-0 pt-0">
      <div className="bg-primary-black shadow-primary h-full w-full rounded-lg">
        <div className="flex flex-col p-4">
          <h3 className="text-sm font-light text-white">{profile?.username}</h3>
          <p className="text-sm font-light text-white/60">{profile?.email}</p>
        </div>
        <div className="flex flex-col">
          <MenuButton xSize="sm" onClick={() => redirect("/", RedirectType.push)}>
            Dashboard
          </MenuButton>
          <hr className="my-2" />
          <MenuButton xSize="sm" onClick={logout}>
            Logout
          </MenuButton>
          <hr className="my-2" />
        </div>
      </div>
    </div>
  );
}
