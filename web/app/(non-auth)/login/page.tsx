import React from "react";
import type { Metadata } from "next";

import { ROUTE_REGISTER } from "@/constants/routes";
import GoogleSignInButton from "@/components/actions/google-signin";
import LinkText from "@/components/actions/link-text";
import LoginForm from "./LoginForm";
import Button from "@/components/actions/button";

export const metadata: Metadata = {
  title: "Login - Labs",
  description: "",
};

export default function LoginPage() {
  return (
    <div className="h-screen w-screen flex items-center justify-center">
      <div className="w-full max-w-md flex flex-col gap-6 p-10 mx-auto">
        <LoginForm />
        <hr/>
        <GoogleSignInButton />
        <Button xVariant="transparent">
          <p className="font-light">
            Don&apos;t have an account? {" "}
            <LinkText href={ROUTE_REGISTER}>Sign up</LinkText>
          </p>
        </Button>
      </div>
    </div>
  );
}
