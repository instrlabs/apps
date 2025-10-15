import React from "react";
import type { Metadata } from "next";

import GoogleSignInButton from "@/components/actions/google-signin";
import LoginForm from "./LoginForm";

export const metadata: Metadata = {
  title: "Login - Labs",
  description: "",
};

export default function LoginPage() {
  return (
    <div className="h-screen w-screen flex items-center justify-center">
      <div className="w-full max-w-md flex flex-col gap-6 p-10 mx-auto">
        <LoginForm />
      </div>
    </div>
  );
}
