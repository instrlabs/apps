import { getImageInstructions, type ImageInstruction } from "@/services/images";
import { APIs } from "@/constants/api";

export default async function HistoriesPage() {
  const res = await getImageInstructions();

  return (
    <div className="p-10">
      <h1 className="text-2xl font-semibold mb-4">Histories</h1>

      {res.success && res.data && res.data.length > 0 && (
        <>
          <h4 className="text-xl font-medium">Images</h4>
        </>
      )}

      {res.success && res.data && res.data.length > 0 && (
        <ul className="space-y-3">
          {res.data.map((item: ImageInstruction) => (
            <li key={item.id} className="rounded border p-4 flex flex-col gap-3">
              <div className="flex items-start justify-between gap-4">
                <div className="min-w-0">
                  <div className="font-medium truncate">ID: {item.id}</div>
                  <div className="text-sm text-gray-600">Status: <span className="font-mono uppercase">{item.status}</span></div>
                  <div className="text-sm text-gray-600">Created: {new Date(item.created_at).toLocaleString()}</div>
                  <div className="text-sm text-gray-600">Updated: {new Date(item.updated_at).toLocaleString()}</div>
                </div>
                <div className="flex flex-col items-end text-sm text-gray-700 shrink-0">
                  <div>Inputs: {item.inputs?.length ?? 0}</div>
                  <div>Outputs: {item.outputs?.length ?? 0}</div>
                </div>
              </div>

              {/* Inputs list */}
              {item.inputs && item.inputs.length > 0 && (
                <div>
                  <div className="text-sm font-medium text-gray-800 mb-1">Inputs</div>
                  <ul className="text-sm text-gray-700 grid gap-1">
                    {item.inputs.map((f, idx) => {
                      const href = `${process.env.API_URL}${APIs.IMAGE_INSTRUCTIONS}/${encodeURIComponent(item.id)}/${encodeURIComponent(f.file_name)}`;
                      return (
                        <li key={`${item.id}-in-${idx}`} className="flex items-center justify-between gap-2">
                          <span className="truncate">{f.file_name}</span>
                          <span className="flex items-center gap-3 shrink-0">
                            <span className="text-gray-500">{Math.round(f.size / 1024)} KB</span>
                            <a
                              href={href}
                              download
                              className="inline-flex items-center gap-1 rounded border px-2 py-1 text-xs hover:bg-gray-50"
                              title="Download file"
                            >
                              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" className="w-4 h-4">
                                <path d="M12 3a1 1 0 0 1 1 1v8.586l2.293-2.293a1 1 0 1 1 1.414 1.414l-4 4a1 1 0 0 1-1.414 0l-4-4A1 1 0 1 1 8.707 10.293L11 12.586V4a1 1 0 0 1 1-1ZM5 19a1 1 0 1 0 0 2h14a1 1 0 1 0 0-2H5Z" />
                              </svg>
                              <span>Download</span>
                            </a>
                          </span>
                        </li>
                      );
                    })}
                  </ul>
                </div>
              )}

              {/* Outputs list with download */}
              {item.outputs && item.outputs.length > 0 && (
                <div>
                  <div className="text-sm font-medium text-gray-800 mb-1">Outputs</div>
                  <ul className="text-sm text-gray-700 grid gap-1">
                    {item.outputs.map((f, idx) => {
                      const href = `${process.env.API_URL}${APIs.IMAGE_INSTRUCTIONS}/${encodeURIComponent(item.id)}/${encodeURIComponent(f.file_name)}`;
                      return (
                        <li key={`${item.id}-out-${idx}`} className="flex items-center justify-between gap-2">
                          <span className="truncate">{f.file_name}</span>
                          <span className="flex items-center gap-3 shrink-0">
                            <span className="text-gray-500">{Math.round(f.size / 1024)} KB</span>
                            <a
                              href={href}
                              download
                              className="inline-flex items-center gap-1 rounded border px-2 py-1 text-xs hover:bg-gray-50"
                              title="Download file"
                            >
                              {/* Download icon */}
                              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" className="w-4 h-4">
                                <path d="M12 3a1 1 0 0 1 1 1v8.586l2.293-2.293a1 1 0 1 1 1.414 1.414l-4 4a1 1 0 0 1-1.414 0l-4-4A1 1 0 1 1 8.707 10.293L11 12.586V4a1 1 0 0 1 1-1ZM5 19a1 1 0 1 0 0 2h14a1 1 0 1 0 0-2H5Z" />
                              </svg>
                              <span>Download</span>
                            </a>
                          </span>
                        </li>
                      );
                    })}
                  </ul>
                </div>
              )}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}

