"use client";

import React, { useCallback, useEffect, useState } from "react";
import useSnackbar from "@/hooks/useSnackbar";
import useSSE from "@/hooks/useSSE";
import DefaultState from "./DefaultState";
import UploadedState from "./UploadedState";
import { CompressState, ExtendedInstructionFile } from "./types";
import {
  createInstruction,
  createInstructionDetail,
  getInstructionDetail,
  getInstructionDetailsFile,
} from "@/services/images";
import { downloadFromArrayBuffer } from "@/utils/download";

export default function ImageCompress() {
  const { showSnackbar } = useSnackbar();
  const { message } = useSSE();
  const [state, setState] = useState<CompressState>("default");
  const [files, setFiles] = useState<ExtendedInstructionFile[]>([]);
  const [instructionId, setInstructionId] = useState<string | null>(null);

  const handleFilesAdded = useCallback(
    (newFiles: File[]) => {
      const newFileItems: ExtendedInstructionFile[] = newFiles.map((file) => ({
        file_name: file.name,
        file_size: file.size,
        mime_type: file.type,
        file,
      }));

      setFiles((prevFiles) => [...prevFiles, ...newFileItems]);
      setState("uploaded");
      showSnackbar({
        type: "info",
        message: `${newFiles.length} file(s) added successfully!`,
      });
    },
    [showSnackbar],
  );

  const handleRemoveFile = useCallback((fileIndex: number) => {
    setFiles((prevFiles) => prevFiles.filter((_, idx) => idx !== fileIndex));
  }, []);

  const handleSubmit = useCallback(async () => {
    if (files.length === 0) {
      showSnackbar({
        type: "error",
        message: "Please select files to compress",
      });
      return;
    }

    try {
      const resCreateInstruction = await createInstruction("68dcb3c6ea3593376916a6a4");

      if (!resCreateInstruction.success) {
        showSnackbar({
          type: "error",
          message: resCreateInstruction.message,
        });
        return;
      }

      const currentInstructionId = resCreateInstruction.data!.instruction.id;
      setInstructionId(currentInstructionId);

      for (let i = 0; i < files.length; i++) {
        const resCreateInstructionDetail = await createInstructionDetail(currentInstructionId, files[i].file);

        if (!resCreateInstructionDetail.success) {
          setFiles((prev) =>
            prev.map((item, idx) => {
              if (idx !== i) return item;
              return { ...item, status: "FAILED" };
            }),
          );

          showSnackbar({
            type: "error",
            message: `Failed to upload ${files[i].file.name}: ${resCreateInstructionDetail.message}`,
          });

          continue;
        }

        const instructionDetailData = resCreateInstructionDetail.data!;

        setFiles((prev) =>
          prev.map((item, idx) => {
            if (idx !== i) return item;

            // Calculate compression stats from input and output
            const inputSize = instructionDetailData.input.file_size;
            const outputSize = instructionDetailData.output.file_size;
            const compressionPercentage = outputSize > 0 ? Math.round(((inputSize - outputSize) / inputSize) * 100) : 0;

            return {
              ...item,
              id: instructionDetailData.input.id,
              instruction_id: instructionDetailData.input.instruction_id,
              file_name: instructionDetailData.input.file_name,
              file_size: instructionDetailData.input.file_size,
              mime_type: instructionDetailData.input.mime_type,
              status: instructionDetailData.input.status,
              output_id: instructionDetailData.input.output_id,
              outputFile: instructionDetailData.output,
              compressedSize: outputSize,
              compressionPercentage,
            };
          }),
        );
      }

      showSnackbar({
        type: "success",
        message: "All files compressed successfully!",
      });
    } catch (error) {
      showSnackbar({
        type: "error",
        message: "Failed to process files. Please try again.",
      });

      console.error("Error processing files:", error);
    }
  }, [files, showSnackbar]);

  const handleReset = useCallback(() => {
    setFiles([]);
    setInstructionId(null);
    setState("default");
  }, []);

  const handleDownload = useCallback(
    async (instructionId: string, detailId: string, fileName: string) => {
      const buff = await getInstructionDetailsFile(instructionId, detailId);
      downloadFromArrayBuffer(buff.data!, fileName);

      showSnackbar({
        type: "success",
        message: `Downloaded ${fileName}!`,
      });
    },
    [showSnackbar],
  );

  useEffect(() => {
    if (files.length === 0) {
      setState("default");
    }
  }, [files]);

  useEffect(() => {
    if (!message || !instructionId) {
      return;
    }

    const data = message.data as {
      user_id?: string;
      instruction_id?: string;
      instruction_detail_id?: string;
    };

    if (message.eventName === "message" && data.instruction_id === instructionId && data.instruction_detail_id) {
      getInstructionDetail(instructionId, data.instruction_detail_id)
        .then((response) => {
          if (!response.success || !response.data) return;
          const detail = response.data.detail;

          const isInput = detail.output_id !== undefined;

          if (isInput) {
            setFiles((prev) =>
              prev.map((f) =>
                f.id === data.instruction_detail_id
                  ? {
                      ...f,
                      status: detail.status,
                    }
                  : f,
              ),
            );

            if (detail.output_id) {
              getInstructionDetail(instructionId, detail.output_id)
                .then((outputResponse) => {
                  if (!outputResponse.success || !outputResponse.data) return;
                  const outputDetail = outputResponse.data.detail;

                  // Update the input file with output file data
                  setFiles((prev) =>
                    prev.map((f) =>
                      f.id === data.instruction_detail_id
                        ? {
                            ...f,
                            outputFile: outputDetail,
                          }
                        : f,
                    ),
                  );

                  if (outputDetail.status === "DONE") {
                    const inputSize = detail.file_size;
                    const outputSize = outputDetail.file_size;
                    const compressedSize = outputSize;
                    const compressionPercentage =
                      outputSize > 0 ? Math.round(((inputSize - outputSize) / inputSize) * 100) : 0;

                    setFiles((prev) =>
                      prev.map((f) =>
                        f.id === data.instruction_detail_id
                          ? {
                              ...f,
                              compressedSize,
                              compressionPercentage,
                            }
                          : f,
                      ),
                    );
                  }
                })
                .catch((error) => {
                  console.error("Error fetching output file:", error);
                });
            }
          } else {
            setFiles((prev) =>
              prev.map((f) => {
                if (f.output_id === data.instruction_detail_id) {
                  const updatedFile = {
                    ...f,
                    outputFile: detail,
                  };

                  if (detail.status === "DONE" && f.file_size) {
                    const inputSize = f.file_size;
                    const outputSize = detail.file_size;
                    const compressedSize = outputSize;
                    const compressionPercentage =
                      outputSize > 0 ? Math.round(((inputSize - outputSize) / inputSize) * 100) : 0;

                    return {
                      ...updatedFile,
                      compressedSize,
                      compressionPercentage,
                    };
                  }

                  return updatedFile;
                }
                return f;
              }),
            );
          }
        })
        .catch((error) => {
          console.error("Error fetching instruction detail:", error);
        });
    }
  }, [message, instructionId]);

  return (
    <div className="flex h-full w-full flex-col gap-2 sm:gap-3">
      <div className="flex h-full flex-col gap-2 sm:gap-4 rounded-lg border border-white/10 bg-white/8 p-2 sm:p-3 md:p-4">
        <h1 className="text-base leading-6 font-semibold text-white">Image Compress</h1>

        {state === "default" && <DefaultState onFilesAdded={handleFilesAdded} />}

        {state === "uploaded" && (
          <UploadedState
            files={files}
            onRemoveFile={handleRemoveFile}
            onSubmit={handleSubmit}
            onReset={handleReset}
            onDownload={handleDownload}
          />
        )}
      </div>
    </div>
  );
}
