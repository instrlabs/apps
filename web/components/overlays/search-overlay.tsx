"use client";

import React, { useCallback, useMemo, useRef, useState } from "react";

import TextField from "@/components/text-field";
import SearchIcon from "@/components/icons/search";
import { imageTools, pdfTools } from "@/constants/tools";
import { useOverlay } from "@/hooks/useOverlay";
import Chip from "@/components/chip";
import clsx from "clsx";

type SearchItem = {
  key: string;
  title: string;
  desc: string;
  href: string;
  icon: string;
  category: "Image" | "PDF";
};

export default function SearchOverlay() {
  const { closeAll } = useOverlay();

  const [query, setQuery] = useState<string>("");
  const inputRef = useRef<HTMLInputElement | null>(null);

  const items: SearchItem[] = useMemo(() => {
    return [
      ...imageTools.map((t) => ({ ...t, category: "Image" as const })),
      ...pdfTools.map((t) => ({ ...t, category: "PDF" as const })),
    ];
  }, []);

  const results = useMemo(() => {
    const q = query.trim().toLowerCase();
    if (!q) return items;
    return items.filter((it) =>
      it.title.toLowerCase().includes(q) ||
      it.desc.toLowerCase().includes(q) ||
      it.key.toLowerCase().includes(q)
    );
  }, [items, query]);

  const handleSelect = useCallback((item: SearchItem) => {
    if (item.href && item.href !== "#") {
      window.location.href = item.href;
    }
    closeAll();
  }, [closeAll]);

  const onKeyDown = useCallback((e: React.KeyboardEvent) => {
    if (!results.length) return;
    if (e.key === "Escape") {
      e.preventDefault();
      closeAll();
    }
  }, [results, closeAll]);

  return (
    <div className="h-[70vh] flex flex-col">
      <div className="sticky top-0 z-10 flex flex-row items-center gap-3 p-4 border-b bg-card">
        <SearchIcon
          className="pointer-events-none w-6 h-6"
          aria-hidden="true"
        />
        <TextField
          ref={inputRef}
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          onKeyDown={onKeyDown}
          placeholder="What are you looking for?"
          className="p-0! rounded-none! border-none shadow-none bg-transparent focus:shadow-none"
          autoComplete="off"
        />
        <Chip xVariant="outlined" xSize="sm">esc</Chip>
      </div>
      <div className="flex-1 min-h-0 overflow-y-auto">
        {(["Image", "PDF"] as const).map((cat) => {
          const section = results.filter((it) => it.category === cat);
          if (!section.length) return null;
          const header = cat === "Image" ? "Image Tools" : "PDF Tools";
          return (
            <div key={cat} className="py-2">
              <div className="px-4 py-2 text-xs font-bold uppercase tracking-wide text-muted">
                {header}
              </div>
              <div className="flex flex-col gap-2">
                {section.map((item) => (
                  <button
                    key={`${item.category}:${item.key}`}
                    type="button"
                    className={clsx(
                      "group flex items-center gap-3 rounded-md p-2 text-left text-sm",
                      "bg-gray-50 border border-gray-100",
                      "transition-colors duration-200 ease-in-out",
                      "cursor-pointer"
                    )}
                    onClick={() => handleSelect(item)}
                  >
                    <span className="inline-flex h-6 min-w-6 items-center justify-center text-base">
                      {item.icon}
                    </span>
                    <div className="flex flex-col">
                      <span className="font-medium">{item.title}</span>
                      <span className="font-light">{item.desc}</span>
                    </div>
                  </button>
                ))}
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
