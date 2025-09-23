"use client";

import React, {
  useState,
  useEffect,
  useContext,
  createContext,
  ReactNode,
} from "react";

type NotificationType = "error" | "warning" | "info";

type NotificationMessage = {
  message: string;
  type?: NotificationType;
  duration?: number;
}

interface NotificationContextProps {
  showNotification: (message: NotificationMessage) => void;
  hideNotification: () => void;
  data: NotificationMessage | null;
  visible: boolean;
}

const NotificationContext = createContext<NotificationContextProps | undefined>(undefined);

export const NotificationProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [visible, setVisible] = useState(false);
  const [data, setData] = useState<NotificationMessage | null>(null);

  const showNotification = (payload: NotificationMessage) => {
    const { message, type = "info", duration = 3000 } = payload;
    setData({ message, type, duration });
    setVisible(true);
  };

  const hideNotification = () => {
    setVisible(false)
  };

  const value = {
    showNotification,
    hideNotification,
    data,
    visible
  };

  return (
    <NotificationContext.Provider value={value}>
      {children}
    </NotificationContext.Provider>
  );
};

const useNotification = (): NotificationContextProps => {
  const context = useContext(NotificationContext);
  if (context === undefined) throw new Error("useNotification must be used within a NotificationProvider");
  return context;
};

export default useNotification;

export const NotificationWidget: React.FC = () => {
  const { data, visible, hideNotification } = useNotification();
  const [shouldRender, setShouldRender] = useState(false);
  const [animationClass, setAnimationClass] = useState("");

  useEffect(() => {
    if (visible && shouldRender) {
      const timer = setTimeout(() => {
        setAnimationClass("animate-notification-out");
        hideNotification();
      }, data?.duration || 3000);
      return () => clearTimeout(timer);
    } else if(!visible && shouldRender) {
      const timer = setTimeout(() => {
        setShouldRender(false);
      }, 300);
      return () => clearTimeout(timer);
    } else if (visible && !shouldRender) {
      setAnimationClass("animate-notification-in");
      setShouldRender(true);
    }
  }, [visible, shouldRender]);

  if (!shouldRender) return null;

  return (
    <div className="fixed bottom-6 right-6 z-50">
      <div className={`w-xs ${animationClass}`}>
        <div
          className={`
          p-4 flex flex-col rounded-lg border
          ${data?.type === "error" && "bg-red-500/10 border-red-500"}
          ${data?.type === "warning" && "bg-yellow-500/10 border-yellow-500"}
          ${data?.type === "info" && "bg-blue-500/10 border-blue-500"}
        `}
        >
          <p className={`text-sm
            ${data?.type === "error" ? "text-red-500" : ""}
            ${data?.type === "warning" ? "text-yellow-500" : ""}
            ${data?.type === "info" ? "text-blue-500" : ""}
          `}>{data?.message}</p>
        </div>
      </div>
    </div>
  );
};
