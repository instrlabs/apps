"use client";

import React from "react";
import Avatar from "@/components/avatar";
import EditIcon from "@/components/icons/edit";
import LockIcon from "@/components/icons/lock";
import LogoutIcon from "@/components/icons/logout";
import MenuButton from "@/components/menu-button";
import { useProfile } from "@/hooks/useProfile";
import { logoutUser } from "@/services/auth";
import { useRouter } from "next/navigation";

export default function ProfileOverlay() {
  const { profile, setProfile } = useProfile();
  const router = useRouter();

  const name = profile?.name || "Guest";
  const email = profile?.email || "";

  const handleLogout = async () => {
    await logoutUser();
    setProfile(null);
    router.replace("/login");
  };

  return (
    <div className="w-full h-full bg-card shadow-primary rounded-xl">
      <div className="flex flex-col gap-4 py-10">
        <div className="mx-auto">
          <Avatar name={name} size={60} />
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
        <MenuButton onClick={handleLogout} icon={<LogoutIcon className="h-5 w-5" aria-hidden="true" />}>
          Logout
        </MenuButton>
      </div>
    </div>
  );
}
