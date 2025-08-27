"use client";

import React from "react";
import Button from "@/components/button";
import { useRouter } from "next/navigation";
import { ROUTES } from "@/constants/routes";

export default function ForgotPasswordCheckPage() {
  const router = useRouter();

  return (
    <div className="w-full h-screen flex justify-center items-center">
      <div className="flex flex-col justify-center items-center p-10 border border-border rounded-xl shadow-primary">
        <h2 className="text-2xl font-bold mb-6">Check your email</h2>
        <p className="text-center mb-6 max-w-sm">
          We have sent password reset instructions to your email address. Please check your inbox and follow the link in the email. If you don&apos;t see it, check your spam folder.
        </p>
        <div className="flex flex-col gap-5 w-full max-w-sm">
          <Button xSize="lg" onClick={() => router.push(ROUTES.LOGIN)}>
            Return to login
          </Button>
        </div>
      </div>
    </div>
  );
}
