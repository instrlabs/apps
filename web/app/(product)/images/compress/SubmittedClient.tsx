"use client";

import { InstructionFile} from "@/services/images";
import {bytesToString} from "@/utils/bytesToString";

export default function SubmittedClient(props: {
  instructionId: string;
  inputFiles: InstructionFile[];
  outputFiles: InstructionFile[];
}) {
  const inputs = props.inputFiles;
  const outputs = props.outputFiles;

  return (
    <div className="w-full max-w-2xl space-y-4">
      {inputs.map((f) => {
        const out = outputs.find(o => o.id === f.output_id);
        const isDone = !!out || (f.status || "").toUpperCase() === "DONE" || (f.output_id && f.status === "COMPLETED");
        const outSize = out?.size;
        const savedPct = isDone && outSize != null ? Math.round((1 - (outSize / f.size)) * 100) : null;
        const imgId = f.id; // always show input thumbnail for the input section
        const fileUrl = `http://localhost:3000/images/instructions/${props.instructionId}/details/${imgId}`;
        const outputUrl = out ? `http://localhost:3000/images/instructions/${props.instructionId}/details/${out.id}` : undefined;
        return (
          <div key={f.id} className="w-full flex items-center justify-between card p-3 gap-4">
            <img
              src={fileUrl}
              alt={f.original_name}
              width={60}
              height={60}
            />
            <div className="min-w-0 flex-1">
              <p className="truncate font-medium">{f.original_name}</p>
              <div className="flex items-center gap-2 text-sm text-gray-600">
                <span className="whitespace-nowrap">{bytesToString(f.size)}</span>
                <span className="text-gray-400">→</span>
                {isDone ? (
                  <span className="whitespace-nowrap">{outSize != null ? bytesToString(outSize) : '—'}</span>
                ) : (
                  <span className="inline-flex items-center gap-2 text-amber-600">
                    <span className="relative block w-20 h-2 rounded bg-amber-100 overflow-hidden">
                      <span className="absolute inset-0 w-1/2 bg-amber-300 animate-pulse" />
                    </span>
                    Processing...
                  </span>
                )}
              </div>
              <p className="text-xs text-gray-400 mt-1">Status: {f.status}</p>
            </div>

            <div className="w-40 text-right">
              {out && outputUrl ? (
                <div className="flex flex-col items-end gap-1">
                  <span className="text-sm text-gray-700">{bytesToString(out.size)}</span>
                  <span className="text-xs text-gray-500">Status: {out.status}</span>
                  <a
                    href={outputUrl}
                    download
                    className="text-xs text-blue-600 hover:underline"
                  >
                    Download
                  </a>
                </div>
              ) : (
                <div className="flex flex-col items-end gap-1">
                  <span className="text-sm text-gray-400">—</span>
                  <span className="text-xs text-gray-400">Status: {f.status}</span>
                  <span className="text-[10px] text-gray-400">Waiting…</span>
                </div>
              )}
            </div>
          </div>
        );
      })}
    </div>
  );
}
