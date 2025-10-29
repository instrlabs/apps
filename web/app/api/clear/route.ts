import { redirect, RedirectType } from "next/navigation";
import { cookies } from "next/headers";

export async function GET() {
  const storeCookie = await cookies();
  storeCookie.delete("access_token");
  storeCookie.delete("refresh_token");
  redirect("/login", RedirectType.push);
}
