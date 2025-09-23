"use client";

import React from "react";


export default function InlineSpinner() {
  return (
    <span
      className={`inline-block animate-spin rounded-full border-2 border-gray-500 border-t-black`.trim()}
      style={{ width: 16, height: 16 }}
    />
  );
}
