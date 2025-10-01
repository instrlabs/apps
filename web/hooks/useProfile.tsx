"use client";

import React, {createContext, useContext, useEffect, useMemo, useState} from "react";
import { redirect } from "next/navigation";

import { User, getProfile, refresh as refreshToken } from "@/services/auth";


type ProfileContextType = {
  profile: User | null;
  loading: boolean;
  setProfile: (p: User) => void;
};

const Profile = createContext<ProfileContextType | undefined>(undefined);

export function ProfileProvider({ children, data }: {
  children: React.ReactNode,
  data: User | null
}) {
  const [profileData, setProfileData] = useState<User | null>(data);
  const [loading, setLoading] = useState<boolean>(false);

  useEffect(() => {
    async function refreshProfile() {
      await refreshToken();
      const { success, data } = await getProfile();

      if (!success) redirect("/login");
      setProfileData((data as { user: User } | null)?.user ?? null);
    }

    if (!profileData) {
      refreshProfile().then();
    }
  }, [profileData])

  const value = useMemo(
    () => ({ profile: profileData, loading, setProfile: setProfileData }),
    [profileData, loading]
  );

  if (!profileData) {
    return <div>Loading...</div>
  }

  return <Profile.Provider value={value}>{children}</Profile.Provider>;
}

export function useProfile(): ProfileContextType {
  const ctx = useContext(Profile);
  if (!ctx) throw new Error("useProfile must be used within a ProfileProvider");
  return ctx;
}
