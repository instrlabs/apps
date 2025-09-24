"use server"
import { bytesToString } from "@/utils/bytesToString";
import {getImageInstruction} from "@/services/images";
import {APIs} from "@/constants/api";
import Notif from "@/app/(site)/images/compress/Notif";

export default async function Compress({ params }: {
  params: Promise<{ id: string }>
}) {
  const getFileUrl = (id: string, filename: string) => {
    return `${process.env.API_URL}${APIs.IMAGES}/${encodeURIComponent(id)}/${encodeURIComponent(filename)}`;
  }

  const res = await getImageInstruction((await params).id);

  const instruction = res.data;
  const error = res.error;

  return (
    <div className="w-full flex flex-col py-10">
      <Notif />
      <h2 className="text-center text-3xl font-bold mt-6">
        Your images have been compressed!
      </h2>

      <div className="w-full mt-8 flex flex-col items-center space-y-6">
        {error && (
          <div className="card max-w-2xl w-full p-4 text-red-600">
            Failed to load instruction: {error}
          </div>
        )}

        {instruction && (
          <div className="w-full max-w-3xl space-y-6">
            <div className="card p-4">
              <p className="text-sm text-gray-500">Instruction ID</p>
              <p className="font-mono text-sm break-all">{instruction.id}</p>
              <div className="mt-2 flex flex-wrap gap-4 text-sm">
                <span className="badge">Status: {instruction.status}</span>
                <span className="text-gray-500">Created: {new Date(instruction.created_at).toLocaleString()}</span>
                <span className="text-gray-500">Updated: {new Date(instruction.updated_at).toLocaleString()}</span>
              </div>
            </div>

            <div className="">
              <h3 className="font-semibold mb-3">Files</h3>
              {instruction.inputs?.length ? (
                <ul className="space-y-3">
                  {instruction.inputs.map((input, idx) => {
                    const outputsByName = new Map(
                      (instruction.outputs || []).map(o => [o.file_name, o])
                    );
                    const matchedOut = outputsByName.get(input.file_name) || instruction.outputs?.[idx];

                    const inSize = input.size;
                    const outSize = matchedOut?.size ?? null;
                    const savedPct = Math.round((1 - outSize / inSize) * 100);
                    const isDone = !!matchedOut;

                    const inHref = getFileUrl(instruction.id, input.file_name);
                    const outHref = isDone && matchedOut ? getFileUrl(instruction.id, matchedOut.file_name) : undefined;

                    return (
                      <li key={`row-${idx}`} className="card flex items-center gap-3 p-3">
                        <img
                          src={inHref}
                          alt={input.file_name}
                          width={60}
                          height={60}
                          className="object-cover rounded-lg aspect-square"
                        />

                        <div className="flex-col flex-1 min-w-0">
                          <span className="truncate font-medium">{input.file_name}</span>
                          <div className="flex items-center gap-2 text-sm text-gray-600">
                            <span className="whitespace-nowrap">{bytesToString(inSize)}</span>
                            <span className="text-gray-400">→</span>
                            {isDone ? (
                              <span className="whitespace-nowrap">{bytesToString(outSize || 0)}</span>
                            ) : (
                              <span className="inline-flex items-center gap-2 text-amber-600">
                              <span className="relative block w-20 h-2 rounded bg-amber-100 overflow-hidden">
                                <span className="absolute inset-0 w-1/2 bg-amber-300 animate-pulse" />
                              </span>
                              Processing...
                            </span>
                            )}
                          </div>
                        </div>

                        <div className="w-24 text-right">
                          {isDone ? (
                            savedPct != null ? (
                              <span className={`badge ${savedPct > 0 ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-700'}`}>
                                {`${savedPct}%`}
                              </span>
                            ) : (
                              <span className="badge">N/A</span>
                            )
                          ) : (
                            <span className="text-xs text-gray-400">Waiting…</span>
                          )}
                        </div>

                        <div className="w-10 flex justify-end">
                          {isDone && outHref ? (
                            <a
                              href={outHref}
                              download
                              className="inline-flex items-center justify-center w-9 h-9 rounded hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
                              title={`Download ${matchedOut?.file_name}`}
                              aria-label={`Download ${matchedOut?.file_name}`}
                            >
                              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" className="w-5 h-5 text-gray-700">
                                <path d="M12 3a1 1 0 011 1v8.586l2.293-2.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 111.414-1.414L11 12.586V4a1 1 0 011-1z" />
                                <path d="M5 20a2 2 0 01-2-2v-1a1 1 0 112 0v1h14v-1a1 1 0 112 0v1a2 2 0 01-2 2H5z" />
                              </svg>
                            </a>
                          ) : null}
                        </div>
                      </li>
                    );
                  })}
                </ul>
              ) : (
                <p className="text-sm text-gray-500">No input files found.</p>
              )}

              {(!instruction.outputs || instruction.outputs.length === 0) && instruction.status !== 'completed' && (
                <p className="mt-3 text-xs text-gray-500">Outputs are being generated. This may take a moment…</p>
              )}
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
