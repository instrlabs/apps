'use client'

import Button from "@/components/button";
import { getProfile, login, logout, refresh, sendPin } from "@/services/auth";
import useNotification from "@/hooks/useNotification";
import { useProfile } from "@/hooks/useProfile";

export default function AuthAction() {
  const { showNotification } = useNotification();
  const { profile } = useProfile();

  async function handleSendPin() {
    const res = await sendPin({ email: "arthadede@gmail.com" });
    if (res.success) {
      showNotification({
        type: "info",
        message: "Pin sent to your email",
      })
    } else {
      showNotification({
        type: "error",
        message: res.message,
      })
    }
  }

  async function handleLogin() {
    const res = await login({ email: "arthadede@gmail.com", pin: "000000"});
    if (res.success) {
      showNotification({
        type: "info",
        message: "Login success",
      })
    } else {
      showNotification({
        type: "error",
        message: res.message,
      })
    }
  }

  async function handleLogout() {
    await logout();
  }

  async function handleRefreshToken() {
    await refresh();
  }

  async function handleUpdateProfile() {
    await getProfile();
  }

  return (
    <>
      <p className="text-sm truncate">Profile: {profile?.email}</p>
      <div className="flex flex-col gap-4">
        <Button onClick={handleSendPin} color="primary" size="sm">
          Send Pin
        </Button>
        <Button onClick={handleLogin} color="primary" size="sm">
          Login
        </Button>
        <Button onClick={handleLogout} color="primary" size="sm">
          Logout
        </Button>
        <Button onClick={handleRefreshToken} color="primary" size="sm">
          Refresh Token - Client
        </Button>
        <Button onClick={handleUpdateProfile} color="primary" size="sm">
          Update Profile - Client
        </Button>
      </div>
    </>
  )
}
