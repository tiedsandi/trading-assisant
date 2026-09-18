import type { Metadata } from "next";
import { redirect } from "next/navigation";
import { Check, Compass } from "lucide-react";
import { Card, CardContent } from "@/components/ui/card";
import { Brand } from "@/features/auth/auth-shell";
import { getCurrentUser } from "@/features/auth/server";
import { LogoutButton } from "@/features/auth/logout-button";

export const metadata: Metadata = { title: "Your workspace" };

export default async function DashboardPage() {
  const user = await getCurrentUser();
  if (!user) redirect("/login");
  return (
    <div className="min-h-dvh">
      <header className="border-b bg-card"><div className="mx-auto flex max-w-6xl items-start justify-between gap-4 px-5 py-5 sm:items-center sm:px-10"><Brand /><LogoutButton /></div></header>
      <main className="mx-auto max-w-6xl px-5 py-12 sm:px-10 sm:py-20">
        <p className="mb-4 text-xs font-semibold tracking-[0.16em] text-primary">YOUR WORKSPACE</p>
        <h1 className="text-3xl font-semibold tracking-tight sm:text-4xl">Welcome to Trading Assistant.</h1>
        <p className="mt-4 max-w-2xl text-base leading-7 text-muted-foreground">You’re signed in and ready for a fresh start.</p>
        <Card className="mt-10 max-w-2xl shadow-none"><CardContent className="flex flex-col gap-6 sm:flex-row sm:items-start">
          <span className="flex size-14 shrink-0 items-center justify-center rounded-2xl bg-secondary text-primary"><Compass className="size-7" aria-hidden="true" /></span>
          <div className="min-w-0">
            <div className="mb-3 inline-flex items-center gap-1.5 rounded-full bg-secondary px-2.5 py-1 text-xs font-medium text-primary"><Check className="size-3" aria-hidden="true" />Signed in</div>
            <h2 className="text-lg font-semibold">Your account is ready</h2>
            <p className="mt-2 break-all text-sm text-muted-foreground">{user.email}</p>
            <p className="mt-5 text-sm leading-6 text-muted-foreground">This is your personal workspace. There’s nothing to set up here yet.</p>
          </div>
        </CardContent></Card>
      </main>
    </div>
  );
}
