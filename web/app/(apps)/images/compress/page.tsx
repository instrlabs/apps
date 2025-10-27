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
  getInstructionDetails,
} from "@/services/images";

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
        const resCreateInstructionDetail = await createInstructionDetail(
          currentInstructionId,
          files[i].file,
        );

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
            const compressionPercentage =
              outputSize > 0 ? Math.round(((inputSize - outputSize) / inputSize) * 100) : 0;

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

  useEffect(() => {
    if (files.length === 0) {
      setState("default");
    }
  }, [files]);

  useEffect(() => {
    if (!message || !instructionId) {
      console.log("SSE Effect: Missing message or instructionId");
      return;
    }

    const data = message.data as {
      user_id?: string;
      instruction_id?: string;
      instruction_detail_id?: string
    };

    if (message.eventName === "message" && data.instruction_id === instructionId) {
      console.log("SSE Effect: Received notification for detail:", data.instruction_detail_id);

      // Update specific file in state when we receive a detail-specific notification
      if (data.instruction_detail_id) {
        setFiles((prevFiles) =>
          prevFiles.map((prevFile) => {
            if (prevFile.id === data.instruction_detail_id) {
              // This file was updated, fetch its latest status
              getInstructionDetails(instructionId)
                .then((response) => {
                  if (response.success && response.data) {
                    const backendFiles = response.data.files;
                    const inputFiles = backendFiles.filter((f) => !f.output_id || f.id !== f.output_id);
                    const outputFiles = backendFiles.filter((f) => f.output_id);

                    const backendInputFile = inputFiles.find((bf) => bf.id === prevFile.id);
                    const backendOutputFile = outputFiles.find(
                      (of) => of.id === backendInputFile?.output_id,
                    );

                    if (backendInputFile) {
                      let compressedSize = prevFile.compressedSize;
                      let compressionPercentage = prevFile.compressionPercentage;

                      if (backendOutputFile && backendOutputFile.status === "DONE") {
                        const inputSize = backendInputFile.file_size;
                        const outputSize = backendOutputFile.file_size;
                        compressedSize = outputSize;
                        compressionPercentage =
                          outputSize > 0 ? Math.round(((inputSize - outputSize) / inputSize) * 100) : 0;
                      }

                      // Update the specific file with fresh data
                      setFiles((currentFiles) =>
                        currentFiles.map((file) =>
                          file.id === data.instruction_detail_id
                            ? {
                                ...file,
                                ...backendInputFile,
                                outputFile: backendOutputFile,
                                compressedSize,
                                compressionPercentage,
                              }
                            : file
                        )
                      );
                    }
                  }
                })
                .catch((error) => {
                  console.error("Error fetching instruction details:", error);
                });

              // Return optimistic update while fetching
              return {
                ...prevFile,
                // Note: We could add optimistic status updates here if needed
              };
            }
            return prevFile;
          })
        );
      } else {
        // Fallback: refetch all files if no specific detail ID (legacy behavior)
        getInstructionDetails(instructionId)
          .then((response) => {
            if (response.success && response.data) {
              const backendFiles = response.data.files;
              const inputFiles = backendFiles.filter((f) => !f.output_id || f.id !== f.output_id);
              const outputFiles = backendFiles.filter((f) => f.output_id);

              setFiles((prevFiles) =>
                prevFiles.map((prevFile) => {
                  const backendInputFile = inputFiles.find((bf) => bf.id === prevFile.id);
                  if (backendInputFile) {
                    const backendOutputFile = outputFiles.find(
                      (of) => of.id === backendInputFile.output_id,
                    );

                    let compressedSize = prevFile.compressedSize;
                    let compressionPercentage = prevFile.compressionPercentage;

                    if (backendOutputFile && backendOutputFile.status === "DONE") {
                      const inputSize = backendInputFile.file_size;
                      const outputSize = backendOutputFile.file_size;
                      compressedSize = outputSize;
                      compressionPercentage =
                        outputSize > 0 ? Math.round(((inputSize - outputSize) / inputSize) * 100) : 0;
                    }

                    return {
                      ...prevFile,
                      ...backendInputFile,
                      outputFile: backendOutputFile,
                      compressedSize,
                      compressionPercentage,
                    };
                  }
                  return prevFile;
                }),
              );
            }
          })
          .catch((error) => {
            console.error("Error fetching instruction details:", error);
          });
      }
    }
  }, [message, instructionId]);

  return (
    <div className="flex h-full w-full flex-col gap-2">
      <div className="flex h-full flex-col gap-4 rounded-lg border border-white/10 bg-white/8 p-4">
        <h1 className="text-base leading-6 font-semibold text-white">Image Compress</h1>

        {state === "default" && <DefaultState onFilesAdded={handleFilesAdded} />}

        {state === "uploaded" && (
          <UploadedState
            files={files}
            onRemoveFile={handleRemoveFile}
            onSubmit={handleSubmit}
            onReset={handleReset}
            onDownload={(fileName) => {
              showSnackbar({
                type: "info",
                message: `Download started for ${fileName}!`,
              });
            }}
          />
        )}
      </div>
    </div>
  );
}
