"use client";

import React from "react";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";

import { registerUser } from "@/services/auth";
import GoogleSignInButton from "@/components/google-signin";
import Button from "@/components/button";
import { useNotification } from "@/components/notification";
import TextField from "@/components/text-field";
import LinkText from "@/components/link-text";
import { ROUTES } from "@/constants/routes";

type RegisterFormValues = {
  email: string;
  password: string;
  verifyPassword: string;
};

export default function RegisterPage() {
  const router = useRouter();
  const { showNotification } = useNotification();

  const {
    register,
    handleSubmit,
    setError,
    getValues,
    formState: { errors },
  } = useForm<RegisterFormValues>({
    defaultValues: { email: "", password: "", verifyPassword: "" },
    mode: "onSubmit",
    reValidateMode: "onChange",
  });

  const onSubmit = async (values: RegisterFormValues) => {
    if (values.password !== values.verifyPassword) {
      setError("verifyPassword", { type: "validate", message: "Passwords do not match" });
      return;
    }

    const { data, error, errorFields } = await registerUser(values.email, values.password);

    if (errorFields && errorFields.length > 0) {
      errorFields.forEach((err: { fieldName: string; errorMessage: string }) => {
        setError(err.fieldName as keyof RegisterFormValues, { type: "server", message: err.errorMessage || "" });
      });
      return;
    }

    if (error) {
      showNotification(error, "error", 3000);
      return;
    }

    if (data) {
      router.replace(ROUTES.LOGIN);
    }
  };

  return (
    <div className="h-screen w-full flex flex-col justify-center items-center p-10">
      <h1 className="text-3xl font-bold mb-10">Create your account</h1>
      <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-7 w-full max-w-sm">
        <TextField
          type="email"
          placeholder="Enter your email address"
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
        <TextField
          type="password"
          placeholder="Create a password"
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
          placeholder="Re-enter your password"
          xIsRounded
          xIsInvalid={!!errors.verifyPassword}
          xErrorMessage={errors.verifyPassword?.message}
          {...register("verifyPassword", {
            required: "Please verify your password",
            validate: (val) => val === getValues("password") || "Passwords do not match",
          })}
        />
        <Button type="submit">Register</Button>
      </form>

      <div className="flex flex-col gap-5 w-sm mt-3">
        <GoogleSignInButton />
        <div className="text-sm text-center">
          Already have an account?{" "}
          <LinkText href={ROUTES.LOGIN}>Sign in</LinkText>
        </div>
      </div>
    </div>
  );
}
