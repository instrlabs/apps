import "./globals.css";

import { ReactNode } from "react";
import { Geist } from "next/font/google";

const geist = Geist({
  subsets: ["latin"],
});

export default function RootLayout({
  children,
}: Readonly<{
  children: ReactNode;
}>) {
  return (
    <html lang="en" className={geist.className}>
      <body className="font-sans">
        {children}
      </body>
    </html>
  );
}
