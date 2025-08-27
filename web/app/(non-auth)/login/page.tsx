import React from "react";
import type { Metadata } from "next";

import { ROUTES } from "@/constants/routes";
import GoogleSignInButton from "@/components/google-signin";
import LinkText from "@/components/link-text";
import LoginForm from "./LoginForm";

export const metadata: Metadata = {
  title: "Log in",
  description: "Log in to your account to access your dashboard and manage your settings.",
  openGraph: {
    title: "Log in",
    description: "Log in to your account to access your dashboard and manage your settings.",
    type: "website",
  },
  twitter: {
    card: "summary",
    title: "Log in",
    description: "Log in to your account to access your dashboard and manage your settings.",
  },
  robots: {
    index: false,
    follow: false,
  },
};

export default function LoginPage() {
  return (
    <div className="h-screen w-full flex flex-col justify-center items-center p-10">
      <h1 className="text-3xl font-bold mb-15">Log in to your account</h1>
      <LoginForm />
      <div className="flex flex-col gap-5 w-sm mt-3">
        <GoogleSignInButton />
        <div className="text-sm text-center">
          Don&apos;t have an account?{" "}
          <LinkText href={ROUTES.REGISTER}>Sign up</LinkText>
        </div>
      </div>
    </div>
  );
}
