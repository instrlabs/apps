import React, {Suspense} from "react";

import {getProfile} from "@/services/auth";
import {ProfileProvider} from "@/hooks/useProfile";
import {ProductProvider} from "@/hooks/useProduct";
import {listProducts} from "@/services/products";
import {OverlayProvider} from "@/hooks/useOverlay";
import OverlayTop from "@/components/overlay-top";
import OverlayBody from "@/components/overlay-body";
import OverlayLeft from "@/components/overlay-left";
import OverlayRight from "@/components/overlay-right";

export default async function LoginLayout({
  children,
}: Readonly<{ children: React.ReactNode }>) {
  const { data: profileData } = await getProfile();
  const { data: productData } = await listProducts();

  return (
    <ProfileProvider data={profileData}>
      <ProductProvider data={productData}>
        <OverlayProvider>
          <Suspense>
            <OverlayBody>
            {children}
            </OverlayBody>
            <OverlayLeft />
            <OverlayRight />
            <OverlayTop />
          </Suspense>
        </OverlayProvider>
      </ProductProvider>
    </ProfileProvider>
  );
}
