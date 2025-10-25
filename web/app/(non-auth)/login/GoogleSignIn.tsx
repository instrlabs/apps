import Button from "@/components/button";
import Icon from "@/components/icon";

import { loginByGoogle } from "@/services/auth";

function GoogleSignIn() {
  return (
      <Button onClick={loginByGoogle} color="secondary" size="lg">
        <div className="flex items-center justify-center gap-2">
          <Icon name="google" size={20} />
          Continue with Google
        </div>
      </Button>
  );
}

export default GoogleSignIn;
