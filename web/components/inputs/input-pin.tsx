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
    "w-12 h-12 text-center text-base rounded bg-white/10 border border-white/30 text-white placeholder:text-white/40 focus:outline-none focus:border-white focus:bg-white/15 transition-colors";

  return (
    <div className={"flex justify-between"}>
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
