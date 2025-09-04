"use server"

import React, {Suspense} from "react";
import Providers from "@/app/providers";
import Widgets from "@/app/widgets";
import OverlayContentWrapper from "@/components/overlay-content-wrapper";
import { getProfile } from "@/services/auth";
import {ProfileProvider} from "@/hooks/useProfile";
import {listProducts} from "@/services/products";
import {ProductProvider} from "@/hooks/useProduct";

export default async function SiteLayout({ children }: Readonly<{
  children: React.ReactNode;
}>) {
  const { data: profileData } = await getProfile();
  const { data: productData } = await listProducts();

  return (
    <ProfileProvider data={profileData}>
      <ProductProvider data={productData} >
        <Providers>
          <OverlayContentWrapper>
            <Suspense>{children}</Suspense>
          </OverlayContentWrapper>
          <Widgets />
        </Providers>
      </ProductProvider>
    </ProfileProvider>
  );
}
