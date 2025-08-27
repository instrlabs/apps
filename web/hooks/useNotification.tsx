"use client";

import React, {
  useState,
  useEffect,
  useContext,
  createContext,
  ReactNode,
} from "react";
import clsx from "clsx";

import InfoIcon from "@/components/icons/InfoIcon";
import ErrorIcon from "@/components/icons/ErrorIcon";
import WarningIcon from "@/components/icons/WarningIcon";

type NotificationType = "error" | "warning" | "info";

interface NotificationContextProps {
  showNotification: (message: string, type: NotificationType, duration?: number) => void;
  hideNotification: () => void;
  message: string;
  type: NotificationType;
  visible: boolean;
}

const NotificationContext = createContext<NotificationContextProps | undefined>(undefined);

export const NotificationProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [visible, setVisible] = useState(false);
  const [message, setMessage] = useState("");
  const [type, setType] = useState<NotificationType>("info");
  const [timeoutId, setTimeoutId] = useState<NodeJS.Timeout | null>(null);

  const showNotification = (
    message: string,
    type: NotificationType = "info",
    duration: number = 3000
  ) => {
    if (timeoutId) clearTimeout(timeoutId);

    setMessage(message);
    setType(type);
    setVisible(true);

    const id = setTimeout(() => hideNotification(), duration);

    setTimeoutId(id);
  };

  const hideNotification = () => {
    setVisible(false);
  };

  return (
    <NotificationContext.Provider
      value={{
        showNotification,
        hideNotification,
        message,
        type,
        visible,
      }}
    >
      {children}
    </NotificationContext.Provider>
  );
};

const useNotification = (): NotificationContextProps => {
  const context = useContext(NotificationContext);

  if (context === undefined) {
    throw new Error("useNotification must be used within a NotificationProvider");
  }

  return context;
};

export default useNotification;

export const NotificationWidget: React.FC = () => {
  const { message, type, visible } = useNotification();
  const [isRendered, setIsRendered] = useState(false);
  const [animationClass, setAnimationClass] = useState("");

  useEffect(() => {
    if (visible) {
      setIsRendered(true);
      setAnimationClass("animate-slide-in");
    } else if (isRendered) {
      setAnimationClass("animate-slide-out");
      const timer = setTimeout(() => setIsRendered(false), 300);
      return () => clearTimeout(timer);
    }
  }, [visible, isRendered]);

  if (!isRendered) return null;

  return (
    <div className="absolute bottom-4 left-1/2 -translate-x-1/2 z-50 w-full flex justify-center px-4">
      <div className={`w-sm ${animationClass}`}>
        <div className={clsx(
          "px-5 py-4",
          "flex flex-row items-center gap-3",
          "rounded-xl shadow-primary bg-white"
        )}>
          {type === "error" && <ErrorIcon className="shrink-0" />}
          {type === "warning" && <WarningIcon className="shrink-0" />}
          {type === "info" && <InfoIcon className="shrink-0" />}
          <div>{message}</div>
        </div>
      </div>
    </div>
  );
};
