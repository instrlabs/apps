"use client";

import React, { useCallback, useEffect, useMemo, useRef, useState } from "react";

import TextField from "@/components/text-field";
import SearchIcon from "@/components/icons/search";
import { imageTools, pdfTools } from "@/constants/tools";
import { useOverlay } from "@/hooks/useOverlay";

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
    <div className="space-y-3">
      <div className="relative">
        <label htmlFor="global-search-input" className="sr-only">Search</label>
        <TextField
          id="global-search-input"
          type="text"
          autoFocus
          placeholder="Search tools…"
          className="pr-10"
          xSize="md"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          onKeyDown={onKeyDown}
          ref={inputRef}
          aria-autocomplete="list"
          role="combobox"
          aria-expanded={results.length > 0}
          aria-controls="global-search-results"
        />
        <SearchIcon
          className="pointer-events-none absolute right-3 top-1/2 -translate-y-1/2 w-5 h-5 text-muted"
          aria-hidden="true"
        />
      </div>

      {results.length === 0 ? (
        <div className="text-sm text-muted">No matches. Try different keywords.</div>
      ) : (
        <ul
          id="global-search-results"
          role="listbox"
          className="divide-y divide-border rounded-xl border border-border overflow-hidden"
        >
          {results.map((item, idx) => (
            <li key={`${item.category}:${item.key}`} role="option" aria-selected={idx === activeIndex}>
              <button
                type="button"
                className={`w-full text-left p-3 flex gap-3 items-start hover:bg-hover ${idx === activeIndex ? "bg-hover" : ""}`}
                onClick={() => handleSelect(item)}
              >
                <span className="inline-flex h-6 min-w-6 items-center justify-center rounded-full bg-muted text-xs px-2">
                  {item.category}
                </span>
                <span>
                  <span className="block font-medium">{item.title}</span>
                  <span className="block text-sm text-muted">{item.desc}</span>
                </span>
              </button>
            </li>
          ))}
        </ul>
      )}

      {query.trim().length === 0 && (
        <div className="text-sm text-muted">Type to search image and PDF tools. Use ↑/↓ to navigate, Enter to open, Esc to close.</div>
      )}
    </div>
  );
}
