"use client";

import React, { createContext, useContext, useState, ReactNode, useEffect } from "react";
import InfoIcon from "./icons/InfoIcon";
import SuccessIcon from "./icons/SuccessIcon";
import ErrorIcon from "./icons/ErrorIcon";
import WarningIcon from "./icons/WarningIcon";

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

export const useNotification = (): NotificationContextProps => {
  const context = useContext(NotificationContext);

  if (context === undefined) {
    throw new Error("useNotification must be used within a NotificationProvider");
  }

  return context;
};

export const Notification: React.FC = () => {
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

  const getColorStyles = () => {
    // Use theme primary for all notification backgrounds to align with app theming
    return "bg-primary text-primary-foreground";
  };

  return (
    <div className="absolute bottom-4 left-1/2 -translate-x-1/2 z-50 w-full flex justify-center px-4">
      <div className={`w-sm ${animationClass}`}>
        <div className={`px-5 py-4 flex flex-row items-center gap-3 rounded-xl shadow-primary ${getColorStyles()}`}>
          {type === "error" && <ErrorIcon className="shrink-0" />}
          {type === "warning" && <WarningIcon className="shrink-0" />}
          {type === "info" && <InfoIcon className="shrink-0" />}
          <div>{message}</div>
        </div>
      </div>
    </div>
  );
};
