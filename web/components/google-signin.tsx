import Script from "next/script";

const GoogleSignInButton = () => {
  return (
    <>
      <div
        id="g_id_onload"
        data-client_id="773835138675-q3ge7t0s64enkmoeqt0rfqidnm41eg6s.apps.googleusercontent.com"
        data-context="signin"
        data-ux_mode="popup"
        data-login_uri="http://localhost:3000/apps"
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
      <Script src="https://accounts.google.com/gsi/client" async />
    </>
  );
};

export default GoogleSignInButton;