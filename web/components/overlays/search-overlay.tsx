"use client";

import React, { useCallback, useMemo, useRef, useState } from "react";

import TextField from "@/components/inputs/text-field";
import SearchIcon from "@/components/icons/search";
import { useOverlay } from "@/hooks/useOverlay";
import Chip from "@/components/chip";
import MenuButton from "@/components/actions/menu-button";
import HashtagIcon from "@/components/icons/hashtag";
import {useProduct} from "@/hooks/useProduct";
import {Product} from "@/services/images";

type Section = {
  key: string;
  title: string;
  items: Product[];
};

export default function SearchOverlay() {
  const { closeRight } = useOverlay();

  const [query, setQuery] = useState<string>("");
  const inputRef = useRef<HTMLInputElement | null>(null);

  const { productsByType } = useProduct();

  const sections: Section[] = useMemo(() => {
    return [
      {
        key: "image",
        title: "Image Tools",
        items: productsByType["image"],
      },
      {
        key: "pdf",
        title: "PDF Tools",
        items: productsByType["pdf"],
      },
    ];
  }, [productsByType]);

  const filteredSections: Section[] = useMemo(() => {
    const q = query.trim().toLowerCase();
    if (!q) return sections;
    return sections
      .map((sec) => ({
        ...sec,
        items: sec.items.filter((it) =>
          it.title.toLowerCase().includes(q) ||
          it.key.toLowerCase().includes(q)),
      }))
      .filter((sec) => sec.items.length > 0);
  }, [sections, query]);

  const totalResults = useMemo(
    () => filteredSections.reduce((acc, s) => acc + s.items.length, 0),
    [filteredSections]
  );

  const handleClick = useCallback((item: Product) => {
    closeRight();
  }, [closeRight]);

  const onKeyDown = useCallback((e: React.KeyboardEvent) => {
    if (!totalResults) return;
    if (e.key === "Escape") {
      e.preventDefault();
      closeRight();
    }
  }, [totalResults, closeRight]);

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
          className="p-0! rounded-none! border-none !shadow-none focus:shadow-none"
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
                  onClick={() => handleClick(item)}
                >
                  <div className="flex items-center gap-3">
                    <HashtagIcon className="w-5 h-5" />
                    <div className="flex flex-col items-start">
                      <span className="text-sm font-medium">{item.title}</span>
                      <span className="text-xs font-light">{item.description}</span>
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
