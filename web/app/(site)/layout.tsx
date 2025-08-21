import React from "react";
import OverlayContent from "@/components/overlay-content";
import Providers from "@/app/providers";
import Widgets from "@/app/widgets";

export default function SiteLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <Providers>
      <OverlayContent>
        {children}
      </OverlayContent>
      <Widgets />
    </Providers>
  );
}
