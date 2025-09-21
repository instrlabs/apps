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
import ListFiles from "@/app/(product)/images/compress/ListFiles";
import Loading from "@/components/feedback/Loading";
import SubmittedClient from "@/app/(product)/images/compress/SubmittedClient";
import useSSE from "@/hooks/useSSE";
import debounce from "lodash.debounce";

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
      showNotification({ title: "Error", message: compressResult.message, type: "error", duration: 3000});
      return;
    } else {
      const instr = compressResult.data.instruction;
      setInstruction(instr);

      for (const file of files) {
        const fileResult = await createInstructionDetails(instr.id, file);
        if (!fileResult.success || !fileResult.data) {
          showNotification({ title: "Error", message: fileResult.message, type: "error", duration: 3000});
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
    <div className="w-full flex flex-col items-center py-10">
      <h2 className="text-center text-3xl font-bold mt-6">
        Compress Images
      </h2>
      <p className="text-center text-lg mt-3">
        Reduce file size while optimizing for maximal image quality.
      </p>

      <div className="w-full max-w-3xl mt-8 flex flex-col items-center space-y-4">
        {progress === 'LOADING' && <Loading size={90} />}

        {progress === 'UPLOAD' && (
          <FileDropzone
            multiple
            accepts={["image/png", "image/jpeg", "image/webp", "image/gif"]}
            maxSize={FILE_SIZE_5MB}
            onFilesAdded={handleUpload}
          />
        )}

        {progress === 'UPLOADED' && (
          <>
            <ListFiles files={files} imagesUrls={urls} removeFile={removeFile} />
            <Button xSize="lg" onClick={handleSubmit}>Submit</Button>
          </>
        )}

        {progress === 'SUBMITTED' && instruction && inputFiles.length > 0 && (
          <SubmittedClient instructionId={instruction.id} inputFiles={inputFiles} outputFiles={outputFiles} />
        )}
      </div>
    </div>
  )
}
