"use server"

import React, { Suspense } from "react";
import { cookies } from "next/headers";

import { getProfile } from "@/services/auth";
import { getProducts } from "@/services/images";

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
  const cookieStore = await cookies();

  if (!cookieStore.has('AccessToken')) {
    return (
      <div className="flex items-center justify-center h-screen">
        <div className="text-center">
          <p className="text-lg font-semibold">Please login first</p>
        </div>
      </div>
    );
  }

  const { data: profileData } = await getProfile();
  const { data: productData } = await getProducts();

  return (
    <ProfileProvider data={profileData?.user || null}>
    <ProductProvider data={productData?.products || null}>
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
