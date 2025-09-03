"use client";

import React, { createContext, useContext, useMemo, useState } from "react";

export type ProfileData = { name: string; email: string };

type ProfileContextType = {
  profile: ProfileData;
  loading: boolean;
  setProfile: (p: ProfileData) => void;
};

const Profile = createContext<ProfileContextType | undefined>(undefined);

export function ProfileProvider({ children, data }: {
  children: React.ReactNode,
  data: ProfileData
}) {
  const [profile, setProfile] = useState<ProfileData>(data);
  const [loading, setLoading] = useState<boolean>(false);

  const value = useMemo(
    () => ({ profile, loading, setProfile }),
    [profile, loading]
  );

  return <Profile.Provider value={value}>{children}</Profile.Provider>;
}

export function useProfile(): ProfileContextType {
  const ctx = useContext(Profile);
  if (!ctx) throw new Error("useProfile must be used within a ProfileProvider");
  return ctx;
}
