"use client";

import React, {
  useState,
  useEffect,
  useContext,
  createContext,
  ReactNode,
} from "react";

import InfoIcon from "@/components/icons/InfoIcon";
import ErrorIcon from "@/components/icons/ErrorIcon";
import WarningIcon from "@/components/icons/WarningIcon";

type NotificationType = "error" | "warning" | "info";

type NotificationMessage = {
  title: string;
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
    const { title, message, type = "info", duration = 3000 } = payload;
    setData({ title, message, type, duration });
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
      setShouldRender(false);
    } else if (visible && !shouldRender) {
      setAnimationClass("animate-notification-in");
      setShouldRender(true);
    }
  }, [visible, shouldRender]);

  if (!shouldRender) return null;

  return (
    <div className="absolute top-[100px] right-0 z-50 w-full flex justify-center">
      <div className={`w-sm ${animationClass}`}>
        <div
          className={`
          p-4 flex flex-row items-center gap-3
          rounded-xl shadow-primary bg-white border
          ${data?.type === "error" && "border-red-300"}
          ${data?.type === "warning" && "border-yellow-300"}
          ${data?.type === "info" && "border-blue-300"}
        `}
        >
          {data?.type === "error" && <ErrorIcon size={32} className="shrink-0 text-red-500" />}
          {data?.type === "warning" && <WarningIcon size={32} className="shrink-0 text-yellow-500" />}
          {data?.type === "info" && <InfoIcon size={32} className="shrink-0 text-blue-500" />}
          <div className="flex flex-col">
            <h5 className="text-base font-medium">{data?.title}</h5>
            <p className="text-sm">{data?.message}</p>
          </div>
        </div>
      </div>
    </div>
  );
};
