"use client";

import { useEffect, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { handleGoogleCallback } from "@/services/api/auth";

export default function GoogleCallbackPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [error, setError] = useState<string | null>(null);

  // useEffect(() => {
  //   async function processCallback() {
  //     try {
  //       // Get the authorization code from the URL query parameters
  //       const code = searchParams.get("code");
  //
  //       if (!code) {
  //         throw new Error("No authorization code received from Google");
  //       }
  //
  //       // Process the callback with the authorization code
  //       const data = await handleGoogleCallback(code);
  //
  //       // Store the authentication token
  //       if (data.access_token) {
  //         localStorage.setItem("authToken", data.access_token);
  //
  //         router.push("/apps");
  //       } else {
  //         throw new Error("No authentication token received");
  //       }
  //     } catch (err) {
  //       console.error("Google authentication error:", err);
  //       setError(err instanceof Error ? err.message : "An error occurred during Google authentication");
  //     }
  //   }
  //
  //   processCallback();
  // }, [router, searchParams]);

  // Show a loading state or error message
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
