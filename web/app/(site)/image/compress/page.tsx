import type { Metadata } from "next";
import ImageCompress from "@/app/(site)/image/compress/ImageCompress";
import History from "@/app/(site)/image/compress/History";

export const metadata: Metadata = {
  title: "Image Compress - Instruction Labs",
  description: "",
};

export default function ImageCompressPage() {
  return (
    <div className="flex flex-row gap-4 p-4">
      <History />
      <ImageCompress />
    </div>
  );
}
