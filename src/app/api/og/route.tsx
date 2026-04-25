import { ImageResponse } from "@vercel/og";
import { NextRequest } from "next/server";

export const runtime = "edge";

export async function GET(req: NextRequest) {
  const { searchParams } = new URL(req.url);
  const title = searchParams.get("title") || "MCPGuard";
  const subtitle =
    searchParams.get("subtitle") ||
    "Security Scanner for MCP Servers";

  return new ImageResponse(
    (
      <div
        style={{
          height: "100%",
          width: "100%",
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
          justifyContent: "center",
          backgroundColor: "#0A0E14",
          position: "relative",
        }}
      >
        {/* Background glow */}
        <div
          style={{
            position: "absolute",
            top: "-20%",
            left: "30%",
            width: "500px",
            height: "500px",
            borderRadius: "50%",
            background:
              "radial-gradient(circle, rgba(0,229,160,0.08) 0%, transparent 70%)",
            display: "flex",
          }}
        />

        {/* Shield icon */}
        <div
          style={{
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            width: "72px",
            height: "72px",
            borderRadius: "16px",
            backgroundColor: "rgba(0,229,160,0.1)",
            border: "2px solid rgba(0,229,160,0.2)",
            marginBottom: "32px",
          }}
        >
          <svg
            width="36"
            height="36"
            viewBox="0 0 24 24"
            fill="none"
            stroke="#00E5A0"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
          >
            <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z" />
          </svg>
        </div>

        {/* Title */}
        <div
          style={{
            display: "flex",
            fontSize: "64px",
            fontWeight: 700,
            color: "#E8ECF0",
            letterSpacing: "-0.02em",
            marginBottom: "12px",
            fontFamily: "monospace",
          }}
        >
          {title}
        </div>

        {/* Subtitle */}
        <div
          style={{
            display: "flex",
            fontSize: "28px",
            color: "#8B99A8",
            fontFamily: "serif",
            fontStyle: "italic",
          }}
        >
          {subtitle}
        </div>

        {/* Terminal preview line */}
        <div
          style={{
            display: "flex",
            alignItems: "center",
            gap: "8px",
            marginTop: "40px",
            padding: "12px 24px",
            borderRadius: "8px",
            backgroundColor: "rgba(26,32,40,0.8)",
            border: "1px solid rgba(42,52,66,0.8)",
            fontFamily: "monospace",
            fontSize: "18px",
          }}
        >
          <span style={{ color: "#5A6878" }}>$</span>
          <span style={{ color: "#E8ECF0" }}>npx mcpguard scan</span>
          <span style={{ color: "#FF4D4D", marginLeft: "16px", fontWeight: 700 }}>
            ✗ RUG-PULL DETECTED
          </span>
        </div>

        {/* Bottom bar */}
        <div
          style={{
            position: "absolute",
            bottom: "32px",
            display: "flex",
            alignItems: "center",
            gap: "8px",
            fontSize: "16px",
            color: "#5A6878",
            fontFamily: "monospace",
          }}
        >
          <span>by AXIOM Collective</span>
          <span>·</span>
          <span style={{ color: "#00E5A0" }}>withaxiom.co</span>
        </div>
      </div>
    ),
    {
      width: 1200,
      height: 630,
    }
  );
}
