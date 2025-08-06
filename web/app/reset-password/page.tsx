"use client";

import React, { useState, useEffect } from "react";
import { useRouter, useSearchParams } from "next/navigation";

import { resetPassword } from "@/services/auth";
import Button from "@/components/button";
import { useNotification } from "@/components/notification";

export default function ResetPasswordPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { showNotification } = useNotification();
  
  const [token, setToken] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [isSubmitted, setIsSubmitted] = useState(false);
  const [tokenError, setTokenError] = useState(false);

  useEffect(() => {
    const tokenParam = searchParams?.get("token");
    if (!tokenParam) {
      setTokenError(true);
      showNotification(
        "Invalid or missing reset token. Please request a new password reset.",
        "error",
        5000
      );
    } else {
      setToken(tokenParam);
    }
  }, [searchParams, showNotification]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (password !== confirmPassword) {
      showNotification("Passwords do not match", "error", 5000);
      return;
    }

    if (password.length < 8) {
      showNotification("Password must be at least 8 characters long", "error", 5000);
      return;
    }

    setIsLoading(true);

    const { data, error } = await resetPassword(token, password);
    
    if (error) {
      showNotification(error, "error", 5000);
    } else {
      setIsSubmitted(true);
      showNotification(
        "Your password has been reset successfully",
        "success",
        5000
      );
    }
    
    setIsLoading(false);
  };

  if (tokenError) {
    return (
      <div className="h-screen w-full flex flex-col justify-center items-center p-10">
        <h2 className="text-2xl font-semibold mb-6">Invalid Reset Link</h2>
        <p className="text-center mb-6 max-w-sm">
          The password reset link is invalid or has expired. Please request a new password reset.
        </p>
        <div className="flex flex-col gap-5 w-full max-w-sm">
          <Button onClick={() => router.push("/forgot-password")}>
            Request New Reset Link
          </Button>
        </div>
      </div>
    );
  }

  if (isSubmitted) {
    return (
      <div className="h-screen w-full flex flex-col justify-center items-center p-10">
        <h2 className="text-2xl font-semibold mb-6">Password Reset Complete</h2>
        <p className="text-center mb-6 max-w-sm">
          Your password has been reset successfully. You can now log in with your new password.
        </p>
        <div className="flex flex-col gap-5 w-full max-w-sm">
          <Button onClick={() => router.push("/login")}>
            Go to Login
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="h-screen w-full flex flex-col justify-center items-center p-10">
      <h2 className="text-2xl font-semibold mb-6">Create New Password</h2>
      <p className="text-center mb-6 max-w-sm">
        Enter your new password below. Password must be at least 8 characters long.
      </p>
      <form onSubmit={handleSubmit} className="flex flex-col gap-5 w-full max-w-sm">
        <div className="space-y-1">
          <label htmlFor="password" className="text-sm font-medium">
            New Password
          </label>
          <input
            id="password"
            type="password"
            placeholder="Enter your new password"
            className="px-2 py-2.5 rounded w-full outline-none text-sm border border-gray-300"
            value={password}
            onChange={e => setPassword(e.target.value)}
            required
            minLength={8}
          />
        </div>
        <div className="space-y-1">
          <label htmlFor="confirmPassword" className="text-sm font-medium">
            Confirm Password
          </label>
          <input
            id="confirmPassword"
            type="password"
            placeholder="Confirm your new password"
            className="px-2 py-2.5 rounded w-full outline-none text-sm border border-gray-300"
            value={confirmPassword}
            onChange={e => setConfirmPassword(e.target.value)}
            required
            minLength={8}
          />
        </div>
        <Button type="submit" isLoading={isLoading} loadingText="Resetting...">
          Reset Password
        </Button>
      </form>
    </div>
  );
}