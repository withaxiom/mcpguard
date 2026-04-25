import { Hero } from "@/components/Hero";
import { TerminalDemo } from "@/components/TerminalDemo";
import { Features } from "@/components/Features";
import { WaitlistForm } from "@/components/WaitlistForm";
import { Footer } from "@/components/Footer";
import { Nav } from "@/components/Nav";

export default function Home() {
  return (
    <main className="relative">
      <div className="absolute inset-0 overflow-hidden pointer-events-none">
        <div className="absolute top-0 left-1/2 -translate-x-1/2 w-[800px] h-[600px] bg-accent/[0.04] rounded-full blur-[120px]" />
        <div className="absolute bottom-1/3 right-0 w-[400px] h-[400px] bg-accent/[0.03] rounded-full blur-[100px]" />
      </div>

      <Nav />
      <Hero />
      <TerminalDemo />
      <Features />
      <WaitlistForm />
      <Footer />
    </main>
  );
}
