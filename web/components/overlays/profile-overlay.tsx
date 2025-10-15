"use client";

import React from "react";
import { useRouter } from "next/navigation";

import { useProfile } from "@/hooks/useProfile";
import MenuButton from "@/components/actions/menu-button";
import { logout } from "@/services/auth";
import Avatar from "@/components/avatar";
import Button from "@/components/actions/button";

export default function ProfileOverlay() {
  const { profile } = useProfile();
  const router = useRouter();

  const username = profile?.username || "Guest";
  const email = profile?.email || "";

  return (
    <div className="w-[400px] h-full ml-4 bg-white/10 rounded p-4 flex flex-col gap-4">
      <div className="flex justify-center">
        <Avatar xSize="lg" name={username} />
      </div>
      <div className="flex flex-col gap">
        <h3 className="text-base font-semibold text-white text-center">{username}</h3>
        <p className="text-sm font-light text-white/60 text-center">{email}</p>
      </div>
      <div className="flex flex-col gap-2">
        <Button onClick={() => router.push("/")}>
          Dashboard
        </Button>
        <Button xColor="secondary" onClick={logout}>
          Logout
        </Button>
      </div>
    </div>
  );
}
