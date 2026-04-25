"use client";

import { useState, FormEvent } from "react";
import { motion } from "framer-motion";

// Mirrors the regex used by the /api/waitlist route so client-side
// feedback matches what the server will accept.
const EMAIL_RE = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

export function WaitlistForm() {
  const [email, setEmail] = useState("");
  const [touched, setTouched] = useState(false);
  const [status, setStatus] = useState<
    "idle" | "loading" | "success" | "error"
  >("idle");
  const [errorMessage, setErrorMessage] = useState("");

  const trimmed = email.trim();
  const inlineError =
    touched && trimmed.length > 0 && !EMAIL_RE.test(trimmed)
      ? "Please enter a valid email address."
      : "";
  const canSubmit = EMAIL_RE.test(trimmed) && status !== "loading";

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setTouched(true);
    if (!canSubmit) return;

    setStatus("loading");
    setErrorMessage("");

    try {
      const res = await fetch("/api/waitlist", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email: trimmed }),
      });

      const data = await res.json();

      if (!res.ok) {
        throw new Error(data.error || "Something went wrong");
      }

      setStatus("success");
      setEmail("");

      // Track conversion
      if (typeof window !== "undefined" && (window as unknown as { posthog?: { capture: (event: string, properties?: Record<string, string>) => void } }).posthog) {
        (window as unknown as { posthog: { capture: (event: string, properties?: Record<string, string>) => void } }).posthog.capture("waitlist_signup", { email });
      }
    } catch (err) {
      setStatus("error");
      setErrorMessage(
        err instanceof Error ? err.message : "Something went wrong"
      );
    }
  };

  return (
    <section id="waitlist" className="relative py-24 px-6">
      <div className="max-w-2xl mx-auto text-center">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ duration: 0.5 }}
        >
          <h2 className="font-serif text-3xl sm:text-4xl tracking-tight mb-4">
            Ship AI agents{" "}
            <span className="italic text-accent">with confidence</span>
          </h2>
          <p className="text-text-muted mb-8 max-w-md mx-auto">
            Get notified when MCPGuard launches. Early access members get the
            managed scanning dashboard free for 6 months.
          </p>

          {status === "success" ? (
            <motion.div
              initial={{ opacity: 0, scale: 0.95 }}
              animate={{ opacity: 1, scale: 1 }}
              className="inline-flex items-center gap-2 px-6 py-3 rounded-lg bg-accent/10 border border-accent/20 text-accent font-mono text-sm"
            >
              <svg
                width="16"
                height="16"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2.5"
                strokeLinecap="round"
                strokeLinejoin="round"
              >
                <path d="M20 6L9 17l-5-5" />
              </svg>
              You&apos;re on the list. We&apos;ll be in touch.
            </motion.div>
          ) : (
            <form
              onSubmit={handleSubmit}
              className="flex flex-col sm:flex-row gap-3 max-w-md mx-auto"
            >
              <div className="relative flex-1">
                <input
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  onBlur={() => setTouched(true)}
                  placeholder="you@company.com"
                  required
                  aria-invalid={inlineError ? true : undefined}
                  aria-describedby={inlineError ? "waitlist-email-error" : undefined}
                  className={`w-full px-4 py-3 rounded-lg bg-surface border text-text placeholder:text-text-dim font-mono text-sm focus:outline-none focus:ring-1 transition-colors ${
                    inlineError
                      ? "border-danger/60 focus:border-danger focus:ring-danger/20"
                      : "border-border focus:border-accent/50 focus:ring-accent/20"
                  }`}
                />
              </div>
              <button
                type="submit"
                disabled={!canSubmit}
                className="px-6 py-3 rounded-lg bg-accent text-void font-semibold text-sm hover:brightness-110 transition-all shadow-lg shadow-accent/20 disabled:opacity-50 disabled:cursor-not-allowed whitespace-nowrap"
              >
                {status === "loading" ? (
                  <span className="flex items-center gap-2">
                    <svg
                      className="animate-spin h-4 w-4"
                      viewBox="0 0 24 24"
                    >
                      <circle
                        className="opacity-25"
                        cx="12"
                        cy="12"
                        r="10"
                        stroke="currentColor"
                        strokeWidth="4"
                        fill="none"
                      />
                      <path
                        className="opacity-75"
                        fill="currentColor"
                        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"
                      />
                    </svg>
                    Joining...
                  </span>
                ) : (
                  "Join Waitlist"
                )}
              </button>
            </form>
          )}

          {inlineError && status !== "error" && (
            <p
              id="waitlist-email-error"
              className="mt-3 text-sm text-danger font-mono"
            >
              {inlineError}
            </p>
          )}

          {status === "error" && (
            <p className="mt-3 text-sm text-danger font-mono">{errorMessage}</p>
          )}

          <p className="mt-6 text-xs text-text-dim">
            No spam. Unsubscribe anytime. We respect your inbox.
          </p>
        </motion.div>
      </div>
    </section>
  );
}
