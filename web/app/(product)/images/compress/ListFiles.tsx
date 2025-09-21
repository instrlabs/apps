import ImagePreview from "@/components/ImagePreview";
import {bytesToString} from "@/utils/bytesToString";
import CloseIcon from "@/components/icons/CloseIcon";

export default function ListFiles(props: {
  files: File[];
  imagesUrls: string[];
  removeFile: (index: number) => void;
}) {
  return (
    <div className="w-full max-w-2xl space-y-4">
      {props.files.map((f, idx) => {
        return (
          <div
            key={idx}
            className="w-full flex items-center justify-between card p-3 gap-4"
          >
            <ImagePreview
              src={props.imagesUrls[idx]}
              alt={f.name}
              width={60}
              height={60}
            />
            <div className="min-w-0 flex-1">
              <p className="truncate font-medium">{f.name}</p>
              <p className="truncate font-light">{bytesToString(f.size)}</p>
            </div>
            <button
              className="text-sm text-red-400 hover:text-red-700 cursor-pointer"
              onClick={(e) => {
                e.stopPropagation()
                props.removeFile(idx)
              }}
            >
              <CloseIcon />
            </button>
          </div>
        )
      })}
    </div>
  )
}
