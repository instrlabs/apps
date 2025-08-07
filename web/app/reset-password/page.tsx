"use client";

import React, { useState, useEffect } from "react";
import { useRouter, useSearchParams } from "next/navigation";

import { resetPassword } from "@/services/auth";
import Button from "@/components/button";
import { useNotification } from "@/components/notification";
import {
  containerStyles,
  headingStyles,
  paragraphStyles,
  formContainerStyles,
  inputGroupStyles,
  inputStyles,
  labelStyles,
  NOTIFICATION_DURATION
} from "@/components/ui-styles";

// Constants
const MIN_PASSWORD_LENGTH = 8;

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
    const tokenParam = searchParams.get("token");
    if (!tokenParam) {
      setTokenError(true);
      showNotification(
        "Invalid or missing reset token. Please request a new password reset.",
        "error",
        NOTIFICATION_DURATION
      );
    } else {
      setToken(tokenParam);
    }
  }, [searchParams, showNotification]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (password !== confirmPassword) {
      showNotification("Passwords do not match", "error", NOTIFICATION_DURATION);
      return;
    }

    if (password.length < MIN_PASSWORD_LENGTH) {
      showNotification(
        `Password must be at least ${MIN_PASSWORD_LENGTH} characters long`, 
        "error", 
        NOTIFICATION_DURATION
      );
      return;
    }

    setIsLoading(true);

    const { data, error } = await resetPassword(token, password);

    if (error) {
      showNotification(error, "error", NOTIFICATION_DURATION);
    }

    if (data) {
      setIsSubmitted(true);
      showNotification(data?.message, "success", NOTIFICATION_DURATION);
    }

    setIsLoading(false);
  };

  // Handle input changes
  const handlePasswordChange = (e: React.ChangeEvent<HTMLInputElement>) => setPassword(e.target.value);
  const handleConfirmPasswordChange = (e: React.ChangeEvent<HTMLInputElement>) => setConfirmPassword(e.target.value);
  
  if (tokenError) {
    return (
      <div className={containerStyles}>
        <h2 className={headingStyles}>Invalid Reset Link</h2>
        <p className={paragraphStyles}>
          The password reset link is invalid or has expired. Please request a new password reset.
        </p>
        <div className={formContainerStyles}>
          <Button onClick={() => router.push("/forgot-password")}>
            Request New Reset Link
          </Button>
        </div>
      </div>
    );
  }

  if (isSubmitted) {
    return (
      <div className={containerStyles}>
        <h2 className={headingStyles}>Password Reset Complete</h2>
        <p className={paragraphStyles}>
          Your password has been reset successfully. You can now log in with your new password.
        </p>
        <div className={formContainerStyles}>
          <Button onClick={() => router.push("/login")}>
            Go to Login
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className={containerStyles}>
      <h2 className={headingStyles}>Create New Password</h2>
      <p className={paragraphStyles}>
        Enter your new password below. Password must be at least {MIN_PASSWORD_LENGTH} characters long.
      </p>
      <form onSubmit={handleSubmit} className={formContainerStyles}>
        <div className={inputGroupStyles}>
          <label htmlFor="password" className={labelStyles}>
            New Password
          </label>
          <input
            id="password"
            type="password"
            placeholder="Enter your new password"
            className={inputStyles}
            value={password}
            onChange={handlePasswordChange}
            required
            minLength={MIN_PASSWORD_LENGTH}
          />
        </div>
        <div className={inputGroupStyles}>
          <label htmlFor="confirmPassword" className={labelStyles}>
            Confirm Password
          </label>
          <input
            id="confirmPassword"
            type="password"
            placeholder="Confirm your new password"
            className={inputStyles}
            value={confirmPassword}
            onChange={handleConfirmPasswordChange}
            required
            minLength={MIN_PASSWORD_LENGTH}
          />
        </div>
        <Button type="submit" isLoading={isLoading} loadingText="Resetting...">
          Reset Password
        </Button>
      </form>
    </div>
  );
}