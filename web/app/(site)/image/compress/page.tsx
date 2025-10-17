import type { Metadata } from "next";
import ImageCompress from "@/app/(site)/image/compress/ImageCompress";
import History from "@/app/(site)/image/compress/History";

export const metadata: Metadata = {
  title: "Image Compress - Instruction Labs",
  description: "",
};

export default function ImageCompressPage() {
  return (
    <div className="h-full rounded-lg bg-secondary border border-primary">
      {/*<History />*/}
      <ImageCompress />
    </div>
  );
}
