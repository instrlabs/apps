import React from "react";
import type { Metadata } from "next";

import LoginForm from "./LoginForm";

export const metadata: Metadata = {
  title: "Login - Labs",
  description: "",
};

export default function LoginPage() {
  return (
    <div className="h-screen w-screen flex items-center justify-center">
      <div className="w-full max-w-md mx-auto">
        <LoginForm />
      </div>
    </div>
  );
}
