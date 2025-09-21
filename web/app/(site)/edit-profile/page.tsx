"use client";

import React, { useEffect } from "react";
import { useForm } from "react-hook-form";
import { useRouter } from "next/navigation";

import TextField from "@/components/inputs/text-field";
import Button from "@/components/actions/button";
import ButtonIcon from "@/components/actions/button-icon";
import ChevronLeftIcon from "@/components/icons/chevron-left";
import useNotification from "@/hooks/useNotification";
import { useProfile } from "@/hooks/useProfile";
import { updateProfile } from "@/services/auth";

type EditProfileFormValues = {
  name: string;
};


export default function EditProfilePage() {
  const router = useRouter();
  const { showNotification } = useNotification();
  const { profile, setProfile } = useProfile();

  const {
    register,
    handleSubmit,
    setError,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<EditProfileFormValues>({
    defaultValues: { name: "" },
    mode: "onSubmit",
    reValidateMode: "onChange",
  });


  useEffect(() => {
    if (profile) {
      reset({ name: profile.name || "" });
    }
  }, [profile, reset]);

  const onSubmit = async (values: EditProfileFormValues) => {
    const { errors, errorFields } = await updateProfile(values.name);

    if (errorFields && errorFields.length > 0) {
      errorFields.forEach((err: { fieldName: string; errorMessage: string }) => {
        if (err.fieldName === "name")
          setError(err.fieldName as keyof EditProfileFormValues, { type: "server", message: err.errorMessage || "" });
      });
      return;
    }

    if (error) {
      showNotification({ title: "Error", message: error, type: "error", duration: 3000 });
      return;
    }

    // Regardless of whether update returns user data, refresh the profile to stay in sync with auth-service
    try {
      // Lazy import to avoid circular issues
      const { getProfile: fetchProfile } = await import("@/services/auth");
      const res = await fetchProfile();
      if (res && !res.errors && res.data) {
        setProfile(res.data.data.user);
      }
    } catch {}

    showNotification({ title: "Info", message: "Profile updated successfully", type: "info", duration: 2500 });
    router.back();
  };


  return (
    <div className="w-full h-full flex flex-col">
      <div className="p-6 flex items-center gap-3">
        <ButtonIcon xSize="sm" xColor="primary" onClick={() => window.history.back()}>
          <ChevronLeftIcon />
        </ButtonIcon>
        <h1 className="text-xl font-bold">Edit Profile</h1>
      </div>
      <form onSubmit={handleSubmit(onSubmit)} className="p-6 flex flex-col gap-6">
        <div>
          <label className="block text-sm font-medium mb-2">Name</label>
          <TextField
            type="text"
            placeholder="Your name"
            xIsInvalid={!!errors.name}
            xErrorMessage={errors.name?.message}
            {...register("name", {
              required: "Name is required",
              minLength: { value: 2, message: "Name must be at least 2 characters" },
            })}
          />
        </div>
        <div className="flex justify-end">
          <Button xSize="md" type="submit">{isSubmitting ? "Saving..." : "Save Changes"}</Button>
        </div>
      </form>
    </div>
  );
}
