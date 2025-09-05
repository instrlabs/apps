"use client"

import { useCallback, useMemo, useRef, useState } from "react"
import Button from "@/components/button";

function formatBytes(bytes: number) {
  if (bytes === 0) return "0 B"
  const k = 1024
  const sizes = ["B", "KB", "MB", "GB", "TB"]
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i]
}

export default function ImageCompressPage() {
  const [files, setFiles] = useState<File[]>([])
  const [isDragging, setIsDragging] = useState(false)
  const inputRef = useRef<HTMLInputElement>(null)

  const accept = useMemo(
    () => [
      "image/png",
      "image/jpeg",
      "image/webp",
      "image/gif",
      "application/pdf",
    ].join(","),
    []
  )

  const onFilesAdded = useCallback((newFiles: FileList | File[]) => {
    const list = Array.from(newFiles)
    // Filter by accepted types
    const accepted = list.filter((f) => accept.split(",").includes(f.type))
    setFiles((prev) => {
      // Avoid duplicates by name+size
      const map = new Map(prev.map((f) => [f.name + ":" + f.size, f]))
      for (const f of accepted) {
        const key = f.name + ":" + f.size
        if (!map.has(key)) map.set(key, f)
      }
      return Array.from(map.values())
    })
  }, [accept])

  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    e.stopPropagation()
    if (!isDragging) setIsDragging(true)
  }, [isDragging])

  const handleDragLeave = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    e.stopPropagation()
    setIsDragging(false)
  }, [])

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    e.stopPropagation()
    setIsDragging(false)
    if (e.dataTransfer?.files && e.dataTransfer.files.length > 0) {
      onFilesAdded(e.dataTransfer.files)
      e.dataTransfer.clearData()
    }
  }, [onFilesAdded])

  const openFileDialog = useCallback(() => {
    inputRef.current?.click()
  }, [])

  const removeFile = useCallback((name: string, size: number) => {
    setFiles((prev) => prev.filter((f) => !(f.name === name && f.size === size)))
  }, [])

  return (
    <div className="w-full flex flex-col py-10">
      <h2 className="text-center text-3xl font-bold mt-6">Compress PDF files</h2>
      <p className="text-center text-lg mt-3">
        Reduce file size while optimizing for maximal PDF quality.
      </p>

      <div className="w-full mt-8 flex flex-col items-center">
        {files.length === 0 && (
          <div
            role="button"
            tabIndex={0}
            aria-label="Upload files by dragging and dropping or by browsing"
            onClick={openFileDialog}
            onKeyDown={(e) => {
              if (e.key === "Enter" || e.key === " ") {
                openFileDialog();
              }
            }}
            onDragOver={handleDragOver}
            onDragEnter={handleDragOver}
            onDragLeave={handleDragLeave}
            onDrop={handleDrop}
            className={
              `w-full max-w-2xl aspect-video flex flex-col items-center justify-center gap-3 ` +
              `border-1 border-dashed rounded-xl p-10 cursor-pointer `
            }
          >
            <div className="text-center">
              <p className="text-base font-light">Maximum file size: 50mb</p>
              <p className="text-base font-light">Supports .PNG, .JPG, .WEBP, .GIF</p>
            </div>
            <input
              ref={inputRef}
              type="file"
              accept={accept}
              multiple
              className="hidden"
              onChange={(e) => {
                if (e.target.files) onFilesAdded(e.target.files)
                e.currentTarget.value = ""
              }}
            />
          </div>
        )}

        {files.length > 0 && (
          <div className="w-full max-w-2xl space-y-4">
            {files.map((f) => (
              <div key={f.name + f.size} className="w-full flex items-center justify-between border border-dashed rounded-xl p-4">
                <div className="min-w-0 space-y-1">
                  <p className="truncate font-medium">{f.name}</p>
                  <p className="text-sm text-gray-500">
                    {f.type.split("/")[1].toUpperCase()} • {formatBytes(f.size)}
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
            ))}
          </div>
        )}

        <Button>PROCEED</Button>
      </div>
    </div>
  )
}
