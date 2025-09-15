import React from "react";
import type { Metadata } from "next";

import { ROUTES } from "@/constants/routes";
import GoogleSignInButton from "@/components/actions/google-signin";
import LinkText from "@/components/actions/link-text";
import RegisterForm from "./RegisterForm";

export const metadata: Metadata = {
  title: "Sign up",
  description: "Create your account to start using the app and manage your settings.",
  openGraph: {
    title: "Sign up",
    description: "Create your account to start using the app and manage your settings.",
    type: "website",
  },
  twitter: {
    card: "summary",
    title: "Sign up",
    description: "Create your account to start using the app and manage your settings.",
  },
  robots: {
    index: false,
    follow: false,
  },
};

export default function RegisterPage() {
  return (
    <div className="h-screen w-full flex flex-col justify-center items-center p-10">
      <h1 className="text-3xl font-bold mb-15">Create your account</h1>
      <RegisterForm />
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
