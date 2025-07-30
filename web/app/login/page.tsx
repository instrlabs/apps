"use client";

import React, { useState } from "react";
import { useRouter } from "next/navigation";

import { loginUser } from "@/services/api/auth";
import GoogleSignInButton from "@/components/google-signin";
import Button from "@/components/button";
import { useNotification } from "@/components/notification";

export default function LoginPage() {
  const router = useRouter();
  const { showNotification } = useNotification();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);

    try {
      const data = await loginUser(email, password);

      if (data.token) {
        localStorage.setItem("authToken", data.token);
      }

      router.push("/apps");
    } catch (err) {
      showNotification(
        err instanceof Error ? err.message : "An error occurred during login",
        "error",
        5000
      );
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="h-screen w-full flex flex-col justify-center items-center p-10">
      <h2 className="text-2xl font-semibold mb-6">Log in to your account</h2>
      <form onSubmit={handleSubmit} className="flex flex-col gap-5 w-full max-w-sm">
        <div className="space-y-1">
          <label htmlFor="email" className="text-sm font-medium">
            Email
          </label>
          <input
            id="email"
            type="email"
            placeholder="Enter your email address"
            className="px-2 py-2.5 rounded w-full outline-none text-sm border border-gray-300"
            value={email}
            onChange={e => setEmail(e.target.value)}
            required
          />
        </div>
        <div className="space-y-1">
          <label htmlFor="password" className="text-sm font-medium">
            Password
          </label>
          <input
            id="password"
            type="password"
            placeholder="Enter your password"
            className="px-2 py-2.5 rounded w-full outline-none text-sm border border-gray-300"
            value={password}
            onChange={e => setPassword(e.target.value)}
            required
          />
          <div className="text-right">
            <a className="text-sm text-blue-600 hover:underline" href="/forgot-password">
              Forgot password?
            </a>
          </div>
        </div>
        <Button type="submit" isLoading={isLoading} loadingText="Signing in...">
          Sign in
        </Button>
      </form>
      <div className="flex flex-col gap-5 w-sm">
        <GoogleSignInButton />
        <div className="text-sm text-center">
          Don&#39;t have an account?{" "}
          <a className="underline" href="/signup">
            Sign up
          </a>
        </div>
      </div>
    </div>
  );
}
