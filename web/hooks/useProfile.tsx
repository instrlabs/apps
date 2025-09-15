"use client";

import React, {createContext, useContext, useEffect, useMemo, useState} from "react";
import { redirect } from "next/navigation";

import { ProfileResponse, getProfile, refreshToken } from "@/services/authentications";


type ProfileContextType = {
  profile: ProfileResponse | null;
  loading: boolean;
  setProfile: (p: ProfileResponse) => void;
};

const Profile = createContext<ProfileContextType | undefined>(undefined);

export function ProfileProvider({ children, data }: {
  children: React.ReactNode,
  data: ProfileResponse | null
}) {
  const [profileData, setProfileData] = useState<ProfileResponse | null>(data);
  const [loading, setLoading] = useState<boolean>(false);

  useEffect(() => {
    async function refresh() {
      await refreshToken();
      const { success, data } = await getProfile();

      if (!success) redirect("/login");
      setProfileData(data as ProfileResponse);
    }

    if (!profileData) {
      refresh().then()
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
