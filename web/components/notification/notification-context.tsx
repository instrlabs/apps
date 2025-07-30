"use client";

import React, { createContext, useContext, useState, ReactNode } from "react";

type NotificationType = "success" | "error" | "warning" | "info";

interface NotificationContextProps {
  showNotification: (message: string, type: NotificationType, duration?: number) => void;
  hideNotification: () => void;
  message: string;
  type: NotificationType;
  visible: boolean;
}

const NotificationContext = createContext<NotificationContextProps | undefined>(undefined);

interface NotificationProviderProps {
  children: ReactNode;
}

export const NotificationProvider: React.FC<NotificationProviderProps> = ({ children }) => {
  const [visible, setVisible] = useState(false);
  const [message, setMessage] = useState("");
  const [type, setType] = useState<NotificationType>("info");
  const [timeoutId, setTimeoutId] = useState<NodeJS.Timeout | null>(null);

  const showNotification = (
    message: string,
    type: NotificationType = "info",
    duration: number = 3000
  ) => {
    // Clear any existing timeout
    if (timeoutId) {
      clearTimeout(timeoutId);
    }

    // Set notifications data
    setMessage(message);
    setType(type);
    setVisible(true);

    // Auto-hide after duration
    const id = setTimeout(() => {
      hideNotification();
    }, duration);

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