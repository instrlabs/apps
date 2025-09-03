"use client";

import React from "react";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";

import { requestPasswordReset } from "@/services/auth";
import Button from "@/components/button";
import useNotification from "@/hooks/useNotification";
import { ROUTES } from "@/constants/routes";
import TextField from "@/components/text-field";
import LinkText from "@/components/link-text";
import type { FieldError } from "@/shared/types";

export type ForgotFormValues = { email: string };

export default function ForgotPasswordForm() {
  const router = useRouter();
  const { showNotification } = useNotification();

  const {
    register,
    handleSubmit,
    setError,
    formState: { errors },
  } = useForm<ForgotFormValues>({
    defaultValues: { email: "" },
    mode: "onSubmit",
    reValidateMode: "onChange",
  });

  const onSubmit = async (values: ForgotFormValues) => {
    const { success, message, errors } = await requestPasswordReset(values.email);

    if (errors && errors.length > 0) {
      errors.forEach((err: FieldError) => {
        setError(err.fieldName as keyof ForgotFormValues, {
          type: "server",
          message: err.errorMessage || "",
        });
      });
    } else if (!success) {
      showNotification(message, "error", 3000);
    } else {
      router.push(`${ROUTES.FORGOT_PASSWORD}/check`);
    }
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-7 w-full max-w-sm">
      <TextField
        type="email"
        placeholder="Enter your email"
        xIsRounded
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
      <Button xSize="lg" type="submit">Send reset password</Button>
      <div className="text-sm text-center">
        Remember your password? <LinkText href={ROUTES.LOGIN}>Sign in</LinkText>
      </div>
    </form>
  );
}
