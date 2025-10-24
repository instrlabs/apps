"use client";

import React from "react";
import { useForm } from "react-hook-form";

import Button from "@/components/button";
import Input from "@/components/inputs/input";
import InputPin from "@/components/inputs/input-pin";
import GoogleSignInButton from "@/components/google-signin";
import LogoIcon from "@/components/icons/logo-icon";
import { login, sendPin } from "@/services/auth";
import useNotification from "@/hooks/useNotification";
import InlineSpinner from "@/components/feedback/InlineSpinner";
import Text from "@/components/text";
import { redirect, RedirectType } from "next/navigation";

type FormEmailValues = {
  email: string;
};

export default function LoginForm() {
  const [email, setEmail] = React.useState("");
  const [state, setState] = React.useState<"email" | "pin">("email");
  return (
    <>
      {state === "email" && (
        <FormEmail setEmail={setEmail} next={() => { setState("pin") }} />
      )}
      {state === "pin" && (
        <FormPin email={email} next={() => redirect("/", RedirectType.replace) } />
      )}
    </>
  )
}

function FormEmail({ setEmail, next }: {
  setEmail: (email: string) => void,
  next: () => void,
}) {
  const [loading, setLoading] = React.useState(false);
  const { showNotification } = useNotification();
  const {
    register,
    handleSubmit,
  } = useForm<FormEmailValues>({
    defaultValues: { email: "" },
    mode: "onSubmit",
    reValidateMode: "onChange",
  });

  const onSubmit = async (values: FormEmailValues) => {
    setLoading(true);

    try {
      const { success, message } = await sendPin({ email: values.email });
      if (success) {
        setEmail(values.email)
        next()
      } else {
        showNotification({
          type: "error",
          message: message,
        })
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="mx-auto w-full max-w-[400px] flex flex-col gap-7 px-6 md:px-0">
      <LogoIcon size={160} className="mx-auto drop-shadow" />
      <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-7">
        <Input
          xSize="lg"
          type="email"
          placeholder="Email address"
          autoComplete="email"
          inputMode="email"
          {...register("email", {
            required: "Email is required",
            pattern: {
              value: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
              message: "Enter a valid email address",
            },
          })}
        />
        <Button
          type="submit"
          color="primary"
          size="lg"
          disabled={loading}
        >
          <div className="flex items-center justify-center gap-2">
            {loading && <InlineSpinner />} <span>Continue with Email</span>
          </div>
        </Button>
        <div className={`
mx-auto
bg-white/40
h-px w-4
        `} />
        <GoogleSignInButton />
      </form>
    </div>
  );
}


function FormPin({ email, next }: {
  email: string,
  next: () => void,
}) {
  const [loading, setLoading] = React.useState(false);
  const [values, setValues] = React.useState<string[]>(Array(6).fill(""));
  const { showNotification } = useNotification();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);

    try {
    const { success, message} = await login({ email, pin: values.join("") });
      if (success) {
        next()
      } else {
        showNotification({
          type: "error",
          message: message,
        })
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={`
mx-auto w-full max-w-[400px] flex flex-col
gap-7 px-6 md:px-0
   `}>
      <Text as="h3" className={`
text-center
      `} isBold xSize="xl">
        Verification
      </Text>
      <Text as="p" className={`
text-center
      `} xColor="secondary">
        If you have an account, we have sent a code to <b>{email}</b>. Enter it below.
      </Text>
      <form onSubmit={handleSubmit} className={`
flex flex-col
gap-7
      `}>
        <InputPin values={values} onChange={setValues} />
        <Button
          type="submit"
          color="primary"
          size="lg"
          disabled={loading}
        >
          <div className={`
            /* layout */
            flex items-center justify-center
            /* spacing */
            gap-2
            /* borders */

            /* colors */

            /* text */

            /* effects */

            /* states */
         `}>
            {loading && <InlineSpinner />} <span>Continue</span>
          </div>
        </Button>
        <Button
          type="button"
          color="secondary"
          size="lg"
          onClick={async () => {
            setLoading(true);
            try {
              const { success, message } = await sendPin({ email });
              if (success) {
                showNotification({ type: "info", message: "Code resent to your email" });
              } else {
                showNotification({ type: "error", message });
              }
            } finally {
              setLoading(false);
            }
          }}
        >
          Resend code
        </Button>
      </form>
    </div>
  );
}
