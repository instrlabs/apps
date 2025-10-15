"use client";

import BrandLink from "@/components/navigation/brand-link";
import Breadcrumbs from "@/components/navigation/breadcrumbs";
import IconButton from "@/components/actions/icon-button";

import { useOverlay } from "@/hooks/useOverlay";
import Avatar from "@/components/avatar";
import { useProfile } from "@/hooks/useProfile";
import NotificationIcon from "@/components/icons/notification-icon";

export default function OverlayTop() {
  const { profile } = useProfile();
  const { openRight } = useOverlay();

  return (
    <div className="relative w-full flex flex-row justify-between items-center bg-background/80">
      <div className="flex items-center gap-2 p-2">
        <BrandLink />
      </div>

      <div className="pointer-events-none absolute inset-x-0 top-1/2 -translate-y-1/2 flex justify-center">
        <Breadcrumbs />
      </div>

      <div className="flex items-center gap-2 p-2">
        <IconButton
          aria-label="Notifications"
          xColor="secondary"
          onClick={() => openRight("notifications")}
        >
          <NotificationIcon />
        </IconButton>
        <Avatar
          xsize="sm"
          name={profile?.username || "Guest"}
          onClick={() => openRight("profile")}
        />
      </div>
    </div>
  );
}
