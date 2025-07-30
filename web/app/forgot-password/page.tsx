"use client";

import React, { useState } from "react";
import { useRouter } from "next/navigation";

import { requestPasswordReset } from "@/services/auth";
import Button from "@/components/button";
import { useNotification } from "@/components/notification";
import { ROUTES } from "@/constants/routes";

export default function ForgotPasswordPage() {
  const router = useRouter();
  const { showNotification } = useNotification();
  const [email, setEmail] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [isSubmitted, setIsSubmitted] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);

    try {
      await requestPasswordReset(email);
      setIsSubmitted(true);
      showNotification(
        "Password reset instructions have been sent to your email",
        "success",
        5000
      );
    } catch (err) {
      showNotification(
        err instanceof Error ? err.message : "An error occurred during password reset request",
        "error",
        5000
      );
    } finally {
      setIsLoading(false);
    }
  };

  if (isSubmitted) {
    return (
      <div className="h-screen w-full flex flex-col justify-center items-center p-10">
        <h2 className="text-2xl font-semibold mb-6">Check your email</h2>
        <p className="text-center mb-6 max-w-sm">
          We&#39;ve sent password reset instructions to <strong>{email}</strong>. 
          Please check your inbox and follow the link in the email.
        </p>
        <div className="flex flex-col gap-5 w-full max-w-sm">
          <Button onClick={() => router.push(ROUTES.LOGIN)}>
            Return to login
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="h-screen w-full flex flex-col justify-center items-center p-10">
      <h2 className="text-2xl font-semibold mb-6">Reset your password</h2>
      <p className="text-center mb-6 max-w-sm">
        Enter your email address and we&#39;ll send you instructions to reset your password.
      </p>
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
        <Button type="submit" isLoading={isLoading} loadingText="Sending...">
          Send reset instructions
        </Button>
        <div className="text-sm text-center">
          Remember your password?{" "}
          <a className="underline" href={ROUTES.LOGIN}>
            Sign in
          </a>
        </div>
      </form>
    </div>
  );
}