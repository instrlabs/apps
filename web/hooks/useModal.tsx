"use client";

import React, { createContext, useContext, useEffect, useState, ReactNode } from "react";

type ModalContextProps = {
  openModal: (content: ReactNode) => void;
  closeModal: () => void;
  content: ReactNode | null;
  visible: boolean;
};

const ModalContext = createContext<ModalContextProps | undefined>(undefined);

export const ModalProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [visible, setVisible] = useState(false);
  const [content, setContent] = useState<ReactNode | null>(null);
  const [isRendered, setIsRendered] = useState(false);
  const [animationClass, setAnimationClass] = useState("");

  const openModal = (node: ReactNode) => {
    setContent(node);
    setVisible(true);
  };

  const closeModal = () => {
    setVisible(false);
  };

  useEffect(() => {
    if (visible) {
      setIsRendered(true);
      setAnimationClass("animate-fade-in");
    } else if (isRendered) {
      setAnimationClass("animate-fade-out");
      const timer = setTimeout(() => {
        setIsRendered(false);
        setContent(null);
      }, 200);
      return () => clearTimeout(timer);
    }
  }, [visible, isRendered]);

  return (
    <ModalContext.Provider value={{ openModal, closeModal, content, visible }}>
      {children}
      <ModalWidgetInternal
        isRendered={isRendered}
        animationClass={animationClass}
        closeModal={closeModal}
        content={content}
      />
    </ModalContext.Provider>
  );
};

const useModal = (): ModalContextProps => {
  const ctx = useContext(ModalContext);
  if (!ctx) throw new Error("useModal must be used within a ModalProvider");
  return ctx;
};

export default useModal;

// Internal widget used by provider so consumers don't have to place a widget manually
const ModalWidgetInternal: React.FC<{
  isRendered: boolean;
  animationClass: string;
  closeModal: () => void;
  content: ReactNode | null;
}> = ({ isRendered, animationClass, closeModal, content }) => {
  if (!isRendered || !content) return null;
  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      <div className="absolute inset-0 bg-black/50" onClick={closeModal} />
      <div className={animationClass}>{content}</div>
    </div>
  );
};

// Optional external widget if a project prefers manual placement (mirrors NotificationWidget style)
export const ModalWidget: React.FC = () => {
  const { content, visible, closeModal } = useModal();
  const [isRendered, setIsRendered] = useState(false);
  const [animationClass, setAnimationClass] = useState("");

  useEffect(() => {
    if (visible) {
      setIsRendered(true);
      setAnimationClass("animate-fade-in");
    } else if (isRendered) {
      setAnimationClass("animate-fade-out");
      const timer = setTimeout(() => setIsRendered(false), 200);
      return () => clearTimeout(timer);
    }
  }, [visible, isRendered]);

  if (!isRendered || !content) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      <div className="absolute inset-0 bg-black/50" onClick={closeModal} />
      <div className={animationClass}>{content}</div>
    </div>
  );
};
