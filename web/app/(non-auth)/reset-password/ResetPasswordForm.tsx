"use client";

import React, { useEffect, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { useForm } from "react-hook-form";

import { resetPassword } from "@/services/auth";
import Button from "@/components/button";
import useNotification from "@/hooks/useNotification";
import TextField from "@/components/text-field";
import { ROUTES } from "@/constants/routes";
import type { FieldError } from "@/shared/types";

type ResetFormValues = { password: string; confirmPassword: string };

export default function ResetPasswordForm() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { showNotification } = useNotification();

  const [token, setToken] = useState("");
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

  const {
    register,
    handleSubmit,
    setError,
    getValues,
    formState: { errors },
  } = useForm<ResetFormValues>({
    defaultValues: { password: "", confirmPassword: "" },
    mode: "onSubmit",
    reValidateMode: "onChange",
  });

  const onSubmit = async (values: ResetFormValues) => {
    if (values.password !== values.confirmPassword) {
      setError("confirmPassword", { type: "validate", message: "Passwords do not match" });
      return;
    }

    const { data, errors, errorFields } = await resetPassword(token, values.password);

    if (errorFields && errorFields.length > 0) {
      errorFields.forEach((err: FieldError) => {
        setError(err.fieldName as keyof ResetFormValues, {
          type: "server",
          message: err.errorMessage || "",
        });
      });
      return;
    }

    if (error) {
      showNotification(error, "error", 3000);
      return;
    }

    if (data) {
      showNotification(data?.message, "info", 3000);
      router.push(ROUTES.LOGIN);
    }
  };

  if (tokenError) {
    return (
      <div className="flex flex-col gap-5 w-full max-w-sm">
        <Button onClick={() => router.push(ROUTES.FORGOT_PASSWORD)}>
          Request New Reset Link
        </Button>
      </div>
    );
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-7 w-full max-w-sm">
      <TextField
        type="password"
        placeholder="Enter your new password"
        xIsRounded
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
        xIsRounded
        xIsInvalid={!!errors.confirmPassword}
        xErrorMessage={errors.confirmPassword?.message}
        {...register("confirmPassword", {
          required: "Please confirm your password",
          validate: (val) => val === getValues("password") || "Passwords do not match",
        })}
      />
      <Button xSize="lg" type="submit">Reset Password</Button>
    </form>
  );
}
