"use client"

import {useCallback, useEffect, useState} from "react"
import Button from "@/components/actions/button";
import { useRouter } from "next/navigation";
import useModal from "@/hooks/useModal";
import ImagePreviewOverlay from "@/components/layouts/image-preview";
import FileDropzone from "@/components/inputs/file-dropzone";
import {compressImage} from "@/services/images";
import {bytesToString} from "@/utils/bytesToString";

export default function ImageCompressPage() {
  const router = useRouter();
  const { openModal } = useModal();
  const [files, setFiles] = useState<File[]>([])
  const [previews, setPreviews] = useState<Record<string, string>>({})
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState<string | null>(null)


  useEffect(() => {
    const next: Record<string, string> = {}
    const urlsToRevoke: string[] = []
    for (const f of files) {
      const url = URL.createObjectURL(f)
      next[`${f.name}:${f.size}`] = url
      urlsToRevoke.push(url)
    }
    setPreviews(next)
    return () => {
      for (const url of urlsToRevoke) URL.revokeObjectURL(url)
    }
  }, [files])

  const removeFile = useCallback((name: string, size: number) => {
    setFiles((prev) => prev.filter((f) => !(f.name === name && f.size === size)))
  }, []);

  return (
    <div className="w-full flex flex-col py-10">
      <h2 className="text-center text-3xl font-bold mt-6">Compress PDF files</h2>
      <p className="text-center text-lg mt-3">
        Reduce file size while optimizing for maximal PDF quality.
      </p>

      <div className="w-full mt-8 flex flex-col items-center">
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
            {files.map((f) => {
              const key = f.name + ":" + f.size
              const preview = previews[key]
              return (
                <div key={key} className="w-full flex items-center justify-between border border-border border-dashed rounded-xl p-4 gap-4">
                  {preview ? (
                    <img
                      src={preview}
                      alt={f.name}
                      className="h-16 w-16 rounded object-contain flex-shrink-0 border cursor-zoom-in bg-white"
                      role="button"
                      tabIndex={0}
                      onClick={(e) => {
                        e.stopPropagation();
                        if (preview) {
                          openModal(<ImagePreviewOverlay src={preview} />);
                        }
                      }}
                      onKeyDown={(e) => {
                        if (e.key === "Enter" || e.key === " ") {
                          e.preventDefault();
                          if (preview) {
                            openModal(<ImagePreviewOverlay src={preview} />);
                          }
                        }
                      }}
                      aria-label={`Preview ${f.name}`}
                    />
                  ) : (
                    <div className="h-16 w-16 rounded bg-gray-100 text-gray-400 flex items-center justify-center flex-shrink-0 border">
                      {f.type.split("/")[1]?.toUpperCase() || "FILE"}
                    </div>
                  )}
                  <div className="min-w-0 flex-1 space-y-1">
                    <p className="truncate font-medium">{f.name}</p>
                    <p className="text-sm text-gray-500">
                      {f.type.split("/")[1].toUpperCase()} • {bytesToString(f.size)}
                    </p>
                  </div>
                  <button
                    className="text-sm text-red-600 hover:text-red-700"
                    onClick={(e) => {
                      e.stopPropagation()
                      removeFile(f.name, f.size)
                    }}
                    aria-label={`Remove ${f.name}`}
                  >
                    Remove
                  </button>
                </div>
              )
            })}
          </div>
        )}

        {error && (
          <div className="text-sm text-red-600 mb-3">{error}</div>
        )}
        {/*<Button*/}
        {/*  disabled={files.length === 0 || submitting}*/}
        {/*  onClick={async () => {*/}
        {/*    if (files.length === 0 || submitting) return;*/}
        {/*    setError(null);*/}
        {/*    setSubmitting(true);*/}
        {/*    try {*/}
        {/*      const res = await compressImage(files);*/}
        {/*      if (res.success) {*/}
        {/*        router.push("/histories");*/}
        {/*      } else {*/}
        {/*        setError(res.message || "Failed to start compression.");*/}
        {/*      }*/}
        {/*    } catch (e) {*/}
        {/*      setError("Unexpected error. Please try again.");*/}
        {/*    } finally {*/}
        {/*      setSubmitting(false);*/}
        {/*    }*/}
        {/*  }}*/}
        {/*  >*/}
        {/*  {submitting ? "Starting..." : "PROCEED"}*/}
        {/*</Button>*/}
      </div>
    </div>
  )
}
