"use server";

import { EventSource } from "@/utils/eventsource";

export default async function DebugSSEPage() {
  const es = EventSource("http://localhost:3001" + "/sse");
  es.onmessage = (ev) => { console.log(ev); };
 return (
   <h1>Testing</h1>
 )
}
