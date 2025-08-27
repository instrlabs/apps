import React from "react";
import type { Metadata } from "next";

import ForgotPasswordForm from "./ForgotPasswordForm";

export const metadata: Metadata = {
  title: "Forgot password",
  description: "Enter your email address and we'll send you instructions to reset your password.",
  openGraph: {
    title: "Forgot password",
    description: "Enter your email address and we'll send you instructions to reset your password.",
    type: "website",
  },
  twitter: {
    card: "summary",
    title: "Forgot password",
    description: "Enter your email address and we'll send you instructions to reset your password.",
  },
  robots: {
    index: false,
    follow: false,
  },
};

export default function ForgotPasswordPage() {
  return (
    <div className="h-screen w-full flex flex-col justify-center items-center p-10">
      <h1 className="text-3xl font-bold mb-15">Forgot password</h1>
      <ForgotPasswordForm />
    </div>
  );
}
