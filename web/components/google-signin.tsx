import Button from "@/components/button";
import { loginByGoogle } from "@/services/auth";
import Icon from "@/components/icon";

function GoogleSignInButton() {
  return (
      <Button onClick={loginByGoogle} color="secondary" size="lg">
        <div className="flex items-center justify-center gap-2">
          <Icon name="google" size={20} />
          Continue with Google
        </div>
      </Button>
  );
};

export default GoogleSignInButton;
