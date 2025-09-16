"use client";

import React from "react";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";

import { ROUTES } from "@/constants/routes";
import useNotification from "@/hooks/useNotification";
import Button from "@/components/actions/button";
import TextField from "@/components/inputs/text-field";
import LinkText from "@/components/actions/link-text";
import { loginUser } from "@/services/authentications";
import { SOMETHING_WENT_WRONG } from "@/constants/errors";
import type { FieldError } from "@/shared/types";

type LoginFormValues = {
  email: string;
  password: string;
};

export default function LoginForm() {
  const router = useRouter();
  const { showNotification } = useNotification();

  const {
    register,
    handleSubmit,
    setError,
    formState: { errors },
  } = useForm<LoginFormValues>({
    defaultValues: { email: "", password: "" },
    mode: "onSubmit",
    reValidateMode: "onChange",
  });

  const onSubmit = async (values: LoginFormValues) => {
    try {
      const { success, message, errors } = await loginUser({
        email: values.email,
        password: values.password,
      });

      if (errors && errors?.length > 0) {
        errors.forEach((err: FieldError) => {
          setError(err.fieldName as keyof LoginFormValues, {
            type: "server",
            message: err.errorMessage,
          });
        });
      } else if (!success) {
        showNotification({ title: "Error", message, type: "error", duration: 3000 });
      } else {
        router.push(ROUTES.HOME);
      }
    } catch {
      showNotification({ title: "Error", message: SOMETHING_WENT_WRONG, type: "error", duration: 3000 });
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

      <div className="space-y-1">
        <TextField
          type="password"
          placeholder="Enter your password"
          xIsRounded
          xIsInvalid={!!errors.password}
          xErrorMessage={errors.password?.message}
          {...register("password", {
            required: "Password is required",
            minLength: { value: 6, message: "Password must be at least 6 characters" },
          })}
        />
        <div className="text-right">
          <LinkText href={ROUTES.FORGOT_PASSWORD}>Forgot password?</LinkText>
        </div>
      </div>

      <Button xSize="lg" type="submit">Sign in</Button>
    </form>
  );
}
