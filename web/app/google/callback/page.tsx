"use client";

import { useEffect, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { handleGoogleCallback } from "@/services/auth";
import ROUTES from "@/constants/routes";

export default function GoogleCallbackPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function processCallback() {
      const code = searchParams.get("code");

      if (!code) {
        setError("No authorization code received from Google");
        return;
      }

      const { data, error } = await handleGoogleCallback(code);

      if (error) {
        setError(error);
        return;
      }

      if (data?.data.access_token) {
        localStorage.setItem("authToken", data?.data.access_token);
        document.cookie = `authToken=${data.data.access_token}; path=/; max-age=86400; samesite=lax`;
        router.push(ROUTES.HOME);
      }
    }

    processCallback().then();
  }, [router, searchParams]);

  return (
    <div className="h-screen w-full flex flex-col justify-center items-center p-10">
      {error ? (
        <div className="text-center">
          <h2 className="text-2xl font-semibold mb-4 text-red-600">Authentication Failed</h2>
          <p className="mb-6">{error}</p>
          <button
            onClick={() => router.push("/login")}
            className="px-4 py-2 bg-black text-white rounded-full hover:bg-gray-800"
          >
            Return to Login
          </button>
        </div>
      ) : (
        <div className="text-center">
          <h2 className="text-2xl font-semibold mb-4">Completing Authentication</h2>
          <p className="mb-6">Please wait while we complete your Google sign-in...</p>
          <div className="w-8 h-8 border-4 border-t-black border-r-gray-200 border-b-gray-200 border-l-gray-200 rounded-full animate-spin mx-auto"></div>
        </div>
      )}
    </div>
  );
}