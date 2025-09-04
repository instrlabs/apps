"use client";

import React, { createContext, useContext, useMemo, useState } from "react";
import {Product} from "@/services/products";

type ProductContextType = {
  products: Product[];
  loading: boolean;
  productsByType: Record<string, Product[]>;
};

const ProductContext = createContext<ProductContextType | undefined>(undefined);

export function ProductProvider({ children, data }: {
  children: React.ReactNode,
  data: Product[]
}) {
  const [products, setProducts] = useState<Product[]>(data);
  const [loading, setLoading] = useState<boolean>(false);

  const productsByType = useMemo(() => {
    return products.reduce((acc, product) => {
      const { productType } = product;
      if (!acc[productType]) {
        acc[productType] = [];
      }

      acc[productType].push(product);
      return acc;
    }, {} as Record<string, Product[]>);
  }, [products]);

  const value = useMemo(
    () => ({ products, loading, productsByType }),
    [products, loading, productsByType]
  );

  return <ProductContext.Provider value={value}>{children}</ProductContext.Provider>;
}

export function useProduct(): ProductContextType {
  const ctx = useContext(ProductContext);
  if (!ctx) throw new Error("useProfile must be used within a ProfileProvider");
  return ctx;
}
