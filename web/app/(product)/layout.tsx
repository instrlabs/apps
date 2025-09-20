import React, {Suspense} from "react";

import {getProfile} from "@/services/authentications";
import {ProfileProvider} from "@/hooks/useProfile";
import {ProductProvider} from "@/hooks/useProduct";
import {listProducts} from "@/services/products";
import {OverlayProvider} from "@/hooks/useOverlay";
import OverlayTop from "@/components/layouts/overlay-top";
import OverlayBody from "@/components/layouts/overlay-body";
import OverlayLeft from "@/components/layouts/overlay-left";
import OverlayRight from "@/components/layouts/overlay-right";
import { ModalProvider } from "@/hooks/useModal";
import {NotificationProvider, NotificationWidget} from "@/hooks/useNotification";
import {SSEProvider} from "@/hooks/useSSE";

export default async function LoginLayout({
  children,
}: Readonly<{ children: React.ReactNode }>) {
  const { data: profileData } = await getProfile();
  const { data: productData } = await listProducts();

  return (
    <ProfileProvider data={profileData}>
    <ProductProvider data={productData}>
    <SSEProvider>
    <NotificationProvider>
    <ModalProvider>
    <OverlayProvider>
      <Suspense>
        <OverlayBody>
        {children}
        </OverlayBody>
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
