import Button from "@/components/button";
import { loginByGoogle } from "@/services/auth";
import GoogleIcon from "@/components/icons/google-icon";

function GoogleSignInButton() {
  return (
      <Button onClick={loginByGoogle} color="secondary" size="lg">
        <div className="flex items-center justify-center gap-2">
          <GoogleIcon className="size-5" />
          Continue with Google
        </div>
      </Button>
  );
};

export default GoogleSignInButton;
