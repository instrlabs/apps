"use client";

import React from "react";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";

import { ROUTES } from "@/constants/routes";
import { registerUser } from "@/services/authentications";
import useNotification from "@/hooks/useNotification";
import Button from "@/components/actions/button";
import TextField from "@/components/inputs/text-field";
import type { FieldError } from "@/shared/types";

type RegisterFormValues = {
  name: string;
  email: string;
  password: string;
  verifyPassword: string;
};

export default function RegisterForm() {
  const router = useRouter();
  const { showNotification } = useNotification();

  const {
    register,
    handleSubmit,
    setError,
    getValues,
    formState: { errors },
  } = useForm<RegisterFormValues>({
    defaultValues: { name: "", email: "", password: "", verifyPassword: "" },
    mode: "onSubmit",
    reValidateMode: "onChange",
  });

  const onSubmit = async (values: RegisterFormValues) => {
    if (values.password !== values.verifyPassword) {
      setError("verifyPassword", { type: "validate", message: "Passwords do not match" });
      return;
    }

    const { success, message, data, errors } = await registerUser({
      name: values.name,
      email: values.email,
      password: values.password
    });

    if (errors && errors.length > 0) {
      errors.forEach((err: FieldError) => {
        setError(err.fieldName as keyof RegisterFormValues, {
          type: "server",
          message: err.errorMessage,
        });
      });
      return;
    } else if (!success) {
      showNotification({ title: "Error", message, type: "error", duration: 3000 });
    } else if (data) {
      router.replace(ROUTES.LOGIN);
    }
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-7 w-full max-w-sm">
      <TextField
        type="name"
        placeholder="Enter your name"
        xIsRounded
        xIsInvalid={!!errors.name}
        xErrorMessage={errors.name?.message}
        {...register("name", { required: "Name is required" })}
      />
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

      <Button xSize="lg" type="submit">Sign up</Button>
    </form>
  );
}
