import React, {Suspense} from "react";
import Providers from "@/app/providers";
import Widgets from "@/app/widgets";
import OverlayContentWrapper from "@/components/overlay-content-wrapper";

export default function SiteLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <Providers>
      <OverlayContentWrapper>
        <Suspense>{children}</Suspense>
      </OverlayContentWrapper>
      <Widgets />
    </Providers>
  );
}
