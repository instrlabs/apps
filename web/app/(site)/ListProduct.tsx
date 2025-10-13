"use client"

import {useProduct} from "@/hooks/useProduct";
import Link from "next/link";

export default function ListProduct() {
  const { productsByType } = useProduct();
  return (
    <div className="p-4">
      <h4 className="mb-4">
        Image Tools
      </h4>
      <div className="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-4">
        {productsByType['image']?.map((product) => (
          <Link
            key={product.key}
            href={`/${product.key.split('-').join('/')}`}
          >
            <div className="p-4 rounded-lg shadow-primary bg-primary-black flex flex-col gap-1">
              <h3 className="text-sm">{product.title}</h3>
              <p className="text-white/50 font-light text-sm">{product.description}</p>
            </div>
          </Link>
        ))}
      </div>
    </div>
  )
}
