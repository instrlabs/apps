"use client";

import Breadcrumbs from "@/components/breadcrumbs";
import { useOverlay } from "@/hooks/useOverlay";
import Avatar from "@/components/avatar";
import { useProfile } from "@/hooks/useProfile";
import Icon from "@/components/icon";
import Button from "@/components/button";
import { redirect } from "next/navigation";

export default function OverlayTop() {
  const { isLoggedIn, profile } = useProfile();
  const { rightKey, openRight, closeRight } = useOverlay();

  function handleToggleProfile() {
    if (rightKey === "right:profile") {
      closeRight();
    } else {
      openRight("profile");
    }
  }

  return (
    <>
      {/* MOBILE BREADCRUMBS */}
      <div className="flex md:hidden flex-col items-center p-2">
        <Breadcrumbs />
      </div>

      {/* HEADER */}
      <div className="flex w-full items-center justify-between gap-2 bg-black p-2">
        {/* LEFT - LOGO */}
        <div className="flex shrink-0 items-center gap-2">
          <Icon name="logo" size={40} />
        </div>

        {/* CENTER - BREADCRUMBS (WEB ONLY) */}
        <div className="hidden flex-1 items-center justify-center md:flex">
          <Breadcrumbs />
        </div>

        {/* RIGHT - AUTH SECTION */}
        <div className="flex shrink-0 items-center gap-2">
          {isLoggedIn ? (
            <Avatar
              size="sm"
              name={profile?.username || "Guest"}
              onClick={handleToggleProfile}
            />
          ) : (
            <Button onClick={() => redirect("/login")}>Login</Button>
          )}
        </div>
      </div>
    </>
  );
}
