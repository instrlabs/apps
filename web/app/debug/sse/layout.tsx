"use server"

import React from "react";
import {SSEProvider} from "@/hooks/useSSE";

export default async function SiteLayout({ children }: Readonly<{
  children: React.ReactNode;
}>) {

  return (
    <SSEProvider>
      {children}
    </SSEProvider>
  );
}
