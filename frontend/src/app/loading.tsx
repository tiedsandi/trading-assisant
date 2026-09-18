import { LoaderCircle } from "lucide-react";

export default function Loading() {
  return <main className="flex min-h-dvh items-center justify-center gap-3 text-sm text-muted-foreground" role="status"><LoaderCircle className="size-5 animate-spin" aria-hidden="true" />Loading your workspace…</main>;
}
