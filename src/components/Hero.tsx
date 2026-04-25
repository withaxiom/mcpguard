"use client";

import { motion } from "framer-motion";

export function Hero() {
  return (
    <section className="relative pt-32 pb-16 px-6">
      <div className="max-w-4xl mx-auto text-center">
        {/* Badge */}
        <motion.div
          initial={{ opacity: 0, y: 12 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5 }}
          className="inline-flex items-center gap-2 px-3 py-1 rounded-full border border-border bg-surface/50 text-xs font-mono text-text-muted mb-8"
        >
          <span className="w-1.5 h-1.5 rounded-full bg-accent animate-pulse" />
          Open Source · CLI · v0.1
        </motion.div>

        {/* Headline */}
        <motion.h1
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6, delay: 0.1 }}
          className="font-serif text-5xl sm:text-6xl md:text-7xl leading-[1.05] tracking-tight mb-6"
        >
          Your MCP servers{" "}
          <br className="hidden sm:block" />
          <span className="italic text-accent">aren&apos;t safe.</span>
        </motion.h1>

        {/* Subheadline */}
        <motion.p
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6, delay: 0.2 }}
          className="text-lg sm:text-xl text-text-muted max-w-2xl mx-auto mb-10 leading-relaxed"
        >
          MCPGuard scans for rug-pulls, prompt injection, and tool poisoning
          before they compromise your AI agents.{" "}
          <span className="text-text">One command. Zero trust.</span>
        </motion.p>

        {/* Install command */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6, delay: 0.3 }}
          className="flex flex-col sm:flex-row items-center justify-center gap-4"
        >
          <div className="group relative">
            <div className="absolute -inset-[1px] rounded-lg bg-gradient-to-r from-accent/20 via-accent/5 to-accent/20 opacity-0 group-hover:opacity-100 transition-opacity blur-sm" />
            <code className="relative flex items-center gap-3 px-5 py-3 rounded-lg bg-surface border border-border font-mono text-sm">
              <span className="text-text-dim select-none">$</span>
              <span className="text-text">npx mcpguard scan</span>
              <button
                onClick={() => {
                  navigator.clipboard.writeText("npx mcpguard scan");
                }}
                className="text-text-dim hover:text-accent transition-colors ml-2"
                title="Copy to clipboard"
              >
                <svg
                  width="14"
                  height="14"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                >
                  <rect x="9" y="9" width="13" height="13" rx="2" ry="2" />
                  <path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1" />
                </svg>
              </button>
            </code>
          </div>

          <a
            href="#waitlist"
            className="px-6 py-3 rounded-lg bg-accent text-void font-semibold text-sm hover:brightness-110 transition-all shadow-lg shadow-accent/20"
          >
            Join the Waitlist →
          </a>
        </motion.div>

        {/* Social proof hint */}
        <motion.p
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 0.6, delay: 0.6 }}
          className="mt-8 text-xs text-text-dim font-mono"
        >
          Built by{" "}
          <a
            href="https://withaxiom.co"
            target="_blank"
            rel="noopener noreferrer"
            className="text-text-muted hover:text-accent transition-colors"
          >
            AXIOM Collective
          </a>
        </motion.p>
      </div>
    </section>
  );
}
