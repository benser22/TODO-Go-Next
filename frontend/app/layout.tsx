import type { Metadata } from "next";
import { Fraunces, IBM_Plex_Sans } from "next/font/google";
import "./globals.css";

const display = Fraunces({
  subsets: ["latin"],
  variable: "--f-display",
});

const text = IBM_Plex_Sans({
  subsets: ["latin"],
  weight: ["400", "500"],
  variable: "--f-text",
});

export const metadata: Metadata = {
  title: "Tasks — ledger",
  description: "TODO clásico: API en Go + Postgres, front en Next.js",
};

export default function RootLayout({
  children,
}: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="es">
      <body className={`${display.variable} ${text.variable}`}>
        {children}
      </body>
    </html>
  );
}
