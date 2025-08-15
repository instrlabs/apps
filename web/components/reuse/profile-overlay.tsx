import React from "react";
import Avatar from "@/components/avatar";
import EditIcon from "@/components/icons/edit";
import LockIcon from "@/components/icons/lock";
import LogoutIcon from "@/components/icons/logout";
import MenuButton from "@/components/menu-button";

export default function ProfileOverlay() {
  return (
    <div className="w-full h-full bg-white shadow-primary rounded-xl">
      <div className="flex flex-col gap-4 py-10">
        <div className="mx-auto">
          <Avatar name="Artha Dede" size={60} />
        </div>
        <div className="flex flex-col">
          <h3 className="text-2xl font-bold text-gray-800 text-center">Artha Suryawan</h3>
          <p className="text-base text-gray-500 text-center">arthadede@gmail.com</p>
        </div>
      </div>
      <div className="flex flex-col gap-2 p-3">
        <MenuButton icon={<EditIcon className="h-5 w-5" aria-hidden="true" />}>
          Edit Profile
        </MenuButton>
        <MenuButton icon={<LockIcon className="h-5 w-5" aria-hidden="true" />}>
          Change Password
        </MenuButton>
        <MenuButton icon={<LogoutIcon className="h-5 w-5" aria-hidden="true" />}>
          Logout
        </MenuButton>
      </div>
    </div>
  );
}
