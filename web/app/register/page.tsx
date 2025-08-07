"use client";

import React, { useState } from "react";
import { useRouter } from "next/navigation";

import { registerUser } from "@/services/auth";
import GoogleSignInButton from "@/components/google-signin";
import Button from "@/components/button";
import { useNotification } from "@/components/notification";
import FormInput from "@/components/form-input";
import { ROUTES } from "@/constants/routes";
import {
  containerStyles,
  headingStyles,
  formContainerStyles,
  NOTIFICATION_DURATION
} from "@/components/ui-styles";

const ERROR_NOTIFICATION_DURATION = NOTIFICATION_DURATION;
const SUCCESS_NOTIFICATION_DURATION = 2000;
const REDIRECT_DELAY = 2500;

export default function RegisterPage() {
  const router = useRouter();
  const { showNotification } = useNotification();
  const [formData, setFormData] = useState({
    email: "",
    password: "",
  });
  const [isLoading, setIsLoading] = useState(false);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { id, value } = e.target;
    setFormData((prev) => ({ ...prev, [id]: value }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);

    const { data, error } = await registerUser(formData.email, formData.password);

    if (error) {
      showNotification(error, "error", ERROR_NOTIFICATION_DURATION);
      return;
    }

    if (data) {
      showNotification(data.message, "success", SUCCESS_NOTIFICATION_DURATION);
      setTimeout(() => router.replace(ROUTES.LOGIN), REDIRECT_DELAY);
    }

    setIsLoading(false);
  }

  return (
    <div className={containerStyles}>
      <h2 className={headingStyles}>Create your account</h2>
      <form onSubmit={handleSubmit} className={formContainerStyles}>
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
        <Button type="submit" isLoading={isLoading} loadingText="Signing up...">
          Sign up
        </Button>
      </form>
      <div className="flex flex-col gap-5 w-sm mt-5">
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