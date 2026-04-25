export function Footer() {
  return (
    <footer className="border-t border-border/50 py-12 px-6">
      <div className="max-w-5xl mx-auto flex flex-col sm:flex-row items-center justify-between gap-6">
        <div className="flex items-center gap-2">
          <div className="w-5 h-5 rounded bg-accent/10 border border-accent/20 flex items-center justify-center">
            <svg
              width="10"
              height="10"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="2.5"
              strokeLinecap="round"
              strokeLinejoin="round"
              className="text-accent"
            >
              <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z" />
            </svg>
          </div>
          <span className="font-mono text-xs text-text-muted">
            MCPGuard
          </span>
        </div>

        <div className="flex items-center gap-6 text-xs text-text-dim">
          <a
            href="https://github.com/withaxiom/mcpguard"
            target="_blank"
            rel="noopener noreferrer"
            className="hover:text-text-muted transition-colors"
          >
            GitHub
          </a>
          <a
            href="https://x.com/withaxiom"
            target="_blank"
            rel="noopener noreferrer"
            className="hover:text-text-muted transition-colors"
          >
            @withaxiom
          </a>
          <span>
            © {new Date().getFullYear()}{" "}
            <a
              href="https://withaxiom.co"
              target="_blank"
              rel="noopener noreferrer"
              className="hover:text-text-muted transition-colors"
            >
              AXIOM Collective
            </a>
          </span>
        </div>
      </div>
    </footer>
  );
}
