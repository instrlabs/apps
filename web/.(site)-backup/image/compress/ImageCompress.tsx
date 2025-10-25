"use client"

import {useCallback, useEffect, useMemo, useState} from "react"
import FileDropzone from "@/components/file-dropzone";
import Text from "@/components/text";
import {
  createImageInstruction,
  Instruction,
  InstructionFile,
  createImageInstructionDetails, getImageInstructionDetails
} from "@/services/images";
import useNotification from "@/hooks/useNotification";
import ListFiles from "@/.(site)-backup/image/compress/ListFiles";
import Loading from "@/components/feedback/Loading";
import useSSE from "@/hooks/useSSE";
import debounce from "lodash.debounce";
import { bytesToString } from "@/utils/bytesToString";
import { useRouter } from "next/navigation";

type ProgressStatus = 'UPLOAD' | 'UPLOADED' | 'LOADING';

const FILE_SIZE_5MB = 5242880;

export default function ImageCompress() {
  const router = useRouter();
  const [progress, setProgress] = useState<ProgressStatus>('UPLOAD');
  const [files, setFiles] = useState<File[]>([]);
  const [urls, setUrls] = useState<string[]>([]);
  const [instruction, setInstruction] = useState<Instruction | null>(null);
  const [inputFiles, setInputFiles] = useState<InstructionFile[]>([]);
  const [outputFiles, setOutputFiles] = useState<InstructionFile[]>([]);
  const [hasFailure, setHasFailure] = useState<boolean>(false);
  const { showNotification } = useNotification();
  const { message } = useSSE();

  const anyFailed = useMemo(() => {
    if (!instruction) return false;
    for (let i = 0; i < files.length; i++) {
      const inStatus = inputFiles[i]?.status;
      const outStatus = outputFiles[i]?.status;
      if (inStatus === 'FAILED' || outStatus === 'FAILED') return true;
    }
    return false;
  }, [instruction, files.length, inputFiles, outputFiles]);

  const allDone = useMemo(() => {
    if (!instruction) return false;
    if (files.length === 0) return false;
    for (let i = 0; i < files.length; i++) {
      if (outputFiles[i]?.status !== 'DONE') return false;
    }
    return true;
  }, [instruction, files.length, outputFiles]);

  useEffect(() => {
    if (!message) return;
    const {eventName, data} = message;

    if (eventName === "message") {
      const payload = data as { instruction_id: string; file_id: string; };
      const debouncedProcessFiles = debounce(async () => {
        const result = await getImageInstructionDetails(payload.instruction_id);
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
    setHasFailure(false);
    const compressResult = await createImageInstruction("image-compress");
    if (!compressResult.success || !compressResult.data) {
      setHasFailure(true);
      showNotification({ message: compressResult.message, type: "error", duration: 3000});
      return;
    } else {
      const instr = compressResult.data.instruction;
      setInstruction(instr);

      for (const file of files) {
        const fileResult = await createImageInstructionDetails(instr.id, file);
        if (!fileResult.success || !fileResult.data) {
          setHasFailure(true);
          showNotification({ message: fileResult.message, type: "error", duration: 3000});
          break;
        } else {
          const { input, output } = fileResult.data;
          setInputFiles(prev => [...prev, input as InstructionFile]);
          setOutputFiles(prev => [...prev, output as InstructionFile]);
        }
      }
    }
  }

  const handleReset = useCallback(() => {
    try {
      urls.forEach(u => {
        try { URL.revokeObjectURL(u); } catch (_) {}
      });
    } finally {
      setFiles([]);
      setUrls([]);
      setInstruction(null);
      setInputFiles([]);
      setOutputFiles([]);
      setHasFailure(false);
      setProgress('UPLOAD');
    }
  }, [urls]);

  const goHome = useCallback(() => {
    router.push('/');
  }, [router]);

  const showInitialActions = progress === 'UPLOADED' && !instruction && !hasFailure;
  const showFinalActions = progress === 'UPLOADED' && ((!instruction && hasFailure) || (instruction && (anyFailed || allDone)));

  return (
    <div className="p-4 flex flex-col gap-4">
      <Text as="h3" xSize="sm" isBold>Image Compress</Text>

        {progress === "LOADING" && <Loading />}

        {progress === "UPLOAD" && (
          <FileDropzone
            multiple
            accepts={["image/png", "image/jpeg", "image/webp", "image/gif"]}
            maxSize={FILE_SIZE_5MB}
            onFilesAdded={handleUpload}
          />
        )}

        {/*  {progress === "UPLOADED" && (*/}
        {/*    <>*/}
        {/*      <div className={`*/}
        {/*                      flex min-h-[300px] w-full*/}
        {/*                      flex-col*/}
        {/*                      rounded-lg*/}
        {/*                      border border-white/10*/}
        {/*                      bg-primary-black*/}
        {/*                      shadow-primary*/}
        {/*                    `}>*/}
        {/*        <div className={`*/}
        {/*                          flex flex-row items-center justify-between*/}
        {/*                          p-4*/}
        {/*                        `}>*/}
        {/*          <span className={`*/}
        {/*            text-sm font-light*/}
        {/*            text-white/50*/}
        {/*          `}>*/}
        {/*            Total files: {files.length}*/}
        {/*          </span>*/}
        {/*          <span className={`*/}
        {/*            text-sm font-light*/}
        {/*            text-white/50*/}
        {/*          `}>*/}
        {/*            Total sizes: {bytesToString(files.reduce((acc, file) => acc + file.size, 0))}*/}
        {/*          </span>*/}
        {/*        </div>*/}
        {/*        <ListFiles*/}
        {/*          files={files}*/}
        {/*          imagesUrls={urls}*/}
        {/*          removeFile={removeFile}*/}
        {/*          submitted={!!instruction}*/}
        {/*          inputFiles={inputFiles}*/}
        {/*          outputFiles={outputFiles}*/}
        {/*        />*/}
        {/*      </div>*/}
        {/*      <div className={`*/}
        {/*                      flex flex-row*/}
        {/*                      gap-2*/}
        {/*                    `}>*/}
        {/*        {showInitialActions && (*/}
        {/*          <>*/}
        {/*            <Button xColor="primary" xSize="sm" onClick={handleSubmit}>*/}
        {/*              Continue*/}
        {/*            </Button>*/}
        {/*            <Button xColor="secondary" xSize="sm" onClick={handleReset}>*/}
        {/*              Reset*/}
        {/*            </Button>*/}
        {/*          </>*/}
        {/*        )}*/}
        {/*        {showFinalActions && (*/}
        {/*          <>*/}
        {/*            <Button xColor="primary" xSize="sm" onClick={goHome}>*/}
        {/*              Back to Homepage*/}
        {/*            </Button>*/}
        {/*            <Button xColor="secondary" xSize="sm" onClick={handleReset}>*/}
        {/*              Reset*/}
        {/*            </Button>*/}
        {/*          </>*/}
        {/*        )}*/}
        {/*      </div>*/}
    </div>
  );
}
