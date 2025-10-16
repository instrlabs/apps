"use client";

import BrandLink from "@/components/navigation/brand-link";
import Breadcrumbs from "@/components/navigation/breadcrumbs";
import IconButton from "@/components/actions/icon-button";

import { useOverlay } from "@/hooks/useOverlay";
import Avatar from "@/components/avatar";
import { useProfile } from "@/hooks/useProfile";
import NotificationIcon from "@/components/icons/notification-icon";
import Button from "@/components/actions/button";
import { redirect } from "next/navigation";

export default function OverlayTop() {
  const { isLoggedIn, profile } = useProfile();
  const { rightKey, openRight, closeRight } = useOverlay();

  function handleToggleNotifications() {
    if (rightKey === "right:notifications") {
      closeRight();
    } else {
      openRight("notifications");
    }
  }

  function handleToggleProfile() {
    if (rightKey === "right:profile") {
      closeRight();
    } else {
      openRight("profile");
    }
  }

  return (
    <div className="relative w-full flex flex-row justify-between items-center bg-background/80 gap-2 p-2">
      <div className="flex items-center gap-2">
        <BrandLink />
      </div>

      <div className="pointer-events-none absolute inset-x-0 top-1/2 -translate-y-1/2 flex justify-center">
        <Breadcrumbs />
      </div>

      <div className="flex items-center gap-2">
        {isLoggedIn ? (
          <>
            <IconButton
              aria-label="Notifications"
              xColor="secondary"
              onClick={handleToggleNotifications}
            >
              <NotificationIcon className="size-6"/>
            </IconButton>
            <Avatar
              xSize="sm"
              name={profile?.username || "Guest"}
              onClick={handleToggleProfile}
            />
          </>
        ) : (
          <Button onClick={() => redirect("/login")}>
            Login
          </Button>
        )}
      </div>
    </div>
  );
}
