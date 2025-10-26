"use server";

import React, { Suspense } from "react";

import { getProfile } from "@/services/auth";
import { getProducts } from "@/services/images";

import { ProductProvider } from "@/hooks/useProduct";
import { ProfileProvider } from "@/hooks/useProfile";
import { OverlayProvider } from "@/hooks/useOverlay";
import { SnackbarProvider, SnackbarWidget } from "@/hooks/useSnackbar";
import { ModalProvider } from "@/hooks/useModal";
import { SSEProvider } from "@/hooks/useSSE";
import OverlayTop from "@/components/layouts/overlay-top";
import OverlayContent from "@/components/layouts/overlay-content";

export default async function SiteLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const { data: profileData } = await getProfile();
  const { data: productData } = await getProducts();

  return (
    <ProfileProvider data={profileData?.user || null}>
      <ProductProvider data={productData?.products || null}>
        <SSEProvider>
          <SnackbarProvider>
            <ModalProvider>
              <OverlayProvider>
                <Suspense>
                  <OverlayTop />
                  <OverlayContent>{children}</OverlayContent>
                  <SnackbarWidget />
                </Suspense>
              </OverlayProvider>
            </ModalProvider>
          </SnackbarProvider>
        </SSEProvider>
      </ProductProvider>
    </ProfileProvider>
  );
}
