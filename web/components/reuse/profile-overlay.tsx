import React from "react";
import Avatar from "@/components/avatar";
import EditIcon from "@/components/icons/edit";
import LockIcon from "@/components/icons/lock";
import LogoutIcon from "@/components/icons/logout";

export default function ProfileOverlay() {
  return (
    <div className="w-full h-full bg-blue-50 p-2">
      <div className="flex flex-col gap-4 my-10">
        <div className="mx-auto">
          <Avatar name="Artha Dede" size={60} />
        </div>
        <div className="flex flex-col">
          <h3 className="text-2xl font-bold text-gray-800 text-center">Artha Suryawan</h3>
          <p className="text-base text-gray-500 text-center">arthadede@gmail.com</p>
        </div>
      </div>
      <div className="flex flex-col gap-1 my-10">
        <button className="w-full flex items-center text-sm gap-3 px-3 py-2 text-gray-800 hover:bg-blue-200 rounded-xl">
          <EditIcon className="h-5 w-5" aria-hidden="true" />
          Edit Profile
        </button>
        <button className="w-full flex items-center text-sm gap-3 px-3 py-2 text-gray-800 hover:bg-blue-200 rounded-xl">
          <LockIcon className="h-5 w-5" aria-hidden="true" />
          Change Password
        </button>
        <button className="w-full flex items-center text-sm gap-3 px-3 py-2 text-gray-800 hover:bg-blue-200 rounded-xl">
          <LogoutIcon className="h-5 w-5" aria-hidden="true" />
          Logout
        </button>
      </div>
    </div>
  );
}
