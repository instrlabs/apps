"use client";

import React, { useState } from "react";
import { useRouter } from "next/navigation";

import { registerUser } from "@/services/api/auth";
import GoogleSignInButton from "@/components/google-signin";
import Button from "@/components/button";
import { useNotification } from "@/components/notification";

export default function SignupPage() {
  const router = useRouter();
  const { showNotification } = useNotification();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);

    try {
      const data = await registerUser(email, password);
      showNotification(data.message, "success", 2000);
      setTimeout(() => router.replace("/login"), 2500);
    } catch (err) {
      showNotification(
        err instanceof Error ? err.message : "An error occurred during registration",
        "error",
        5000
      );
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="h-screen w-full flex flex-col justify-center items-center p-10">
      <h2 className="text-2xl font-semibold mb-6">Create your account</h2>
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
            placeholder="Create a password"
            className="px-2 py-2.5 rounded w-full outline-none text-sm border border-gray-300"
            value={password}
            onChange={e => setPassword(e.target.value)}
            required
          />
        </div>
        <Button type="submit" isLoading={isLoading} loadingText="Signing up...">
          Sign up
        </Button>
      </form>
      <div className="flex flex-col gap-5 w-sm mt-5">
        <div className="text-sm text-center">
          Already have an account?{" "}
          <a className="underline" href="/login">
            Sign in
          </a>
        </div>
      </div>
    </div>
  );
}
