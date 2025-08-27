"use client";

import React, { createContext, useCallback, useContext, useEffect, useMemo, useState } from "react";
import { profile as fetchProfile } from "@/services/auth";

export type ProfileData = { name: string; email: string };

type ProfileContextType = {
  profile: ProfileData | null;
  loading: boolean;
  error: string | null;
  refresh: () => Promise<void>;
  setProfile: (p: ProfileData | null) => void;
};

const UseProfile = createContext<ProfileContextType | undefined>(undefined);

export function ProfileProvider({ children }: { children: React.ReactNode }) {
  const [profile, setProfile] = useState<ProfileData | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  const load = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const res = await fetchProfile();
      if (res && !res.error && res.data) {
        // services/auth.profile() returns WrapperResponse<ProfileResponse>
        // where ProfileResponse = { message: string; data: { name; email } }
        setProfile(res.data.data.user);
      } else {
        setProfile(null);
        setError(res?.error ?? "Failed to load profile");
      }
    } catch {
      setProfile(null);
      setError("Failed to load profile");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    let active = true;
    (async () => {
      await load();
    })();
    return () => {
      active = false;
    };
  }, [load]);

  const value = useMemo(
    () => ({ profile, loading, error, refresh: load, setProfile }),
    [profile, loading, error, load]
  );

  return <UseProfile.Provider value={value}>{children}</UseProfile.Provider>;
}

export function useProfile(): ProfileContextType {
  const ctx = useContext(UseProfile);
  if (!ctx) throw new Error("useProfile must be used within a ProfileProvider");
  return ctx;
}
