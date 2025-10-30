"use server";

import { Suspense } from "react";
import CookieActions from "./CookieActions";

export default async function CookiesPage() {
  return (
    <div className="container mx-auto max-w-6xl p-6">
      <div className="mb-8">
        <h1 className="mb-2 text-3xl font-bold">Test Cookies</h1>
        <p className="mb-6 text-gray-400">Simulate create/delete cookies for testing</p>
      </div>

      <Suspense
        fallback={
          <div className="rounded-lg border border-white/20 bg-white/5 p-4">
            <div className="animate-pulse">
              <div className="mb-2 h-4 rounded bg-gray-700"></div>
              <div className="mb-2 h-4 rounded bg-gray-700"></div>
              <div className="h-4 rounded bg-gray-700"></div>
            </div>
          </div>
        }
      >
        <CookieActions />
      </Suspense>
    </div>
  );
}
