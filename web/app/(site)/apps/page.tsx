"use client";

import { imageTools } from "@/constants/tools";
import MenuButton from "@/components/menu-button";
import HashtagIcon from "@/components/icons/hashtag";

export default function AppsPage() {
  return (
    <div className="container mx-auto">
      <div className="p-6">
        <div className="mb-4">
          <h1 className="text-lg font-bold">Image Tools</h1>
        </div>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
          {imageTools.map((tool) => (
            <MenuButton key={tool.key} xSize="lg">
              <div className="flex flex-col items-start gap-2">
                <div className="flex items-center gap-2">
                  <HashtagIcon className="w-5 h-5" />
                  <h2 className="text-base font-semibold">{tool.title}</h2>
                </div>
                <span className="text-sm font-light text-wrap text-left">{tool.desc}</span>
              </div>
            </MenuButton>
          ))}
        </div>
      </div>
    </div>
  );
}
