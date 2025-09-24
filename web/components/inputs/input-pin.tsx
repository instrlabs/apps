"use client";

import React from "react";

export type InputPinProps = {
  values: string[];
  onChange: (values: string[]) => void;
  length?: number;
};

export default function InputPin({
  values,
  onChange,
  length = 6,
}: InputPinProps) {
  const inputsRef = React.useRef<Array<HTMLInputElement | null>>([]);

  React.useEffect(() => {
    if (inputsRef.current.length !== length) {
      inputsRef.current = Array(length).fill(null);
    }
  }, [length]);

  const setRef = (el: HTMLInputElement | null, idx: number) => {
    inputsRef.current[idx] = el;
  };

  const handleChange = (idx: number) => (e: React.ChangeEvent<HTMLInputElement>) => {
    const val = e.target.value.replace(/\D/g, "");
    const next = [...values];
    if (!val) {
      next[idx] = "";
      onChange(next);
      return;
    }
    next[idx] = val[0];
    onChange(next);
    if (idx < length - 1) {
      inputsRef.current[idx + 1]?.focus();
    }
  };

  const handleKeyDown = (idx: number) => (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Backspace" && !values[idx] && idx > 0) {
      inputsRef.current[idx - 1]?.focus();
    }
    if (e.key === "ArrowLeft" && idx > 0) {
      e.preventDefault();
      inputsRef.current[idx - 1]?.focus();
    }
    if (e.key === "ArrowRight" && idx < length - 1) {
      e.preventDefault();
      inputsRef.current[idx + 1]?.focus();
    }
  };

  const handlePaste = (e: React.ClipboardEvent<HTMLInputElement>) => {
    e.preventDefault();
    const txt = e.clipboardData.getData("text").replace(/\D/g, "");
    if (!txt) return;
    const next = [...values];
    for (let i = 0; i < length && i < txt.length; i++) {
      next[i] = txt[i];
    }
    onChange(next);
    const nextIdx = Math.min(txt.length, length - 1);
    inputsRef.current[nextIdx]?.focus();
  };

  const baseInputClass =
    "w-10 h-12 sm:w-12 sm:h-14 text-center text-lg rounded-lg bg-primary-black shadow-primary focus:outline-none focus:ring-2 focus:ring-white";

  return (
    <div className={"flex gap-2 sm:gap-3"}>
      {Array.from({ length }).map((_, i) => (
        <input
          key={i}
          ref={(el) => setRef(el, i)}
          type="text"
          inputMode="numeric"
          maxLength={1}
          className={baseInputClass}
          value={values[i] || ""}
          onChange={handleChange(i)}
          onKeyDown={handleKeyDown(i)}
          onPaste={handlePaste}
        />
      ))}
    </div>
  );
}
