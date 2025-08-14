import {ReactNode} from "react";

export default function OverlayContent({ children }: { children: ReactNode }) {
  return (
    <div
      className="absolute top-0 bottom-0 p-3 pt-[90px] transition-[left,right] duration-300 ease-in-out"
      style={{
        left: 'var(--overlay-left-width, 300px)',
        right: 'var(--overlay-right-width, 300px)'
      }}
    >
      <div className="w-full h-full rounded-3xl bg-neutral-50">
        <div className="h-full w-full flex items-center justify-center text-gray-700">
          {children}
        </div>
      </div>
    </div>
  )
}
