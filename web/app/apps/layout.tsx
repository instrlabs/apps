"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import clsx from "clsx";

import SidebarLeft from "@/components/sidebar-left";
import AppBar from "@/components/appbar";

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const router = useRouter();
  //
  // useEffect(() => {
  //   // Check if we're in a browser environment
  //   if (typeof window !== "undefined") {
  //     // Check if user is authenticated
  //     const authToken = localStorage.getItem("authToken");
  //     if (!authToken) {
  //       // Redirect to login page if not authenticated
  //       router.push("/login");
  //     }
  //   }
  // }, [router]);

  return (
    <>
      <SidebarLeft />
      <div className={clsx("absolute inset-0 w-full h-full overflow-hidden")}>
        <AppBar />
        {children}
      </div>
    </>
  );
}
