import React, {Suspense} from "react";
import Providers from "@/app/providers";
import Widgets from "@/app/widgets";
import OverlayContentWrapper from "@/components/overlay-content-wrapper";

export default async function SiteLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const profile = await fetch("http://localhost:8000/api/auth/profile");

  return (
    <Providers>
      <OverlayContentWrapper>
        <Suspense>{children}</Suspense>
      </OverlayContentWrapper>
      <Widgets />
    </Providers>
  );
}
