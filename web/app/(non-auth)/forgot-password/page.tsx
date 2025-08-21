"use client";

import React, { useState } from "react";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";

import { requestPasswordReset } from "@/services/auth";
import Button from "@/components/button";
import { useNotification } from "@/components/notification";
import { ROUTES } from "@/constants/routes";
import TextField from "@/components/text-field";
import LinkText from "@/components/link-text";

type ForgotFormValues = { email: string };

export default function ForgotPasswordPage() {
  const router = useRouter();
  const { showNotification } = useNotification();
  const [isSubmitted, setIsSubmitted] = useState(false);

  const { register, handleSubmit, setError, watch, formState: { errors } } = useForm<ForgotFormValues>({
    defaultValues: { email: "" },
    mode: "onSubmit",
    reValidateMode: "onChange",
  });

  const onSubmit = async (values: ForgotFormValues) => {
    const { data, error, errorFields } = await requestPasswordReset(values.email);

    if (errorFields && errorFields.length > 0) {
      errorFields.forEach((err: { fieldName: string; errorMessage: string }) => {
        setError(err.fieldName as keyof ForgotFormValues, { type: "server", message: err.errorMessage || "" });
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

  const emailValue = watch("email");

  if (isSubmitted) {
    return (
      <div className="h-screen w-full flex flex-col justify-center items-center p-10">
        <h1 className="text-3xl font-bold mb-10">Check your email</h1>
        <p className="text-center mb-6 max-w-sm">
          We&#39;ve sent password reset instructions to <strong>{emailValue}</strong>.
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
      <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-7 w-full max-w-sm">
        <TextField
          type="email"
          placeholder="Enter your email address"
          xIsInvalid={!!errors.email}
          xErrorMessage={errors.email?.message}
          {...register("email", {
            required: "Email is required",
            pattern: {
              value: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
              message: "Enter a valid email address",
            },
          })}
        />
        <Button type="submit">Send reset instructions</Button>
        <div className="text-sm text-center">
          Remember your password?{" "}
          <LinkText href={ROUTES.LOGIN}>Sign in</LinkText>
        </div>
      </form>
    </div>
  );
}
