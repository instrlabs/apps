import type { Metadata } from "next";
import ImageCompress from "@/.(site)-backup/image/compress/ImageCompress";
import History from "@/.(site)-backup/image/compress/History";

export const metadata: Metadata = {
  title: "Image Compress - Instruction Labs",
  description: "",
};

export default function ImageCompressPage() {
  return (
    <div className="h-full rounded-lg bg-secondary border border-primary">
      <ImageCompress />
    </div>
  );
}
