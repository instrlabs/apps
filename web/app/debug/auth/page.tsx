'use server';

import { cookies } from "next/headers";
import AuthAction from "@/app/debug/AuthAction";

export default async function DebugAuthPage() {
  const cookieStore = await cookies();
  const accessToken = cookieStore.get('AccessToken');
  const refreshToken = cookieStore.get('RefreshToken');

  return (
    <div className="flex flex-col gap-4 w-lg mx-auto my-10">
      <p className="text-sm truncate">AccessToken: {accessToken ? accessToken.value : 'Not set'}</p>
      <p className="text-sm truncate">RefreshToken: {refreshToken ? refreshToken.value : 'Not set'}</p>
      <AuthAction />
    </div>
  );
}
