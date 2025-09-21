"use client"

import {useCallback, useState} from "react"
import FileDropzone from "@/components/inputs/file-dropzone";
import Button from "@/components/actions/button";
import {
  createInstruction,
  Instruction,
  InstructionFile,
  createInstructionDetails
} from "@/services/images";
import useNotification from "@/hooks/useNotification";
import Notif from "@/app/(product)/images/compress/Notif";
import ListFiles from "@/app/(product)/images/compress/ListFiles";
import Loading from "@/components/feedback/Loading";

type ProgressStatus = 'UPLOAD' | 'UPLOADED' | 'SUBMITTED' | 'LOADING';

const FILE_SIZE_5MB = 5242880;

export default function ImageCompressPage() {
  const [progress, setProgress] = useState<ProgressStatus>('UPLOAD');
  const [files, setFiles] = useState<File[]>([]);
  const [urls, setUrls] = useState<string[]>([]);
  const [instruction, setInstruction] = useState<Instruction | null>(null);
  const [instructionFiles, setInstructionFiles] = useState<InstructionFile[]>([]);

  const { showNotification } = useNotification();

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
        console.log(fileResult);
        if (!fileResult.success || !fileResult.data) {
          showNotification({ title: "Error", message: fileResult.message, type: "error", duration: 3000});
          break;
        } else setInstructionFiles(prev => [...prev, fileResult.data?.file as InstructionFile]);
      }
    }

    setProgress('SUBMITTED');
  }

  return (
    <div className="w-full flex flex-col items-center py-10">
      <Notif />

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

        {progress === 'SUBMITTED' && !!instructionFiles && (
          <div className="w-full space-y-6">
            <p>
              {JSON.stringify(instructionFiles, null, 2)}
            </p>
            {/*<ul className="space-y-3">*/}
            {/*  {instruction.inputs.map((input, idx) => {*/}
            {/*    const outputsByName = new Map(*/}
            {/*      (instruction.outputs || []).map(o => [o.file_name, o])*/}
            {/*    );*/}
            {/*    const matchedOut = outputsByName.get(input.file_name) || instruction.outputs?.[idx];*/}

            {/*    const inSize = input.size;*/}
            {/*    const outSize = matchedOut?.size ?? null;*/}
            {/*    const savedPct = outSize != null ? Math.round((1 - (outSize as number) / inSize) * 100) : null;*/}
            {/*    const isDone = !!matchedOut;*/}

            {/*    const inHref = getFileUrl(instruction.id, input.file_name);*/}
            {/*    const outHref = isDone && matchedOut ? getFileUrl(instruction.id, matchedOut.file_name) : undefined;*/}

            {/*    return (*/}
            {/*      <li key={`row-${idx}`} className="card flex items-center gap-3 p-3">*/}
            {/*        <img*/}
            {/*          src={inHref}*/}
            {/*          alt={input.file_name}*/}
            {/*          width={60}*/}
            {/*          height={60}*/}
            {/*          className="object-cover rounded-lg aspect-square"*/}
            {/*        />*/}

            {/*        <div className="flex-col flex-1 min-w-0">*/}
            {/*          <span className="truncate font-medium">{input.file_name}</span>*/}
            {/*          <div className="flex items-center gap-2 text-sm text-gray-600">*/}
            {/*            <span className="whitespace-nowrap">{bytesToString(inSize)}</span>*/}
            {/*            <span className="text-gray-400">→</span>*/}
            {/*            {isDone ? (*/}
            {/*              <span className="whitespace-nowrap">{bytesToString(outSize || 0)}</span>*/}
            {/*            ) : (*/}
            {/*              <span className="inline-flex items-center gap-2 text-amber-600">*/}
            {/*                <span className="relative block w-20 h-2 rounded bg-amber-100 overflow-hidden">*/}
            {/*                  <span className="absolute inset-0 w-1/2 bg-amber-300 animate-pulse" />*/}
            {/*                </span>*/}
            {/*                Processing...*/}
            {/*              </span>*/}
            {/*            )}*/}
            {/*          </div>*/}
            {/*        </div>*/}

            {/*        <div className="w-24 text-right">*/}
            {/*          {isDone ? (*/}
            {/*            savedPct != null ? (*/}
            {/*              <span className={`${savedPct > 0 ? 'badge bg-green-100 text-green-700' : 'badge bg-gray-100 text-gray-700'}`}>*/}
            {/*                {`${savedPct}%`}*/}
            {/*              </span>*/}
            {/*            ) : (*/}
            {/*              <span className="badge">N/A</span>*/}
            {/*            )*/}
            {/*          ) : (*/}
            {/*            <span className="text-xs text-gray-400">Waiting…</span>*/}
            {/*          )}*/}
            {/*        </div>*/}

            {/*        <div className="w-10 flex justify-end">*/}
            {/*          {isDone && outHref ? (*/}
            {/*            <a*/}
            {/*              href={outHref}*/}
            {/*              download*/}
            {/*              className="inline-flex items-center justify-center w-9 h-9 rounded hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500"*/}
            {/*              title={`Download ${matchedOut?.file_name}`}*/}
            {/*              aria-label={`Download ${matchedOut?.file_name}`}*/}
            {/*            >*/}
            {/*              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" className="w-5 h-5 text-gray-700">*/}
            {/*                <path d="M12 3a1 1 0 011 1v8.586l2.293-2.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 111.414-1.414L11 12.586V4a1 1 0 011-1z" />*/}
            {/*                <path d="M5 20a2 2 0 01-2-2v-1a1 1 0 112 0v1h14v-1a1 1 0 112 0v1a2 2 0 01-2 2H5z" />*/}
            {/*              </svg>*/}
            {/*            </a>*/}
            {/*          ) : null}*/}
            {/*        </div>*/}
            {/*      </li>*/}
            {/*    );*/}
            {/*  })}*/}
            {/*</ul>*/}
          </div>
        )}
      </div>
    </div>
  )
}
