import ImagePreview from "@/components/ImagePreview";
import { bytesToString } from "@/utils/bytesToString";
import CloseIcon from "@/components/icons/CloseIcon";
import Chip from "@/components/chip";

export default function ListFiles(props: {
  files: File[];
  imagesUrls: string[];
  removeFile: (index: number) => void;
}) {
  return (
    <div className="flex w-full flex-col px-4">
      {props.files.map((f, idx) => {
        return (
          <div
            key={idx}
            className="card flex w-full items-center justify-between gap-4 rounded-lg p-2 transition-colors hover:bg-white/10"
          >
            <ImagePreview src={props.imagesUrls[idx]} alt={f.name} width={40} height={40} />
            <div className="min-w-0 flex-1">
              <p className="truncate text-sm">{f.name}</p>
              <div className="flex flex-row gap-2">
                <p className="truncate text-sm font-light text-white/50">{bytesToString(f.size)}</p>
                ---
                <p className="truncate text-sm">{bytesToString(f.size)}</p>
              </div>
            </div>
            {/* CLOSE ICON ONLY SHOW WHEN NOT SUBMITTED */}
            <CloseIcon className="ml-auto size-4 shrink-0 cursor-pointer text-white/70" />
            {/* FOR DISPLAY STATUS INPUT FILE AND OUTPUT FILE */}
            <Chip xSize="sm" xVariant="outline" xColor="warning" loading>Processing...</Chip>
          {/*  ADD DOWNLOAD ICON*/}
          </div>
        );
      })}
    </div>
  );
}
