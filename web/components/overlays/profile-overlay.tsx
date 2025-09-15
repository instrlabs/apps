"use client";

import React from "react";
import { useRouter } from "next/navigation";

import Avatar from "@/components/avatar";
import EditIcon from "@/components/icons/edit";
import LockIcon from "@/components/icons/lock";
import LogoutIcon from "@/components/icons/logout";
import MenuButton from "@/components/actions/menu-button";
import { useProfile } from "@/hooks/useProfile";
import { logoutUser } from "@/services/authentications";

export default function ProfileOverlay() {
  const router = useRouter();
  const { profile } = useProfile();

  const name = profile?.name || "Guest";
  const email = profile?.email || "";

  return (
    <div className="w-[400px] h-full bg-card shadow-primary rounded-xl">
      <div className="flex flex-col gap-4 py-10">
        <div className="mx-auto">
          <Avatar name={name} size="lg" />
        </div>
        <div className="flex flex-col">
          <h3 className="text-2xl font-bold text-foreground text-center">{name}</h3>
          <p className="text-base text-muted text-center">{email}</p>
        </div>
      </div>
      <div className="flex flex-col gap-2 p-3">
        <MenuButton onClick={() => router.push("/edit-profile")} icon={<EditIcon className="h-5 w-5" aria-hidden="true" />}>
          Edit Profile
        </MenuButton>
        <MenuButton onClick={() => router.push("/change-password")} icon={<LockIcon className="h-5 w-5" aria-hidden="true" />}>
          Change Password
        </MenuButton>
        <MenuButton onClick={() => logoutUser()} icon={<LogoutIcon className="h-5 w-5" aria-hidden="true" />}>
          Logout
        </MenuButton>
      </div>
    </div>
  );
}
