"use client";

import React from "react";
import { useForm } from "react-hook-form";

import Button from "@/components/actions/button";
import TextField from "@/components/inputs/text-field";
import InputPin from "@/components/inputs/input-pin";
import { login, sendPin } from "@/services/auth";
import useNotification from "@/hooks/useNotification";
import InlineSpinner from "@/components/feedback/InlineSpinner";
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
    <>
      <h3 className="text-center text-3xl font-semibold">Log in to Labs</h3>
      <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4">
        <TextField
          type="email"
          placeholder="Email Address"
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
          xVariant="solid"
          disabled={loading}
        >
          <div className="flex items-center justify-center gap-2">
            {loading && <InlineSpinner />} <span>Continue with Email</span>
          </div>
        </Button>
      </form>
    </>
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
    <>
      <h3 className="text-center text-3xl font-semibold">Verification</h3>
      <p className="text-center text-white/70">
        If you have an account, we have sent a code to <b>{email}</b>. Enter it below.
      </p>
      <form onSubmit={handleSubmit} className="flex flex-col gap-4">
        <div className="mx-auto">
          <InputPin values={values} onChange={setValues} />
        </div>
        <Button
          type="submit"
          xVariant="solid"
          disabled={loading}
        >
          <div className="flex items-center justify-center gap-2">
            {loading && <InlineSpinner />} <span>Continue</span>
          </div>
        </Button>
      </form>
    </>
  );
}
