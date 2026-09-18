"use client";

import { Button } from "@/components/ui/button";

export default function ErrorPage({ reset }: { error: Error & { digest?: string }; reset: () => void }) {
  return <main className="flex min-h-dvh flex-col items-center justify-center gap-4 px-6 text-center"><h1 className="text-2xl font-semibold">We couldn’t load your workspace</h1><p className="max-w-md text-muted-foreground">There was a connection problem. Please try again in a moment.</p><Button onClick={reset}>Try again</Button></main>;
}
