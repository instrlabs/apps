"use server"

import React, { Suspense } from "react";

import { getProfile } from "@/services/auth";
import { getProducts } from "@/services/products";

import { ProductProvider } from "@/hooks/useProduct";
import { ProfileProvider } from "@/hooks/useProfile";
import { OverlayProvider } from "@/hooks/useOverlay";
import { NotificationProvider, NotificationWidget } from "@/hooks/useNotification";
import { ModalProvider } from "@/hooks/useModal";
import { SSEProvider } from "@/hooks/useSSE";
import OverlayBody from "@/components/layouts/overlay-body";
import OverlayTop from "@/components/layouts/overlay-top";
import OverlayRight from "@/components/layouts/overlay-right";

export default async function SiteLayout({ children }: Readonly<{
  children: React.ReactNode;
}>) {
  const { data: profileData } = await getProfile();
  const { data: productData } = await getProducts();

  return (
    <ProfileProvider data={profileData}>
    <ProductProvider data={productData}>
    <SSEProvider>
    <NotificationProvider>
    <ModalProvider>
    <OverlayProvider>
      <Suspense>
        <OverlayBody>{children}</OverlayBody>
        <OverlayTop />
        <OverlayRight />
        <NotificationWidget />
      </Suspense>
    </OverlayProvider>
    </ModalProvider>
    </NotificationProvider>
    </SSEProvider>
    </ProductProvider>
    </ProfileProvider>
  );
}
