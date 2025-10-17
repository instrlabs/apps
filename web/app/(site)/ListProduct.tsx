"use client"

import { useMemo, useState } from "react";
import { useProduct } from "@/hooks/useProduct";
import Input from "@/components/inputs/input";
import Text from "@/components/text";
import AppsCard from "@/components/cards/apps-card";
import { Product } from "@/services/images";
import InputWithIcon from "@/components/inputs/input-with-icon";
import SearchIcon from "@/components/icons/search";

export default function ListProduct() {
  const { productsByType } = useProduct();
  const [query, setQuery] = useState("");

  const images = productsByType["image"];
  const imagesFiltered = useMemo(() => {
    const q = query.trim().toLowerCase();
    if (!q) return images;
    return images.filter(p =>
      [p.title, p.description, p.key]
        .filter(Boolean)
        .some((v: string) => v.toLowerCase().includes(q))
    );
  }, [images, query]);


  return (
    <div className="flex flex-col gap-4">
      <div className="flex w-full justify-center">
        <InputWithIcon
          icon={<SearchIcon className="size-6" />}
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Search"
          aria-label="Search image tools"
          className="min-w-[500px]"
        />
      </div>

      {imagesFiltered.length > 0 && (
        <Text as="h4" xSize="sm" className="font-semibold">
          Images
        </Text>
      )}

      <div className="@container">
        <div className="grid w-full grid-cols-1 gap-2 @2xl:grid-cols-3 @5xl:grid-cols-4">
          {imagesFiltered.map((product: Product) => (
            <AppsCard
              key={product.key}
              href={`/${product.key.split("-").join("/")}`}
              title={product.title}
              description={product.description}
            />
          ))}
        </div>
      </div>
    </div>
  );
}
