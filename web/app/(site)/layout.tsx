import React, {Suspense} from "react";
import Providers from "@/app/providers";
import Widgets from "@/app/widgets";
import OverlayContentWrapper from "@/components/overlay-content-wrapper";
import { profile } from "@/services/auth";
import {ProfileProvider} from "@/hooks/useProfile";
import {redirect} from "next/navigation";
import {listProducts} from "@/services/products";
import {ProductProvider} from "@/hooks/useProduct";

export default async function SiteLayout({ children }: Readonly<{
  children: React.ReactNode;
}>) {
  const { success: successProfile, data: profileData } = await profile();
  const { data: productData } = await listProducts();

  if (!successProfile || !profileData) {
    return redirect("/login");
  }

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
