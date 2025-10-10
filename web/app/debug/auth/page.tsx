'use server';

import { cookies } from "next/headers";
import AuthAction from "@/app/debug/auth/AuthAction";

export default async function DebugAuthPage() {
  const cookieStore = await cookies();
  const accessToken = cookieStore.get('access_token');
  const refreshToken = cookieStore.get('refresh_token');

  return (
    <div className="flex flex-col gap-4 w-lg mx-auto my-10">
      <p className="text-sm truncate">AccessToken: {accessToken ? accessToken.value : 'Not set'}</p>
      <p className="text-sm truncate">RefreshToken: {refreshToken ? refreshToken.value : 'Not set'}</p>
      <AuthAction />
    </div>
  );
}
