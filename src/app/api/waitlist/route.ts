import { NextRequest, NextResponse } from "next/server";

export async function POST(req: NextRequest) {
  try {
    const { email } = await req.json();

    if (!email || typeof email !== "string") {
      return NextResponse.json(
        { error: "Email is required" },
        { status: 400 }
      );
    }

    // Basic email validation
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(email)) {
      return NextResponse.json(
        { error: "Invalid email format" },
        { status: 400 }
      );
    }

    const normalizedEmail = email.toLowerCase().trim();

    // Store in Supabase
    if (
      process.env.NEXT_PUBLIC_SUPABASE_URL &&
      process.env.SUPABASE_SERVICE_ROLE_KEY
    ) {
      const { createClient } = await import("@supabase/supabase-js");
      const supabase = createClient(
        process.env.NEXT_PUBLIC_SUPABASE_URL,
        process.env.SUPABASE_SERVICE_ROLE_KEY
      );

      const { error: dbError } = await supabase
        .from("waitlist")
        .upsert(
          {
            email: normalizedEmail,
            source: "landing_page",
            created_at: new Date().toISOString(),
          },
          { onConflict: "email" }
        );

      if (dbError) {
        console.error("Supabase error:", dbError);
        // Don't fail — still send the email
      }
    }

    // Send welcome email via Resend
    if (process.env.RESEND_API_KEY) {
      const { Resend } = await import("resend");
      const resend = new Resend(process.env.RESEND_API_KEY);

      await resend.emails.send({
        from: "MCPGuard <noreply@mcpguard.com>",
        to: normalizedEmail,
        subject: "You're on the MCPGuard waitlist 🛡️",
        html: `
          <div style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; max-width: 480px; margin: 0 auto; padding: 40px 20px;">
            <h1 style="font-size: 24px; color: #111; margin-bottom: 16px;">Welcome to MCPGuard</h1>
            <p style="font-size: 16px; color: #555; line-height: 1.6;">
              You've been added to the MCPGuard early access list. We'll notify you as soon as the CLI and managed dashboard are ready.
            </p>
            <p style="font-size: 16px; color: #555; line-height: 1.6;">
              In the meantime, star us on <a href="https://github.com/withaxiom/mcpguard" style="color: #00C488;">GitHub</a> 
              and follow <a href="https://x.com/withaxiom" style="color: #00C488;">@withaxiom</a> for updates.
            </p>
            <p style="font-size: 14px; color: #999; margin-top: 32px;">— The AXIOM Collective team</p>
          </div>
        `,
      });
    }

    return NextResponse.json({ success: true });
  } catch (error) {
    console.error("Waitlist error:", error);
    return NextResponse.json(
      { error: "Internal server error" },
      { status: 500 }
    );
  }
}
