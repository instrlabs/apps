import type { Metadata } from "next";
import ListProduct from "@/.(site)-backup/ListProduct";

export const metadata: Metadata = {
  title: "Home - Instruction Labs",
  description: "",
};

export default function HomeContent() {
  return (
    <>
      <ListProduct />
    </>
  )
}
