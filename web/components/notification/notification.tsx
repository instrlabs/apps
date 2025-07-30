"use client";

import React, { useEffect, useState } from "react";
import { useNotification } from "./notification-context";
import clsx from "clsx";
import Image from "next/image";

export const Notification: React.FC = () => {
  const { message, type, visible, hideNotification } = useNotification();
  const [isRendered, setIsRendered] = useState(false);
  const [animationClass, setAnimationClass] = useState("");

  useEffect(() => {
    if (visible) {
      setIsRendered(true);
      setAnimationClass("animate-slide-in");
    } else if (isRendered) {
      setAnimationClass("animate-slide-out");
      const timer = setTimeout(() => {
        setIsRendered(false);
      }, 300); // Match this with the animation duration
      return () => clearTimeout(timer);
    }
  }, [visible, isRendered]);

  if (!isRendered) return null;

  const getTypeStyles = () => {
    switch (type) {
      case "success":
        return "text-green-500";
      case "error":
        return "text-red-500";
      case "warning":
        return "text-amber-500";
      case "info":
      default:
        return "text-blue-500";
    }
  };

  const getBorderStyles = () => {
    switch (type) {
      case "success":
        return "border-green-100";
      case "error":
        return "border-red-100";
      case "warning":
        return "border-amber-100";
      case "info":
      default:
        return "border-blue-100";
    }
  };

  const getIcon = () => {
    switch (type) {
      case "success":
        return (
          <Image src="/notifications/success-icon.svg" alt="success-icon" width={24} height={24} />
        );
      case "error":
        return (
          <Image src="/notifications/error-icon.svg" alt="error-icon" width={24} height={24} />
        );
      case "warning":
        return (
          <Image src="/notifications/warning-icon.svg" alt="warning-icon" width={24} height={24} />
        );
      case "info":
      default:
        return <Image src="/notifications/info-icon.svg" alt="info-icon" width={24} height={24} />;
    }
  };

  const getTitle = () => {
    switch (type) {
      case "success":
        return "Success";
      case "error":
        return "Error";
      case "warning":
        return "Warning";
      case "info":
      default:
        return "Information";
    }
  };

  return (
    <div className="fixed top-0 right-0 z-50 flex justify-end px-4 pt-4">
      <div className={`w-[350px] ${animationClass}`}>
        <div
          className={clsx(
            "flex flex-row",
            "p-2 gap-2 rounded-md bg-white",
            `border ${getBorderStyles()}`
          )}
          role="alert"
        >
          <div className="flex items-start">{getIcon()}</div>
          <div className="flex-1">
            <h4 className={`font-medium ${getTypeStyles()}`}>{getTitle()}</h4>
            <div className="text-sm">{message}</div>
          </div>
          <div className="flex items-center">
            <button
              type="button"
              className={clsx(
                "flex items-center justify-center",
                "rounded-full p-1.5 h-8 w-8 hover:bg-white/20"
              )}
              onClick={hideNotification}
              aria-label="Close"
            >
              <Image src="/notifications/close-icon.svg" alt="close-icon" width={16} height={16} />
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};
