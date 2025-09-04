import React, {Suspense} from "react";
import NonAuthGuard from "@/components/non-auth-guard";

export default function LoginLayout({
  children,
}: Readonly<{ children: React.ReactNode }>) {
  return (
    <>
      <NonAuthGuard />
      <Suspense>
        {children}
      </Suspense>
    </>
  );
}
