'use client'

import Button from "@/components/actions/button";
import { login, logout, sendPin } from "@/services/auth";
import useNotification from "@/hooks/useNotification";

export default function AuthAction() {
  const { showNotification } = useNotification();

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

  return (
    <div className="flex flex-col gap-4">
      <Button onClick={handleSendPin} xVariant="solid" xSize="sm">
        Send Pin
      </Button>
      <Button onClick={handleLogin} xVariant="solid" xSize="sm">
        Login
      </Button>
      <Button onClick={handleLogout} xVariant="solid" xSize="sm">
        Logout
      </Button>
    </div>
  )
}
