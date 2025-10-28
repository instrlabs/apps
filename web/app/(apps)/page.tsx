"use client";

import React from "react";
import { useProduct } from "@/hooks/useProduct";
import DashboardCard from "./DashboardCard";
import DashboardSearch from "./DashboardSearch";

export default function AppsPage() {
  const { productsByType } = useProduct();
  const images = productsByType["image"] || [];

  return (
    <div className="flex flex-col gap-2 w-full">
      <div className="flex justify-center w-full">
        <DashboardSearch />
      </div>

      {images.length > 0 && (
        <>
          <h2 className="text-base leading-6 font-semibold text-white">
            Images
          </h2>
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-2 sm:gap-3 lg:gap-4 w-full">
            {images.map((product) => (
              <DashboardCard
                key={product.id}
                title={product.title}
                description={product.description}
                href={`/${product.key.split("-").join("/")}`}
              />
            ))}
          </div>
        </>
      )}
    </div>
  );
}
