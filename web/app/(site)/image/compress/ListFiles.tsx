import ImagePreview from "@/components/ImagePreview";
import { bytesToString } from "@/utils/bytesToString";
import CloseIcon from "@/components/icons/CloseIcon";
import Chip from "@/components/chip";
import ButtonIcon from "@/components/actions/button-icon";
import DownloadIcon from "@/components/icons/DownloadIcon";
import { InstructionFile, getInstructionFileBytes } from "@/services/images";
import Button from "@/components/actions/button";

export default function ListFiles(props: {
  files: File[];
  imagesUrls: string[];
  removeFile: (index: number) => void;
  submitted?: boolean;
  inputFiles?: InstructionFile[];
  outputFiles?: InstructionFile[];
}) {
  const handleDownload = async (output: InstructionFile) => {
    const res = await getInstructionFileBytes(output.instruction_id, output.id);
    if (!res?.success || !res.data) {
      console.error("Failed to download file:", res?.message);
      return;
    }

    const filename = output.original_name;
    const lower = filename.toLowerCase();
    let mime = "application/octet-stream";
    if (lower.endsWith(".png")) mime = "image/png";
    else if (lower.endsWith(".jpg") || lower.endsWith(".jpeg")) mime = "image/jpeg";
    else if (lower.endsWith(".webp")) mime = "image/webp";
    else if (lower.endsWith(".gif")) mime = "image/gif";
    else if (lower.endsWith(".svg")) mime = "image/svg+xml";

    const blob = new Blob([res.data], { type: mime });
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = filename;
    document.body.appendChild(a);
    a.click();
    a.remove();
    window.URL.revokeObjectURL(url);
  };

  return (
    <div className="flex w-full flex-col px-4">
      {props.files.map((f, idx) => {
        const input = props.inputFiles?.[idx] || null;
        const output = props.outputFiles?.[idx] || null;
        const isInputFailed = input?.status === "FAILED";
        const isInputPending = input?.status === "PENDING";
        const isInputUploading = input?.status === "UPLOADING";
        const isInputProcessing = input?.status === "PROCESSING";
        const isInputDone = input?.status === "DONE";
        const isOutputFailed = output?.status === "FAILED";
        const isOutputPending = output?.status === "PENDING";
        const isOutputUploading = output?.status === "UPLOADING";
        const isOutputProcessing = output?.status === "PROCESSING";
        const isOutputDone = output?.status === "DONE";

        return (
          <div
            key={idx}
            className="card flex w-full items-center justify-between gap-4 rounded-lg p-2 transition-colors hover:bg-white/10"
          >
            <ImagePreview src={props.imagesUrls[idx]} alt={f.name} width={40} height={40} />
            <div className="min-w-0 flex-1">
              <p className="truncate text-sm">{f.name}</p>
              <div className="flex flex-row gap-1">
                <p className="truncate text-sm font-light text-white/50">{bytesToString(f.size)}</p>
                {output && (
                  <>
                    <svg xmlns="http://www.w3.org/2000/svg" height="20px" viewBox="0 -960 960 960" width="20px" fill="#FFFFFF"><path d="m569.85-301.85-37.16-36.38L648.46-454H212v-52h436.46L532.69-621.77l37.16-36.38L748-480 569.85-301.85Z"/></svg>
                    <p className="truncate text-sm text-white">{bytesToString(output?.size)}</p>
                  </>
                )}
              </div>
            </div>
            {/* CLOSE ICON ONLY SHOW WHEN NOT SUBMITTED */}
            {!props.submitted && (
              <Button onClick={() => props.removeFile(idx)} xVariant="transparent">
                <CloseIcon className="ml-auto size-4 shrink-0 cursor-pointer text-white/70" />
              </Button>
            )}
            {/* FOR DISPLAY STATUS INPUT FILE AND OUTPUT FILE */}
            {props.submitted && (() => {
              const failed = isInputFailed || isOutputFailed;
              const uploading = isInputUploading || isOutputUploading;
              const pending = isInputPending || isOutputPending;
              const processing = isInputProcessing || isOutputProcessing || (isInputDone && !isOutputDone);

              if (failed) return <Chip xSize="sm" xVariant="outline" xColor="danger">Failed</Chip>;
              if (uploading) return <Chip xSize="sm" xVariant="outline" xColor="warning" loading>Uploading...</Chip>;
              if (pending) return <Chip xSize="sm" xVariant="outline" xColor="info" loading>Pending...</Chip>;
              if (processing) return <Chip xSize="sm" xVariant="outline" xColor="warning" loading>Processing...</Chip>;
            })()}
          {/*  ADD DOWNLOAD ICON*/}
            {props.submitted && isOutputDone && (
              <ButtonIcon
                xsize="sm"
                xVariant="solid"
                onClick={() => handleDownload(output)}
              >
                <DownloadIcon className="size-6" />
              </ButtonIcon>
            )}
          </div>
        );
      })}
    </div>
  );
}
