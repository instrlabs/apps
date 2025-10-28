"use client";

import { useEffect, useMemo, useState } from "react";

export const tailwindScreens = {
  sm: 640,
  md: 768,
  lg: 1024,
  xl: 1280,
  "2xl": 1536,
} as const;

export type TailwindBreakpoint = keyof typeof tailwindScreens; // "sm" | "md" | "lg" | "xl" | "2xl"

function toMediaQuery(input: string): string {
  if (/[()]/.test(input)) return input;

  const isMax = input.startsWith("max-");
  const token = (isMax ? input.slice(4) : input) as TailwindBreakpoint;

  const px = tailwindScreens[token];
  if (!px) return input;

  if (isMax) {
    const maxPx = px - 1;
    return `(max-width: ${maxPx}px)`;
  }

  return `(min-width: ${px}px)`;
}

export default function useMediaQuery(queryOrToken: string): boolean {
  const query = useMemo(() => toMediaQuery(queryOrToken), [queryOrToken]);

  const getMatch = () => {
    if (typeof window === "undefined" || typeof window.matchMedia !== "function") {
      return false;
    }
    return window.matchMedia(query).matches;
  };

  const [matches, setMatches] = useState<boolean>(getMatch);

  useEffect(() => {
    if (typeof window === "undefined" || typeof window.matchMedia !== "function") {
      return;
    }

    const mql = window.matchMedia(query);
    const onChange = (e: MediaQueryListEvent) => setMatches(e.matches);

    setMatches(mql.matches);

    mql.addEventListener?.("change", onChange);
    mql.addListener?.(onChange);

    return () => {
      mql.removeEventListener?.("change", onChange);
      mql.removeListener?.(onChange);
    };
  }, [query]);

  return matches;
}

export function useMobile(): boolean {
  return useMediaQuery("max-md");
}
