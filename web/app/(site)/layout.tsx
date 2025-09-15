"use server"

import React, {Suspense} from "react";
import OverlayBody from "@/components/overlay-body";
import { getProfile } from "@/services/auth";
import {ProfileProvider} from "@/hooks/useProfile";
import {listProducts} from "@/services/products";
import {ProductProvider} from "@/hooks/useProduct";
import OverlayLeft from "@/components/overlay-left";
import OverlayRight from "@/components/overlay-right";
import OverlayTop from "@/components/overlay-top";

import { NotificationProvider, NotificationWidget } from "@/hooks/useNotification";
import { OverlayProvider } from "@/hooks/useOverlay";
import { ModalProvider } from "@/hooks/useModal";

export default async function SiteLayout({ children }: Readonly<{
  children: React.ReactNode;
}>) {
  const { data: profileData } = await getProfile();
  const { data: productData } = await listProducts();

  return (
    <ProfileProvider data={profileData}>
    <ProductProvider data={productData}>
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
    </ProductProvider>
    </ProfileProvider>
  );
}
