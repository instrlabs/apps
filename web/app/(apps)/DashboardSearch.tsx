"use client";

import React from "react";
import Input from "@/components/input";

export default function DashboardSearch() {
  return (
    <div className="w-[500px]">
      <Input
        variant="primary"
        size="lg"
        placeholder="Search"
        leftIconName="search"
        hasLeftIcon
      />
    </div>
  );
}
