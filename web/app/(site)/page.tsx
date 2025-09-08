import { Suspense } from "react";

import ListProduct from "@/app/(site)/ListProduct";

export default function HomeContent() {
  return (
    <>
      <Suspense>
        <ListProduct />
      </Suspense>
    </>
  )
}
