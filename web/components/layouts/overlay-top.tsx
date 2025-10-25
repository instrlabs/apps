"use client";

import Breadcrumbs from "@/components/breadcrumbs";
import IconButton from "@/components/icon-button";

import { useOverlay } from "@/hooks/useOverlay";
import Avatar from "@/components/avatar";
import { useProfile } from "@/hooks/useProfile";
import Icon from "@/components/icon";
import Button from "@/components/button";
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
    <>
      {/* ACTIVE - MOBILE SECTION */}
      <div className="flex flex-col items-center md:hidden p-2">
        <Breadcrumbs />
      </div>

      <div className="relative w-full flex flex-row items-center justify-between gap-2 p-2 bg-background/80">
        <div className="flex items-center gap-2">
          <Icon name="logo" className="size-6" />
        </div>

        {/* ACTIVE - WEB SECTION */}
        <div className="absolute inset-x-0 top-1/2 -translate-y-1/2 hidden md:flex md:justify-center pointer-events-none">
          <Breadcrumbs />
        </div>

        <div className="flex items-center gap-2">
          {isLoggedIn ? (
            <>
              <IconButton
                aria-label="Notifications"
                xColor="transparent"
                onClick={handleToggleNotifications}
              >
                <Icon name="rectangle" className="size-6"/>
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
    </>
  );
}
