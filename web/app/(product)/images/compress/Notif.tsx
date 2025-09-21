"use client"

import useSSE from "@/hooks/useSSE";
import { useEffect } from "react";

export default function Notif() {
  const { message } = useSSE()

  useEffect(() => {
    console.log(message);
  }, [message]);

  return <div />
}
