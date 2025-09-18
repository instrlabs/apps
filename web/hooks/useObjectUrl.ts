"use client"

import { useEffect, useState } from "react";

/**
 * Create object URLs for a list of Files and keep them in the same order.
 * URLs are revoked automatically when files change or the component unmounts.
 *
 * @param files Array of File objects
 * @returns Array of object URL strings (same length/order as files)
 */
const useObjectUrl = (files: File[]): string[] => {
  const [urls, setUrls] = useState<string[]>([]);

  useEffect(() => {
    if (!files || files.length === 0) {
      setUrls([]);
      return;
    }

    const newUrls = files.map((f) => URL.createObjectURL(f));
    setUrls(newUrls);

    return () => {
      for (const url of newUrls) URL.revokeObjectURL(url);
    };
  }, [files]);

  return urls;
};

export default useObjectUrl;
export { useObjectUrl };
