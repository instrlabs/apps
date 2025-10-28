'use client'

import Button from "@/components/button";
import { getProfile, login, logout, refresh, sendPin } from "@/services/auth";
import useSnackbar from "@/hooks/useSnackbar";
import { useProfile } from "@/hooks/useProfile";

export default function AuthAction() {
  const { showSnackbar } = useSnackbar();
  const { profile } = useProfile();

  async function handleSendPin() {
    const res = await sendPin({ email: "arthadede@gmail.com" });
    if (res.success) {
      showSnackbar({
        type: "info",
        message: "Pin sent to your email",
      })
    } else {
      showSnackbar({
        type: "error",
        message: res.message,
      })
    }
  }

  async function handleLogin() {
    const res = await login({ email: "arthadede@gmail.com", pin: "000000"});
    if (res.success) {
      showSnackbar({
        type: "info",
        message: "Login success",
      })
    } else {
      showSnackbar({
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
        <Button onClick={handleSendPin} variant="primary" size="sm">
          Send Pin
        </Button>
        <Button onClick={handleLogin} variant="primary" size="sm">
          Login
        </Button>
        <Button onClick={handleLogout} variant="primary" size="sm">
          Logout
        </Button>
        <Button onClick={handleRefreshToken} variant="primary" size="sm">
          Refresh Token - Client
        </Button>
        <Button onClick={handleUpdateProfile} variant="primary" size="sm">
          Update Profile - Client
        </Button>
      </div>
    </>
  )
}
