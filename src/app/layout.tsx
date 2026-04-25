import type { Metadata } from "next";
import "./globals.css";
import { PostHogProvider } from "@/components/PostHogProvider";

export const metadata: Metadata = {
  title: "MCPGuard — Security Scanner for MCP Servers",
  description:
    "Detect rug-pulls, prompt injection, and tool poisoning in MCP servers before they compromise your AI agents. Open-source CLI security scanner.",
  metadataBase: new URL(
    process.env.NEXT_PUBLIC_SITE_URL || "https://mcpguard.com"
  ),
  openGraph: {
    title: "MCPGuard — Security Scanner for MCP Servers",
    description:
      "Detect rug-pulls, prompt injection, and tool poisoning before they compromise your AI agents.",
    url: "/",
    siteName: "MCPGuard",
    images: [{ url: "/api/og", width: 1200, height: 630 }],
    type: "website",
  },
  twitter: {
    card: "summary_large_image",
    title: "MCPGuard — Security Scanner for MCP Servers",
    description:
      "Detect rug-pulls, prompt injection, and tool poisoning before they compromise your AI agents.",
    images: ["/api/og"],
  },
  robots: { index: true, follow: true },
  icons: {
    icon: [
      { url: "/favicon.svg", type: "image/svg+xml" },
    ],
  },
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" className="noise-overlay">
      <head>
        <link rel="preconnect" href="https://fonts.googleapis.com" />
        <link
          rel="preconnect"
          href="https://fonts.gstatic.com"
          crossOrigin="anonymous"
        />
        <link
          href="https://fonts.googleapis.com/css2?family=Instrument+Serif:ital@0;1&family=JetBrains+Mono:wght@400;500;600;700&family=Inter:wght@400;500;600;700&display=swap"
          rel="stylesheet"
        />
      </head>
      <body className="min-h-screen bg-void text-text antialiased">
        <PostHogProvider>{children}</PostHogProvider>
      </body>
    </html>
  );
}
