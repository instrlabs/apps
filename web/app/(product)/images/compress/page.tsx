"use client"

import {useCallback, useState} from "react"
import { useRouter } from "next/navigation";
import FileDropzone from "@/components/inputs/file-dropzone";
import {bytesToString} from "@/utils/bytesToString";
import CloseIcon from "@/components/icons/CloseIcon";
import ImagePreview from "@/components/ImagePreview";
import useObjectUrl from "@/hooks/useObjectUrl";
import Button from "@/components/actions/button";
import {compressImage} from "@/services/images";
import useNotification from "@/hooks/useNotification";

export default function ImageCompressPage() {
  const router = useRouter();
  const [files, setFiles] = useState<File[]>([])
  const urls = useObjectUrl(files)
  const [submitting, setSubmitting] = useState(false)
  const { showNotification } = useNotification();


  const removeFile = useCallback((name: string, size: number) => {
    setFiles((prev) => prev.filter((f) => !(f.name === name && f.size === size)))
  }, []);

  return (
    <div className="w-full flex flex-col py-10">
      <h2 className="text-center text-3xl font-bold mt-6">Compress PDF files</h2>
      <p className="text-center text-lg mt-3">
        Reduce file size while optimizing for maximal PDF quality.
      </p>

      <div className="w-full mt-8 flex flex-col items-center space-y-4">
        {files.length === 0 && (
          <FileDropzone
            multiple
            accepts={["image/png", "image/jpeg", "image/webp", "image/gif"]}
            onFilesAdded={setFiles}
            maxFileSize={5242880}
          />
        )}

        {files.length > 0 && (
          <div className="w-full max-w-2xl space-y-4">
            {files.map((f, idx) => {
              const preview = urls[idx]
              if (!preview) return null
              return (
                <div
                  key={idx}
                  className="w-full flex items-center justify-between card p-3 gap-4"
                >
                  <ImagePreview
                    src={preview}
                    alt={f.name}
                    width={60}
                    height={60}
                  />
                  <div className="min-w-0 flex-1">
                    <p className="truncate font-medium">{f.name}</p>
                    <p className="truncate font-light">{bytesToString(f.size)}</p>
                  </div>
                  <button
                    className="text-sm text-red-400 hover:text-red-700 cursor-pointer"
                    onClick={(e) => {
                      e.stopPropagation()
                      removeFile(f.name, f.size)
                    }}
                  >
                    <CloseIcon />
                  </button>
                </div>
              )
            })}
          </div>
        )}

        {files.length > 0 && (
          <Button
            xSize="lg"
            onClick={async () => {
              setSubmitting(true);
              const res = await compressImage(files);
              if (res.success && res.data) {
                router.push("/images/compress/" + res.data.instruction_id);
              } else {
                showNotification({ title: "Error", message: res.message, type: "error", duration: 3000});
              }
            }}
          >
            {submitting ? "Starting..." : "Submit"}
          </Button>
        )}
      </div>
    </div>
  )
}
