"use client";

import { imageTools, pdfTools } from "@/constants/tools";

export default function AppsPage() {
  return (
    <div className="container mx-auto p-4">
      <div className="mb-6">
        <h1 className="text-2xl font-semibold">Image Tools</h1>
        <p className="text-sm text-muted">Quick utilities to process your images.</p>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
        {imageTools.map((tool) => (
          <a
            key={tool.key}
            href={tool.href}
            className="group block rounded-xl bg-card border border-border p-4"
          >
            <div className="flex items-start gap-3">
              <div className="text-2xl leading-none select-none">
                <span aria-hidden>{tool.icon}</span>
              </div>
              <div>
                <h2 className="font-medium group-hover:text-primary">{tool.title}</h2>
                <p className="mt-1 text-sm text-muted">{tool.desc}</p>
              </div>
            </div>
          </a>
        ))}
      </div>

      <div className="mt-10 mb-6">
        <h1 className="text-2xl font-semibold">PDF Tools</h1>
        <p className="text-sm text-muted">Handy utilities to work with PDF files.</p>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
        {pdfTools.map((tool) => (
          <a
            key={tool.key}
            href={tool.href}
            className="group block rounded-xl bg-card border border-border p-4"
          >
            <div className="flex items-start gap-3">
              <div className="text-2xl leading-none select-none">
                <span aria-hidden>{tool.icon}</span>
              </div>
              <div>
                <h2 className="font-medium group-hover:text-primary">{tool.title}</h2>
                <p className="mt-1 text-sm text-muted">{tool.desc}</p>
              </div>
            </div>
          </a>
        ))}
      </div>
    </div>
  );
}
