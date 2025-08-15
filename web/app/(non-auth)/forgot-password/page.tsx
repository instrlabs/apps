"use client";

import React, { useState } from "react";
import { useRouter } from "next/navigation";

import { requestPasswordReset } from "@/services/auth";
import Button from "@/components/button";
import { useNotification } from "@/components/notification";
import { ROUTES } from "@/constants/routes";
import FormInput from "@/components/form-input";

export default function ForgotPasswordPage() {
  const router = useRouter();
  const { showNotification } = useNotification();
  const [formData, setFormData] = useState({ email: "" });
  const [isLoading, setIsLoading] = useState(false);
  const [isSubmitted, setIsSubmitted] = useState(false);
  const [fieldErrors, setFieldErrors] = useState<{ email?: string }>({});

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { id, value } = e.target;
    setFormData((prev) => ({ ...prev, [id]: value }));
    setFieldErrors((prev) => ({ ...prev, [id]: undefined }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    setFieldErrors({});
    setIsLoading(true);
    try {
      const { data, error, errors } = await requestPasswordReset(formData.email);

      if (errors && errors.length > 0) {
        const mapped: { email?: string } = {};
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

  if (isSubmitted) {
    return (
      <div className="h-screen w-full flex flex-col justify-center items-center p-10">
        <h1 className="text-3xl font-bold mb-10">Check your email</h1>
        <p className="text-center mb-6 max-w-sm">
          We&#39;ve sent password reset instructions to <strong>{formData.email}</strong>.
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
      <div className="mb-10 space-y-2">
        <h1 className="text-center text-3xl font-bold">Reset your password</h1>
        <p className="text-center text-sm text-gray-500 max-w-sm">
          Enter your email address and we&#39;ll send you instructions to reset your password.
        </p>
      </div>
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
