"use client";

import React, { useState } from "react";
import { useRouter } from "next/navigation";

import { registerUser } from "@/services/auth";
import GoogleSignInButton from "@/components/google-signin";
import Button from "@/components/button";
import { useNotification } from "@/components/notification";
import FormInput from "@/components/form-input";
import { ROUTES } from "@/constants/routes";

const ERROR_NOTIFICATION_DURATION = 5000;
const SUCCESS_NOTIFICATION_DURATION = 2000;
const REDIRECT_DELAY = 2500;

export default function RegisterPage() {
  const router = useRouter();
  const { showNotification } = useNotification();
  const [formData, setFormData] = useState({
    email: "",
    password: "",
    verifyPassword: "",
  });
  const [isLoading, setIsLoading] = useState(false);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { id, value } = e.target;
    setFormData((prev) => ({ ...prev, [id]: value }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (formData.password !== formData.verifyPassword) {
      showNotification("Passwords do not match", "error", ERROR_NOTIFICATION_DURATION);
      return;
    }

    setIsLoading(true);
    try {
      const { data, error } = await registerUser(formData.email, formData.password);

      if (error) {
        showNotification(error, "error", ERROR_NOTIFICATION_DURATION);
        return;
      }

      if (data) {
        showNotification(data.message, "success", SUCCESS_NOTIFICATION_DURATION);
        setTimeout(() => router.replace(ROUTES.LOGIN), REDIRECT_DELAY);
      }
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="h-screen w-full flex flex-col justify-center items-center p-10">
      <h2 className="text-2xl font-bold mb-6">Create your account</h2>
      <form onSubmit={handleSubmit} className="flex flex-col gap-5 w-full max-w-sm">
        <FormInput
          id="email"
          type="email"
          label="Email"
          value={formData.email}
          onChange={handleInputChange}
          placeholder="Enter your email address"
        />
        <FormInput
          id="password"
          type="password"
          label="Password"
          value={formData.password}
          onChange={handleInputChange}
          placeholder="Create a password"
        />
        <FormInput
          id="verifyPassword"
          type="password"
          label="Verify Password"
          value={formData.verifyPassword}
          onChange={handleInputChange}
          placeholder="Re-enter your password"
        />

        <Button type="submit" isLoading={isLoading} loadingText="Signing up...">
          Register
        </Button>
      </form>

      <div className="flex flex-col gap-5 w-sm mt-3">
        <GoogleSignInButton />
        <div className="text-sm text-center">
          Already have an account?{" "}
          <a className="underline" href={ROUTES.LOGIN}>
            Sign in
          </a>
        </div>
      </div>
    </div>
  );
}
