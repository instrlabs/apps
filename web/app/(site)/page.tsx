import type { Metadata } from "next";
import ListProduct from "@/app/(site)/ListProduct";

export const metadata: Metadata = {
  title: "Home - Labs",
  description: "",
};

export default function HomeContent() {
  return (
    <div className="w-full h-full px-4 pb-4">
      <div className="w-full h-full flex flex-col">
        <ListProduct />
      </div>
    </div>
  )
}
