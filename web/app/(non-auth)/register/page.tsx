"use client";

import React, { useState } from "react";
import { useRouter } from "next/navigation";

import { registerUser } from "@/services/auth";
import GoogleSignInButton from "@/components/google-signin";
import Button from "@/components/button";
import { useNotification } from "@/components/notification";
import FormInput from "@/components/form-input";
import { ROUTES } from "@/constants/routes";

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
  const [fieldErrors, setFieldErrors] = useState<{ email?: string; password?: string; verifyPassword?: string }>({});

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { id, value } = e.target;
    setFormData((prev) => ({ ...prev, [id]: value }));
    // Clear field-specific error on change
    setFieldErrors((prev) => ({ ...prev, [id]: undefined }));
  };

  const passwordsMismatch =
    formData.password.length > 0 &&
    formData.verifyPassword.length > 0 &&
    formData.password !== formData.verifyPassword;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    setFieldErrors({});

    if (formData.password !== formData.verifyPassword) {
      setFieldErrors({ verifyPassword: "Passwords do not match" });
      return;
    }

    setIsLoading(true);
    try {
      const { data, error, errors } = await registerUser(formData.email, formData.password);

      if (errors && errors.length > 0) {
        const mapped: { email?: string; password?: string; verifyPassword?: string } = {};
        errors.forEach((err: { fieldName?: string; errorMessage?: string }) => {
          const key = (err.fieldName || "").toString();
          if (key === "email" || key === "password" || key === "verifyPassword") {
            mapped[key as keyof typeof mapped] = err.errorMessage || "";
          }
        });
        setFieldErrors(mapped);
        return;
      }

      if (error) {
        showNotification(error, "error", 3000);
        return;
      }

      if (data) {
        router.replace(ROUTES.LOGIN);
      }
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="h-screen w-full flex flex-col justify-center items-center p-10">
      <h1 className="text-3xl font-bold mb-10">Create your account</h1>
      <form onSubmit={handleSubmit} className="flex flex-col gap-7 w-full max-w-sm">
        <FormInput
          id="email"
          type="email"
          label="Email"
          value={formData.email}
          onChange={handleInputChange}
          placeholder="Enter your email address"
          isInvalid={!!fieldErrors.email}
          errorMessage={fieldErrors.email}
        />
        <FormInput
          id="password"
          type="password"
          label="Password"
          value={formData.password}
          onChange={handleInputChange}
          placeholder="Create a password"
          isInvalid={!!fieldErrors.password}
          errorMessage={fieldErrors.password}
        />
        <FormInput
          id="verifyPassword"
          type="password"
          label="Verify Password"
          value={formData.verifyPassword}
          onChange={handleInputChange}
          placeholder="Re-enter your password"
          isInvalid={!!fieldErrors.verifyPassword || passwordsMismatch}
          errorMessage={fieldErrors.verifyPassword}
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
