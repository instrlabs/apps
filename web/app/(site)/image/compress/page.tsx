import type { Metadata } from "next";
import ImageCompress from "@/app/(site)/image/compress/ImageCompress";
import { getImageInstructions } from "@/services/images";
import { redirect } from "next/navigation";

export const metadata: Metadata = {
  title: "Image Compress - Labs",
  description: "",
};

export default async function ImageCompressPage() {
  const { success, data } = await getImageInstructions();

  if (!success) redirect("/")

  return (
    <div className="grid grid-cols-2 gap-4">
      <p>{JSON.stringify(data!.instructions)}</p>
      <ImageCompress />
    </div>
  )
}
