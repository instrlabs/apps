import type { Metadata } from "next";
import ImageCompress from "@/app/(site)/image/compress/ImageCompress";

export const metadata: Metadata = {
  title: "Image Compress - Labs",
  description: "",
};

export default function ImageCompressPage() {
  return (
    <>
      <ImageCompress />
    </>
  )
}
