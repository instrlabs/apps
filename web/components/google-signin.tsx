"use client";

import Script from "next/script";
import { useEffect, useState } from "react";

const GoogleSignInButton = () => {
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);
  
  return (
    <>
      {mounted && (
        <>
          <div
            id="g_id_onload"
            data-client_id={process.env.NEXT_PUBLIC_GOOGLE_CLIENT_ID}
            data-context="signin"
            data-ux_mode="redirect"
            data-login_uri={`${process.env.NEXT_PUBLIC_API_URL}/auth/google`}
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