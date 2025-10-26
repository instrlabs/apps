"use client";

import React, {
  useState,
  useEffect,
  useContext,
  createContext,
  ReactNode,
} from "react";
import SnackbarIcon from "@/components/snackbar-icon";

type NotificationType = "error" | "warning" | "info" | "success";

type Position =
  | "left-top"
  | "right-top"
  | "left-bottom"
  | "right-bottom";

type SnackbarMessage = {
  message: string;
  type?: NotificationType;
  duration?: number;
  position?: Position;
}

interface SnackbarContextProps {
  showSnackbar: (message: SnackbarMessage) => void;
  hideSnackbar: () => void;
  data: SnackbarMessage | null;
  visible: boolean;
}

const SnackbarContext = createContext<SnackbarContextProps | undefined>(undefined);

export const SnackbarProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [visible, setVisible] = useState(false);
  const [data, setData] = useState<SnackbarMessage | null>(null);

  const showSnackbar = (payload: SnackbarMessage) => {
    const {
      message,
      type = "success",
      duration = 3000,
      position = "right-bottom"
    } = payload;
    setData({ message, type, duration, position });
    setVisible(true);
  };

  const hideSnackbar = () => {
    setVisible(false)
  };

  const value = {
    showSnackbar,
    hideSnackbar,
    data,
    visible
  };

  return (
    <SnackbarContext.Provider value={value}>
      {children}
    </SnackbarContext.Provider>
  );
};

const useSnackbar = (): SnackbarContextProps => {
  const context = useContext(SnackbarContext);
  if (context === undefined) throw new Error("useSnackbar must be used within a SnackbarProvider");
  return context;
};

export default useSnackbar;

export const SnackbarWidget: React.FC = () => {
  const { data, visible, hideSnackbar } = useSnackbar();
  const [shouldRender, setShouldRender] = useState(false);
  const [animationClass, setAnimationClass] = useState("");

  useEffect(() => {
    if (visible && shouldRender) {
      const timer = setTimeout(() => {
        setAnimationClass("animate-notification-out");
        hideSnackbar();
      }, data?.duration || 3000);
      return () => clearTimeout(timer);
    } else if (!visible && shouldRender) {
      const timer = setTimeout(() => {
        setShouldRender(false);
      }, 300);
      return () => clearTimeout(timer);
    } else if (visible && !shouldRender) {
      setAnimationClass("animate-notification-in");
      setShouldRender(true);
    }
  }, [visible, shouldRender, data?.duration, hideSnackbar]);

  if (!shouldRender) return null;

  const positionClasses = {
    "left-top": "fixed top-6 left-6",
    "right-top": "fixed top-6 right-6",
    "left-bottom": "fixed bottom-6 left-6",
    "right-bottom": "fixed bottom-6 right-6"
  };

  const getColorClasses = (type?: NotificationType) => {
    switch (type) {
      case "error":
        return {
          bg: "bg-red-500/10",
          border: "border-white/10",
          text: "text-red-500",
          icon: "text-red-500"
        };
      case "warning":
        return {
          bg: "bg-yellow-500/10",
          border: "border-white/10",
          text: "text-yellow-500",
          icon: "text-yellow-500"
        };
      case "info":
        return {
          bg: "bg-blue-500/10",
          border: "border-white/10",
          text: "text-blue-500",
          icon: "text-blue-500"
        };
      case "success":
      default:
        return {
          bg: "bg-green-500/10",
          border: "border-white/10",
          text: "text-green-500",
          icon: "text-green-500"
        };
    }
  };


  const colors = getColorClasses(data?.type);

  return (
    <div className={`${positionClasses[data?.position || "right-bottom"]} z-50`}>
      <div className={animationClass}>
        <div
          className={[
            "flex items-center gap-2.5 max-w-sm p-2 rounded border",
            colors.bg,
            colors.border
          ]
            .filter(Boolean)
            .join(" ")}
        >
          <div className={`shrink-0 ${colors.icon}`}>
            <SnackbarIcon type={data?.type} className="w-5 h-5" />
          </div>
          <p
            className={[
              "flex-1 text-sm leading-5 font-normal",
              colors.text
            ]
              .filter(Boolean)
              .join(" ")}
          >
            {data?.message}
          </p>
        </div>
      </div>
    </div>
  );
};

// Backwards compatibility aliases
export const NotificationProvider = SnackbarProvider;
export const NotificationWidget = SnackbarWidget;