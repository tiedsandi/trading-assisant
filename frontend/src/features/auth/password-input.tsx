"use client";

import { useState, type ComponentProps } from "react";
import { Eye, EyeOff } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

export function PasswordInput({ label = "password", ...props }: ComponentProps<typeof Input> & { label?: string }) {
  const [visible, setVisible] = useState(false);
  const Icon = visible ? EyeOff : Eye;
  return (
    <div className="relative">
      <Input {...props} type={visible ? "text" : "password"} className="h-12 bg-card pr-12 text-base" />
      <Button type="button" variant="ghost" size="icon" className="absolute top-1 right-1 size-10 text-muted-foreground" disabled={props.disabled} onClick={() => setVisible(!visible)} aria-label={`${visible ? "Hide" : "Show"} ${label}`} aria-pressed={visible}>
        <Icon className="size-4" aria-hidden="true" />
      </Button>
    </div>
  );
}
