"use client";

import React, { useState } from "react";
import { useRouter } from "next/navigation";
import { loginUser } from "@/services/auth";
import GoogleSignInButton from "@/components/google-signin";
import Button from "@/components/button";
import { useNotification } from "@/components/notification";
import TextField from "@/components/text-field";
import LinkText from "@/components/link-text";
import { ROUTES } from "@/constants/routes";
import { useForm } from "react-hook-form";

type LoginFormValues = {
  email: string;
  password: string;
};

export default function LoginPage() {
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
    const { data, error, errorFields } = await loginUser(values.email, values.password);

    if (errorFields && errorFields.length > 0) {
      errorFields.forEach((err: { fieldName: string; errorMessage: string }) => {
        setError(err.fieldName as keyof LoginFormValues, { type: "server", message: err.errorMessage || "" });
      });
      return;
    }

    if (error) {
      showNotification(error, "error", 3000);
      return;
    }

    if (data) {
      router.push(ROUTES.HOME);
    }
  };

  return (
    <div className="h-screen w-full flex flex-col justify-center items-center p-10">
      <h1 className="text-3xl font-bold mb-15">Log in to your account</h1>
      <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-7 w-full max-w-sm">
        <TextField
          type="email"
          placeholder="Enter your email"
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
            xIsInvalid={!!errors.password}
            xErrorMessage={errors.password?.message}
            {...register("password", {
              required: "Password is required",
              minLength: { value: 6, message: "Password must be at least 6 characters" },
            })}
          />
          <div className="text-right">
            <LinkText href={ROUTES.FORGOT_PASSWORD}>
              Forgot password?
            </LinkText>
          </div>
        </div>

        <Button type="submit">Sign in</Button>
      </form>

      <div className="flex flex-col gap-5 w-sm mt-3">
        <GoogleSignInButton />
        <div className="text-sm text-center">
          Don&#39;t have an account?{" "}
          <LinkText href={ROUTES.REGISTER}>
            Sign up
          </LinkText>
        </div>
      </div>
    </div>
  );
}
