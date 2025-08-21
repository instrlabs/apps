"use client";

import React, { useState, useEffect, Suspense } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { useForm } from "react-hook-form";

import { resetPassword } from "@/services/auth";
import Button from "@/components/button";
import { useNotification } from "@/components/notification";
import TextField from "@/components/text-field";
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
  const [tokenError, setTokenError] = useState(false);
  const [isSubmitted, setIsSubmitted] = useState(false);

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

  type ResetFormValues = { password: string; confirmPassword: string };

  const { register, handleSubmit, setError, getValues, formState: { errors } } = useForm<ResetFormValues>({
    defaultValues: { password: "", confirmPassword: "" },
    mode: "onSubmit",
    reValidateMode: "onChange",
  });

  const onSubmit = async (values: ResetFormValues) => {
    if (values.password !== values.confirmPassword) {
      setError("confirmPassword", { type: "validate", message: "Passwords do not match" });
      return;
    }

    const { data, error, errorFields } = await resetPassword(token, values.password);

    if (errorFields && errorFields.length > 0) {
      errorFields.forEach((err: { fieldName: string; errorMessage: string }) => {
        setError(err.fieldName as keyof ResetFormValues, { type: "server", message: err.errorMessage || "" });
      });
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
      <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-7 w-full max-w-sm">
        <TextField
          type="password"
          placeholder="Enter your new password"
          xIsInvalid={!!errors.password}
          xErrorMessage={errors.password?.message}
          {...register("password", {
            required: "Password is required",
            minLength: { value: 6, message: "Password must be at least 6 characters" },
          })}
        />
        <TextField
          type="password"
          placeholder="Confirm your new password"
          xIsInvalid={!!errors.confirmPassword}
          xErrorMessage={errors.confirmPassword?.message}
          {...register("confirmPassword", {
            required: "Please confirm your password",
            validate: (val) => val === getValues("password") || "Passwords do not match",
          })}
        />
        <Button type="submit">Reset Password</Button>
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
