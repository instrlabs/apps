"use client";

import React, { useState } from "react";
import { useRouter } from "next/navigation";
import { loginUser } from "@/services/auth";
import GoogleSignInButton from "@/components/google-signin";
import Button from "@/components/button";
import { useNotification } from "@/components/notification";
import FormInput from "@/components/form-input";
import { ROUTES } from "@/constants/routes";

export default function LoginPage() {
  const router = useRouter();
  const { showNotification } = useNotification();
  const [formData, setFormData] = useState({
    email: "",
    password: "",
  });
  const [isLoading, setIsLoading] = useState(false);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { id, value } = e.target;
    setFormData((prev) => ({ ...prev, [id]: value }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    setIsLoading(true);

    const { data, error } = await loginUser(formData.email, formData.password);

    if (error) showNotification(error, "error", 5000);
    else if (data) router.push(ROUTES.HOME);

    setIsLoading(false);
  };

  return (
    <div className="h-screen w-full flex flex-col justify-center items-center p-10">
      <h2 className="text-2xl font-bold mb-6">Log in to your account</h2>
      <form onSubmit={handleSubmit} className="flex flex-col gap-5 w-full max-w-sm">
        <FormInput
          id="email"
          type="email"
          label="Email"
          value={formData.email}
          onChange={handleInputChange}
          placeholder="Enter your email address"
        />

        <div className="space-y-1">
          <FormInput
            id="password"
            type="password"
            label="Password"
            value={formData.password}
            onChange={handleInputChange}
            placeholder="Enter your password"
          />
          <div className="text-right">
            <a
              className="text-sm text-blue-600 hover:underline"
              href={ROUTES.FORGOT_PASSWORD}
            >
              Forgot password?
            </a>
          </div>
        </div>

        <Button type="submit" isLoading={isLoading} loadingText="Signing in...">
          Sign in
        </Button>
      </form>

      <div className="flex flex-col gap-5 w-sm">
        <GoogleSignInButton />
        <div className="text-sm text-center">
          Don&#39;t have an account?{" "}
          <a className="underline" href={ROUTES.REGISTER}>
            Sign up
          </a>
        </div>
      </div>
    </div>
  );
}
