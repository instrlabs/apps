"use client";

import Script from "next/script";
import ROUTES from "@/constants/routes";
import { AUTH_ENDPOINTS } from "@/constants/api";
import { useEffect, useState } from "react";

const GoogleSignInButton = () => {
  const [mounted, setMounted] = useState(false);
  
  useEffect(() => {
    setMounted(true);
  }, []);
  // Get the absolute URL for the frontend callback route
  const callbackUrl = `${window.location.origin}${ROUTES.GOOGLE_CALLBACK}`;
  
  // Only render the Google Sign-In button after the component has mounted
  // to avoid window is not defined errors
  return (
    <>
      {mounted && (
        <>
          <div
            id="g_id_onload"
            data-client_id="773835138675-q3ge7t0s64enkmoeqt0rfqidnm41eg6s.apps.googleusercontent.com"
            data-context="signin"
            data-ux_mode="redirect"
            data-login_uri={`${AUTH_ENDPOINTS.GOOGLE}?redirect_uri=${encodeURIComponent(callbackUrl)}`}
            data-auto_prompt="false"
          ></div>
          <div
            className="g_id_signin"
            data-type="standard"
            data-shape="rectangular"
            data-theme="outline"
            data-text="signin_with"
            data-size="large"
            data-logo_alignment="center"
          ></div>
        </>
      )}
      <Script src="https://accounts.google.com/gsi/client" async />
    </>
  );
};

export default GoogleSignInButton;