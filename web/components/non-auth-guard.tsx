"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import ROUTES from "@/constants/routes";

/**
 * NonAuthGuard
 *
 * Client-side guard to prevent authenticated users from accessing non-auth pages
 * (e.g., login, register, forgot/reset password). If a local auth token exists,
 * we redirect the user to the home page.
 *
 * Minimal implementation to avoid server-side changes. This will apply to all
 * pages under the (non-auth) route group when included in its layout.
 */
export default function NonAuthGuard() {
  const router = useRouter();

  useEffect(() => {
    try {
      const token = typeof window !== "undefined" ? localStorage.getItem("auth_token") : null;
      if (token) {
        router.replace(ROUTES.HOME);
      }
    } catch {
      // Ignore storage access errors
    }
  }, [router]);

  return null;
}
