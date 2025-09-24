"use client"

import {useCallback, useEffect, useState} from "react"
import FileDropzone from "@/components/inputs/file-dropzone";
import Button from "@/components/actions/button";
import {
  createInstruction,
  Instruction,
  InstructionFile,
  createInstructionDetails, getInstructionDetails
} from "@/services/images";
import useNotification from "@/hooks/useNotification";
import ListFiles from "@/app/(site)/images/compress/ListFiles";
import Loading from "@/components/feedback/Loading";
import SubmittedClient from "@/app/(site)/images/compress/SubmittedClient";
import useSSE from "@/hooks/useSSE";
import debounce from "lodash.debounce";
import ButtonIcon from "@/components/actions/button-icon";
import TrashIcon from "@/components/icons/TrashIcon";
import CloseIcon from "@/components/icons/CloseIcon";
import CloudUploadIcon from "@/components/icons/CloudUploadIcon";

type ProgressStatus = 'UPLOAD' | 'UPLOADED' | 'SUBMITTED' | 'LOADING';

const FILE_SIZE_5MB = 5242880;

export default function ImageCompressPage() {
  const [progress, setProgress] = useState<ProgressStatus>('UPLOAD');
  const [files, setFiles] = useState<File[]>([]);
  const [urls, setUrls] = useState<string[]>([]);
  const [instruction, setInstruction] = useState<Instruction | null>(null);
  const [inputFiles, setInputFiles] = useState<InstructionFile[]>([]);
  const [outputFiles, setOutputFiles] = useState<InstructionFile[]>([]);
  const { showNotification } = useNotification();
  const { message } = useSSE();

  useEffect(() => {
    if (!message) return;
    const {eventName, data} = message;
    if (eventName === "message") {
      const payload = data as { instruction_id: string; file_id: string; };
      const debouncedProcessFiles = debounce(async () => {
        const result = await getInstructionDetails(payload.instruction_id);
        if (result.success && result.data) {
          const inputs = result.data.files.filter(f => f.output_id);
          const outputs = result.data.files.filter(f => !f.output_id);
          setInputFiles(inputs);
          setOutputFiles(outputs);
        }
      }, 500);

      debouncedProcessFiles();

      return () => {
        debouncedProcessFiles.cancel();
      };
    }
  }, [message]);

  const removeFile = useCallback((idx: number) => {
    setFiles((prev) => {
      const next = [...prev.slice(0, idx), ...prev.slice(idx + 1)];
      if (next.length === 0) setProgress('UPLOAD');
      return next;
    })
  }, []);

  const handleUpload = (files: File[]) => {
    setProgress('LOADING');
    setFiles(files);
    setUrls(files.map(f => URL.createObjectURL(f)));
    setProgress('UPLOADED');
  }

  const handleSubmit = async () => {
    setProgress('LOADING');
    const compressResult = await createInstruction("images-compress");
    if (!compressResult.success || !compressResult.data) {
      showNotification({ message: compressResult.message, type: "error", duration: 3000});
      return;
    } else {
      const instr = compressResult.data.instruction;
      setInstruction(instr);

      for (const file of files) {
        const fileResult = await createInstructionDetails(instr.id, file);
        if (!fileResult.success || !fileResult.data) {
          showNotification({ message: fileResult.message, type: "error", duration: 3000});
          break;
        } else {
          const { input, output } = fileResult.data;
          setInputFiles(prev => [...prev, input as InstructionFile]);
          setOutputFiles(prev => [...prev, output as InstructionFile]);
        }
      }
    }

    setProgress('SUBMITTED');
  }

  return (
    <div className="flex h-full w-full flex-row gap-4 px-4 pb-4">
      <div className="flex w-[340px] flex-col gap-4">
        <div className="flex flex-col gap-2">
          <h4 className="text-sm">Histories</h4>
        </div>
        <div className="bg-primary-black shadow-primary flex w-full flex-col gap-4 rounded-lg p-4">
          <span className="text-sm">Last 30 days</span>
          <div className="flex flex-col items-center gap-2">
            <div className="flex w-full flex-row items-center gap-2">
              <CloseIcon className="size-4 shrink-0 cursor-pointer text-white/70" />
              <span className="shrink truncate text-sm font-light text-white/50">
                WhatsApp Image 2025-07-31 at 18.39.20.jpeg
              </span>
              <span className="shrink-0 text-sm font-light text-white/50">1.2 MB</span>
            </div>
          </div>
        </div>
      </div>

      <div className="flex h-full flex-1 flex-col gap-4">
        <div className="flex flex-col gap-2">
          <h4 className="text-sm">Image Compress</h4>
        </div>

        <div className="flex h-full w-full flex-col">
          {progress === "LOADING" && (
            <div className={`
              flex h-full w-full flex-col items-center justify-center
              bg-primary-black shadow-primary rounded-lg
            `}>
              <Loading />
            </div>
          )}

          {progress === "UPLOAD" && (
            <FileDropzone
              multiple
              accepts={["image/png", "image/jpeg", "image/webp", "image/gif"]}
              maxSize={FILE_SIZE_5MB}
              onFilesAdded={handleUpload}
            />
          )}

          {progress === "UPLOADED" && (
            <>
              <ListFiles files={files} imagesUrls={urls} removeFile={removeFile} />
              <Button xSize="lg" onClick={handleSubmit}>
                Submit
              </Button>
            </>
          )}

          {progress === "SUBMITTED" && instruction && inputFiles.length > 0 && (
            <SubmittedClient
              instructionId={instruction.id}
              inputFiles={inputFiles}
              outputFiles={outputFiles}
            />
          )}
        </div>
      </div>
    </div>
  );
}
