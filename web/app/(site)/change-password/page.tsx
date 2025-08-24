"use client";

import React from "react";
import { useForm } from "react-hook-form";

import TextField from "@/components/text-field";
import Button from "@/components/button";
import { useNotification } from "@/components/notification";
import { changePassword } from "@/services/auth";
import ButtonIcon from "@/components/button-icon";
import ChevronLeftIcon from "@/components/icons/chevron-left";

type ChangePasswordFormValues = {
  current_password: string;
  new_password: string;
  confirm_password: string;
};

export default function ChangePasswordPage() {
  const { showNotification } = useNotification();

  const {
    register,
    handleSubmit,
    setError,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<ChangePasswordFormValues>({
    defaultValues: { current_password: "", new_password: "", confirm_password: "" },
    mode: "onSubmit",
    reValidateMode: "onChange",
  });

  const onSubmit = async (values: ChangePasswordFormValues) => {
    if (values.new_password !== values.confirm_password) {
      setError("confirm_password", { type: "validate", message: "Passwords do not match" });
      return;
    }

    const { data, error, errorFields } = await changePassword(values.current_password, values.new_password);

    if (errorFields && errorFields.length > 0) {
      const fieldError = errorFields[0];
      showNotification(fieldError?.errorMessage || "Failed to change password", "error", 3000);
      return;
    }

    if (error) {
      if (/current password/i.test(error)) {
        setError("current_password", { type: "server", message: error });
      } else {
        showNotification(error, "error", 3000);
      }
      return;
    }

    if (data) {
      showNotification("Password changed successfully", "info", 2500);
      reset({ current_password: "", new_password: "", confirm_password: "" });
    }
  };

  return (
    <div className="w-full h-full flex flex-col">
      <div className="p-6 flex items-center gap-3">
        <ButtonIcon xSize="sm" xColor="primary" onClick={() => window.history.back()}>
          <ChevronLeftIcon />
        </ButtonIcon>
        <h1 className="text-xl font-bold">Change Password</h1>
      </div>
      <form onSubmit={handleSubmit(onSubmit)} className="p-6 flex flex-col gap-6">
        <div>
          <label className="block text-sm font-medium mb-2">Current Password</label>
          <TextField
            type="password"
            placeholder="Current password"
            xSize="md"
            xIsInvalid={!!errors.current_password}
            xErrorMessage={errors.current_password?.message}
            {...register("current_password", { required: "Current password is required" })}
          />
        </div>
        <div>
          <label className="block text-sm font-medium mb-2">New Password</label>
          <TextField
            type="password"
            placeholder="New password"
            xSize="md"
            xIsInvalid={!!errors.new_password}
            xErrorMessage={errors.new_password?.message}
            {...register("new_password", { required: "New password is required", minLength: { value: 6, message: "Password must be at least 6 characters" } })}
          />
        </div>
        <div>
          <label className="block text-sm font-medium mb-2">Confirm New Password</label>
          <TextField
            type="password"
            placeholder="Confirm new password"
            xSize="md"
            xIsInvalid={!!errors.confirm_password}
            xErrorMessage={errors.confirm_password?.message}
            {...register("confirm_password", { required: "Please confirm your new password" })}
          />
        </div>
        <div className="flex justify-end">
          <Button xSize="md" type="submit">{isSubmitting ? "Changing..." : "Change Password"}</Button>
        </div>
      </form>
    </div>
  );
}
