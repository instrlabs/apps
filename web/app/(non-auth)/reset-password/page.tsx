import React from "react";
import type { Metadata } from "next";

import ResetPasswordForm from "./ResetPasswordForm";

export const metadata: Metadata = {
  title: "Reset password",
  description: "Create a new password to access your account.",
  openGraph: {
    title: "Reset password",
    description: "Create a new password to access your account.",
    type: "website",
  },
  twitter: {
    card: "summary",
    title: "Reset password",
    description: "Create a new password to access your account.",
  },
  robots: {
    index: false,
    follow: false,
  },
};

export default function ResetPasswordPage() {
  return (
    <div className="h-screen w-full flex flex-col justify-center items-center p-10">
      <h1 className="text-3xl font-bold mb-15">Create New Password</h1>
      <ResetPasswordForm />
    </div>
  );
}
