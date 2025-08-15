"use client";

import React, { useState, useEffect, Suspense } from "react";
import { useRouter, useSearchParams } from "next/navigation";

import { resetPassword } from "@/services/auth";
import Button from "@/components/button";
import { useNotification } from "@/components/notification";
import FormInput from "@/components/form-input";
import { ROUTES } from "@/constants/routes";

function LoadingFallback() {
  return (
    <div className="h-screen w-full flex flex-col justify-center items-center p-10">
      <div className="text-center">
        <h2 className="text-2xl font-semibold mb-4">Loading...</h2>
        <div className="w-8 h-8 border-4 border-t-black border-r-gray-200 border-b-gray-200 border-l-gray-200 rounded-full animate-spin mx-auto"></div>
      </div>
    </div>
  );
}

// Component that uses useSearchParams
function ResetPasswordContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { showNotification } = useNotification();

  const [token, setToken] = useState("");
  const [formData, setFormData] = useState({ password: "", confirmPassword: "" });
  const [fieldErrors, setFieldErrors] = useState<{ password?: string; confirmPassword?: string }>({});
  const [isLoading, setIsLoading] = useState(false);
  const [isSubmitted, setIsSubmitted] = useState(false);
  const [tokenError, setTokenError] = useState(false);

  useEffect(() => {
    const tokenParam = searchParams.get("token");
    if (!tokenParam) {
      setTokenError(true);
      showNotification(
        "Invalid or missing reset token. Please request a new password reset.",
        "error",
        3000
      );
    } else {
      setToken(tokenParam);
    }
  }, [searchParams, showNotification]);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { id, value } = e.target;
    setFormData((prev) => ({ ...prev, [id]: value }));
    setFieldErrors((prev) => ({ ...prev, [id]: undefined }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (formData.password !== formData.confirmPassword) {
      setFieldErrors((prev) => ({ ...prev, confirmPassword: "Passwords do not match" }));
      return;
    }

    setFieldErrors({});
    setIsLoading(true);
    try {
      const { data, error, errors } = await resetPassword(token, formData.password);

      if (errors && errors.length > 0) {
        const mapped: { password?: string; confirmPassword?: string } = {};
        errors.forEach((err: { fieldName: string; errorMessage: string }) => {
          const key = err.fieldName || "";
          mapped[key as keyof typeof mapped] = err.errorMessage || "";
        });
        setFieldErrors(mapped);
        return;
      }

      if (error) {
        showNotification(error, "error", 3000);
        return;
      }

      if (data) {
        setIsSubmitted(true);
        showNotification(data?.message, "info", 3000);
      }
    } finally {
      setIsLoading(false);
    }
  };

  if (tokenError) {
    return (
      <div className="h-screen w-full flex flex-col justify-center items-center p-10">
        <h1 className="text-3xl font-bold mb-10">Invalid Reset Link</h1>
        <p className="text-center mb-6 max-w-sm">
          The password reset link is invalid or has expired. Please request a new password reset.
        </p>
        <div className="flex flex-col gap-5 w-full max-w-sm">
          <Button onClick={() => router.push(ROUTES.FORGOT_PASSWORD)}>
            Request New Reset Link
          </Button>
        </div>
      </div>
    );
  }

  if (isSubmitted) {
    return (
      <div className="h-screen w-full flex flex-col justify-center items-center p-10">
        <h1 className="text-3xl font-bold mb-10">Password Reset Complete</h1>
        <p className="text-center mb-6 max-w-sm">
          Your password has been reset successfully. You can now log in with your new password.
        </p>
        <div className="flex flex-col gap-5 w-full max-w-sm">
          <Button onClick={() => router.push(ROUTES.LOGIN)}>
            Go to Login
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="h-screen w-full flex flex-col justify-center items-center p-10">
      <div className="mb-10 space-y-2">
        <h1 className="text-center text-3xl font-bold">Create New Password</h1>
        <p className="text-center text-sm text-gray-500 max-w-sm">
          Enter your new password below. Password must be at least 8 characters long.
        </p>
      </div>
      <form onSubmit={handleSubmit} className="flex flex-col gap-7 w-full max-w-sm">
        <FormInput
          id="password"
          type="password"
          label="New Password"
          value={formData.password}
          onChange={handleInputChange}
          placeholder="Enter your new password"
          isInvalid={!!fieldErrors.password}
          errorMessage={fieldErrors.password}
        />
        <FormInput
          id="confirmPassword"
          type="password"
          label="Confirm Password"
          value={formData.confirmPassword}
          onChange={handleInputChange}
          placeholder="Confirm your new password"
          isInvalid={!!fieldErrors.confirmPassword}
          errorMessage={fieldErrors.confirmPassword}
        />
        <Button type="submit" isLoading={isLoading} loadingText="Resetting...">
          Reset Password
        </Button>
      </form>
    </div>
  );
}

// Main page component with Suspense boundary
export default function ResetPasswordPage() {
  return (
    <Suspense fallback={<LoadingFallback />}>
      <ResetPasswordContent />
    </Suspense>
  );
}
