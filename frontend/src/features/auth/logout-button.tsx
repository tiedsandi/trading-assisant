"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { LoaderCircle, LogOut } from "lucide-react";
import { Button } from "@/components/ui/button";
import { submitAuth } from "./client";

export function LogoutButton() {
  const router = useRouter();
  const [pending, setPending] = useState(false);
  const [error, setError] = useState("");
  async function logout() {
    if (pending) return;
    setPending(true);
    setError("");
    try {
      await submitAuth("logout");
      router.replace("/login");
      router.refresh();
    } catch {
      setError("Unable to sign out. Please try again.");
      setPending(false);
    }
  }
  return (
    <div className="flex flex-col items-end gap-2">
      <Button variant="outline" className="h-10" disabled={pending} onClick={logout}>
        {pending ? <LoaderCircle className="size-4 animate-spin" aria-hidden="true" /> : <LogOut className="size-4" aria-hidden="true" />}{pending ? "Signing out…" : "Sign out"}
      </Button>
      {error && <p className="max-w-56 text-right text-sm text-destructive" role="alert">{error}</p>}
    </div>
  );
}
