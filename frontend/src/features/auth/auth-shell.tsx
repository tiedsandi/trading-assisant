import Link from "next/link";
import { ArrowUpRight, Compass } from "lucide-react";
import { Card, CardContent, CardDescription, CardHeader } from "@/components/ui/card";

export function Brand({ light = false }: { light?: boolean }) {
  return (
    <span className="inline-flex items-center gap-3 text-sm font-semibold tracking-tight">
      <span className={`flex size-9 items-center justify-center rounded-xl ${light ? "bg-white/10 text-[#d8edb1]" : "bg-primary text-primary-foreground"}`}>
        <Compass className="size-5" aria-hidden="true" />
      </span>
      Trading Assistant
    </span>
  );
}

export function AuthShell({ mode, children }: { mode: "login" | "register"; children: React.ReactNode }) {
  const register = mode === "register";
  return (
    <main className="grid min-h-dvh lg:grid-cols-[0.9fr_1.1fr]">
      <aside className="relative hidden flex-col justify-between overflow-hidden bg-[#1c3529] p-12 text-white lg:flex xl:p-16" aria-label="Trading Assistant">
        <Brand light />
        <div className="relative my-16 max-w-md">
          <div className="mb-10 flex size-16 items-center justify-center rounded-2xl border border-white/20 text-[#d8edb1]"><ArrowUpRight className="size-8" aria-hidden="true" /></div>
          <p className="mb-5 text-xs font-medium tracking-[0.2em] text-[#c0d1bf]">YOUR TRADING WORKSPACE</p>
          <h2 className="text-5xl leading-[1.12] font-medium tracking-tight xl:text-6xl">A clear place<br />to begin.</h2>
          <p className="mt-6 max-w-xs text-base leading-7 text-[#c0d1bf]">Your own account. Your own space.<br />One step closer to a more focused routine.</p>
        </div>
        <p className="text-xs text-[#c0d1bf]">Built for your next chapter.</p>
        <div className="pointer-events-none absolute -right-36 -bottom-44 size-[430px] rounded-full border border-white/10" aria-hidden="true" />
        <div className="pointer-events-none absolute -right-24 -bottom-32 size-[335px] rounded-full border border-white/10" aria-hidden="true" />
      </aside>
      <section className="flex flex-col px-5 py-7 sm:px-10 lg:px-12">
        <div className="mb-12 lg:hidden"><Brand /></div>
        <div className="flex flex-1 items-center justify-center py-6">
          <div className="w-full max-w-[420px]">
            <Card className="gap-7 border-0 bg-transparent py-0 shadow-none">
              <CardHeader className="gap-3 px-0">
                <p className="text-xs font-semibold tracking-[0.16em] text-primary">{register ? "GET STARTED" : "WELCOME BACK"}</p>
                <h1 className="text-[2rem] leading-tight font-semibold tracking-tight sm:text-4xl">{register ? "Create your account" : "Sign in to your space"}</h1>
                <CardDescription className="text-base leading-6">{register ? "A fresh start for your trading routine." : "Good to see you. Pick up where you left off."}</CardDescription>
              </CardHeader>
              <CardContent className="px-0">{children}</CardContent>
            </Card>
            <p className="mt-8 text-center text-sm text-muted-foreground">
              {register ? "Already have an account?" : "New to Trading Assistant?"}{" "}
              <Link className="rounded-sm font-semibold text-primary underline-offset-4 hover:underline focus-visible:outline-2" href={register ? "/login" : "/register"}>{register ? "Sign in" : "Create an account"}</Link>
            </p>
          </div>
        </div>
        <p className="mt-8 text-center text-xs text-muted-foreground">Trading Assistant · Your personal workspace</p>
      </section>
    </main>
  );
}
