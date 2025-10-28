"use client";

import React, {createContext, useContext, useEffect, useMemo, useState} from "react";

import { User, logout } from "@/services/auth";


type ProfileContextType = {
  isLoggedIn: boolean;
  profile: User | null;
  setProfile: (p: User) => void;
};

const Profile = createContext<ProfileContextType | undefined>(undefined);

export function ProfileProvider({ children, data }: {
  children: React.ReactNode,
  data: User | null
}) {
  const [profileData, setProfileData] = useState<User | null>(data);

  const value = useMemo(
    () => ({
      isLoggedIn: !!profileData,
      profile: profileData,
      setProfile: setProfileData
    }),
    [profileData]
  );

  return <Profile.Provider value={value}>{children}</Profile.Provider>;
}

export function useProfile(): ProfileContextType {
  const ctx = useContext(Profile);
  if (!ctx) throw new Error("useProfile must be used within a ProfileProvider");
  return ctx;
}
