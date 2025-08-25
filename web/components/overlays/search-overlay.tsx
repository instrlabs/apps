"use client";

import React, { useCallback, useMemo, useRef, useState } from "react";

import TextField from "@/components/text-field";
import SearchIcon from "@/components/icons/search";
import { imageTools, pdfTools } from "@/constants/tools";
import { useOverlay } from "@/hooks/useOverlay";
import Chip from "@/components/chip";
import MenuButton from "@/components/menu-button";
import HashtagIcon from "@/components/icons/hashtag";

type SearchItem = {
  key: string;
  title: string;
  desc: string;
  href: string;
};

type Section = {
  key: string;
  title: string;
  items: SearchItem[];
};

export default function SearchOverlay() {
  const { closeAll } = useOverlay();

  const [query, setQuery] = useState<string>("");
  const inputRef = useRef<HTMLInputElement | null>(null);

  const sections: Section[] = useMemo(() => {
    return [
      {
        key: "image",
        title: "Image Tools",
        items: imageTools,
      },
      {
        key: "pdf",
        title: "PDF Tools",
        items: pdfTools,
      },
    ];
  }, []);

  // Apply search only to items within each section (keep sections markup intact)
  const filteredSections: Section[] = useMemo(() => {
    const q = query.trim().toLowerCase();
    if (!q) return sections;
    return sections
      .map((sec) => ({
        ...sec,
        items: sec.items.filter(
          (it) =>
            it.title.toLowerCase().includes(q) ||
            it.desc.toLowerCase().includes(q) ||
            it.key.toLowerCase().includes(q)
        ),
      }))
      .filter((sec) => sec.items.length > 0);
  }, [sections, query]);

  const totalResults = useMemo(
    () => filteredSections.reduce((acc, s) => acc + s.items.length, 0),
    [filteredSections]
  );

  const handleSelect = useCallback((item: SearchItem) => {
    if (item.href && item.href !== "#") {
      window.location.href = item.href;
    }
    closeAll();
  }, [closeAll]);

  const onKeyDown = useCallback((e: React.KeyboardEvent) => {
    if (!totalResults) return;
    if (e.key === "Escape") {
      e.preventDefault();
      closeAll();
    }
  }, [totalResults, closeAll]);

  return (
    <div className="h-[70vh] flex flex-col">
      <div className="sticky top-0 z-10 flex flex-row items-center gap-3 p-5 border-b border-border">
        <SearchIcon className="pointer-events-none w-6 h-6" />
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
        {filteredSections.map((section) => (
          <div key={section.key} className="py-2">
            <div className="px-4 py-3 text-sm font-bold tracking-wide text-muted">
              {section.title}
            </div>
            <div className="px-4 flex flex-col gap-3">
              {section.items.map((item) => (
                <MenuButton
                  key={`${section.key}:${item.key}`}
                  onClick={() => handleSelect(item)}
                >
                  <div className="flex items-center gap-3">
                    <HashtagIcon className="w-5 h-5" />
                    <div className="flex flex-col items-start">
                      <span className="text-sm font-medium">{item.title}</span>
                      <span className="text-xs font-light">{item.desc}</span>
                    </div>
                  </div>
                </MenuButton>
              ))}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
