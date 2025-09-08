"use client"

import {useProduct} from "@/hooks/useProduct";
import Link from "next/link";

export default function ListProduct() {
  const { productsByType } = useProduct();

  return (
    <div className="container mx-auto p-10">
      <h4 className="text-xl font-bold mb-4">Image Tools</h4>
      <div className="grid grid-cols-4 gap-4">
        {productsByType['image']?.map((product) => (
          <Link
            key={product.key}
            href={`/${product.key}`}
          >
            <div
              className={
                `p-6 border border-border rounded-lg shadow-primary hover:bg-gray-100 `
              }
            >
              <h3 className="text-lg font-semibold">{product.name}</h3>
              <p className="text-gray-600">{product.description}</p>
            </div>
          </Link>
        ))}
      </div>
    </div>
  )
}
