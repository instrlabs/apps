"use client";

import React, { useCallback, useEffect, useMemo, useRef, useState } from "react";

import TextField from "@/components/text-field";
import SearchIcon from "@/components/icons/search";
import { imageTools, pdfTools } from "@/constants/tools";
import { useOverlay } from "@/hooks/useOverlay";
import Chip from "@/components/chip";

type SearchItem = {
  key: string;
  title: string;
  desc: string;
  href: string;
  category: "Image" | "PDF";
};

export default function SearchOverlay() {
  const { closeAll } = useOverlay();
  const [query, setQuery] = useState<string>("");
  const [activeIndex, setActiveIndex] = useState<number>(0);
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

  useEffect(() => {
    // reset active index when results change
    setActiveIndex(0);
  }, [query]);

  const handleSelect = useCallback((item: SearchItem) => {
    if (item.href && item.href !== "#") {
      window.location.href = item.href;
    }
    closeAll();
  }, [closeAll]);

  const onKeyDown = useCallback((e: React.KeyboardEvent) => {
    if (!results.length) return;
    if (e.key === "ArrowDown") {
      e.preventDefault();
      setActiveIndex((prev) => (prev + 1) % results.length);
    } else if (e.key === "ArrowUp") {
      e.preventDefault();
      setActiveIndex((prev) => (prev - 1 + results.length) % results.length);
    } else if (e.key === "Enter") {
      e.preventDefault();
      const item = results[activeIndex];
      if (item) handleSelect(item);
    } else if (e.key === "Escape") {
      e.preventDefault();
      closeAll();
    }
  }, [results, activeIndex, handleSelect, closeAll]);

  return (
    <div className="h-[70vh] flex-col">
      <div className="flex flex-row items-center gap-3 p-4 border-b bg-card">
        <SearchIcon
          className="pointer-events-none w-5 h-5"
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
        <Chip xVariant="outlined" xSize="sm">
          esc
        </Chip>
      </div>
      <div className="max-h-full overflow-y-auto divide-y divide-border">
          {results.map((item, idx) => (
            <div key={`${item.category}:${item.key}`} role="option" aria-selected={idx === activeIndex}>
              <button
                type="button"
                className={`w-full text-left p-3 flex gap-3 items-start hover:bg-hover ${idx === activeIndex ? "bg-hover" : ""}`}
                onClick={() => handleSelect(item)}
              >
                <span
                  className="inline-flex h-6 min-w-6 items-center justify-center rounded-full bg-muted text-xs px-2">
                  {item.category}
                </span>
                <span>
                  <span className="block font-medium">{item.title}</span>
                  <span className="block text-sm text-muted">{item.desc}</span>
                </span>
              </button>
            </div>
          ))}
      </div>
    </div>
  );
}
