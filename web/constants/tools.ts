export type ToolItem = {
  key: string;
  title: string;
  desc: string;
  href: string;
};

export const imageToolsDict: Record<string, ToolItem> = {
  compress: {
    key: "compress",
    title: "Compress",
    desc: "Reduce image file size while preserving quality.",
    href: "/apps/image-compress",
  },
  resize: {
    key: "resize",
    title: "Resize",
    desc: "Change width/height, keep aspect ratio or set exact.",
    href: "#",
  },
  crop: {
    key: "crop",
    title: "Crop",
    desc: "Trim images to a selected region or aspect ratio.",
    href: "#",
  },
  convert: {
    key: "convert",
    title: "Convert",
    desc: "Convert between JPG, PNG, WEBP, and more.",
    href: "#",
  },
  rotate: {
    key: "rotate",
    title: "Rotate",
    desc: "Rotate images by 90°, 180°, or custom angles.",
    href: "#",
  },
  flip: {
    key: "flip",
    title: "Flip",
    desc: "Flip images horizontally or vertically.",
    href: "#",
  },
  watermark: {
    key: "watermark",
    title: "Watermark",
    desc: "Add text or image watermarks with controls.",
    href: "#",
  },
  optimize: {
    key: "optimize",
    title: "Optimize",
    desc: "Auto-optimize images for web performance.",
    href: "#",
  },
};

// PDF tools as a filename-keyed dictionary
export const pdfToolsDict: Record<string, ToolItem> = {
  "compress-pdf": {
    key: "compress-pdf",
    title: "Compress PDF",
    desc: "Reduce PDF size while keeping quality readable.",
    href: "#",
  },
  "merge-pdf": {
    key: "merge-pdf",
    title: "Merge PDFs",
    desc: "Combine multiple PDFs into a single file.",
    href: "#",
  },
  "split-pdf": {
    key: "split-pdf",
    title: "Split PDF",
    desc: "Extract pages or split by ranges.",
    href: "#",
  },
  "rotate-pdf": {
    key: "rotate-pdf",
    title: "Rotate PDF",
    desc: "Rotate pages by 90°, 180°, or custom angles.",
    href: "#",
  },
  "reorder-pages": {
    key: "reorder-pages",
    title: "Reorder Pages",
    desc: "Change the order of pages quickly.",
    href: "#",
  },
  protect: {
    key: "protect",
    title: "Protect PDF",
    desc: "Add a password to restrict access.",
    href: "#",
  },
  unlock: {
    key: "unlock",
    title: "Unlock PDF",
    desc: "Remove password from a PDF you own.",
    href: "#",
  },
  "extract-text": {
    key: "extract-text",
    title: "Extract Text",
    desc: "Pull text content from PDF pages.",
    href: "#",
  },
  "extract-images": {
    key: "extract-images",
    title: "Extract Images",
    desc: "Export embedded images from PDF.",
    href: "#",
  },
  "pdf-to-images": {
    key: "pdf-to-images",
    title: "PDF to Images",
    desc: "Render pages to JPG/PNG/WebP.",
    href: "#",
  },
  "convert-to-pdf": {
    key: "convert-to-pdf",
    title: "Convert to PDF",
    desc: "Turn images/Docs into a single PDF.",
    href: "#",
  },
};

export const imageTools = Object.values(imageToolsDict);
export const pdfTools = Object.values(pdfToolsDict);
