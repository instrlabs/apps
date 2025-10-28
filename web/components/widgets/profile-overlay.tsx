"use client";

import React from "react";

import { useProfile } from "@/hooks/useProfile";
import { logout } from "@/services/auth";
import Avatar from "@/components/avatar";
import Button from "@/components/button";

export default function ProfileOverlay() {
  const { profile } = useProfile();

  const username = profile?.username || "Guest";
  const email = profile?.email || "";

  return (
    <div className="flex w-full flex-col items-center gap-3 rounded-lg border border-white/10 bg-white/8 p-4 md:w-[300px]">
      <Avatar size="xl" name={username} />
      <div className="flex w-full flex-col items-center gap-0">
        <h3 className="text-xl font-semibold text-white">{username}</h3>
        <p className="text-base font-normal text-white/60">{email}</p>
      </div>
      <Button
        variant="secondary"
        size="sm"
        className="w-full"
        onClick={async () => {
          await logout();
          document.location.href = "/";
        }}
      >
        Logout
      </Button>
    </div>
  );
}
