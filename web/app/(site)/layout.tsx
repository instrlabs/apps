"use server"

import React, {Suspense} from "react";
import OverlayBody from "@/components/layouts/overlay-body";
import { getProfile } from "@/services/auth";
import {ProfileProvider} from "@/hooks/useProfile";
import {listProducts} from "@/services/products";
import {ProductProvider} from "@/hooks/useProduct";
import OverlayLeft from "@/components/layouts/overlay-left";
import OverlayRight from "@/components/layouts/overlay-right";
import OverlayTop from "@/components/layouts/overlay-top";

import { NotificationProvider, NotificationWidget } from "@/hooks/useNotification";
import { OverlayProvider } from "@/hooks/useOverlay";
import { ModalProvider } from "@/hooks/useModal";
import {SSEProvider} from "@/hooks/useSSE";

export default async function SiteLayout({ children }: Readonly<{
  children: React.ReactNode;
}>) {
  const { data: profileData } = await getProfile();
  const { data: productData } = await listProducts();

  return (
    <ProfileProvider data={profileData}>
    <ProductProvider data={productData}>
    <SSEProvider>
    <NotificationProvider>
    <ModalProvider>
    <OverlayProvider defaultLeft="left:navigation">
      <Suspense>
        <OverlayBody>{children}</OverlayBody>
        <OverlayLeft />
        <OverlayRight />
        <OverlayTop />
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
