"use server";

import { Suspense } from "react";
import CookieActions from "./CookieActions";

export default async function CookiesPage() {
  return (
    <div className="container mx-auto max-w-6xl p-6">
      <div className="mb-8">
        <h1 className="mb-2 text-3xl font-bold">Cookie Simulation</h1>
        <p className="mb-6 text-gray-400">Demonstrate server-side and client-side cookie operations in Next.js</p>

        <div className="mb-6 rounded-lg border border-blue-500/30 bg-blue-900/20 p-4">
          <h2 className="mb-2 text-lg font-semibold text-blue-300">Understanding Cookies</h2>
          <div className="grid gap-4 text-sm md:grid-cols-2">
            <div>
              <h3 className="mb-2 font-medium">Server-side Operations</h3>
              <ul className="space-y-1 text-gray-300">
                <li>
                  • Using <code className="rounded bg-gray-800 px-1">cookies()</code> from next/headers
                </li>
                <li>• HttpOnly, Secure, SameSite options</li>
                <li>• Automatic encryption/signed cookies</li>
                <li>• Path and domain scoping</li>
              </ul>
            </div>
            <div>
              <h3 className="mb-2 font-medium">Client-side Operations</h3>
              <ul className="space-y-1 text-gray-300">
                <li>
                  • Using <code className="rounded bg-gray-800 px-1">document.cookie</code>
                </li>
                <li>• Limited access to HttpOnly cookies</li>
                <li>• Browser security restrictions</li>
                <li>• Size limitations (4KB per cookie)</li>
              </ul>
            </div>
          </div>
        </div>
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
