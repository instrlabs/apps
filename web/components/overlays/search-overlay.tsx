import React from "react";

import TextField from "@/components/text-field";
import SearchIcon from "@/components/icons/search";

export default function SearchOverlay() {
  return (
    <div className="space-y-3">
      <div className="relative">
        <label htmlFor="global-search-input" className="sr-only">Search</label>
        <TextField
          id="global-search-input"
          type="text"
          autoFocus
          placeholder="Search..."
          className="pr-10"
          xSize="md"
        />
        <SearchIcon
          className="pointer-events-none absolute right-3 top-1/2 -translate-y-1/2 w-5 h-5 text-muted"
          aria-hidden="true"
        />
      </div>
      <div className="text-sm text-muted">Type to search…</div>
    </div>
  );
}
