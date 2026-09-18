"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { ArrowRight, LoaderCircle } from "lucide-react";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { PasswordInput } from "./password-input";
import { loginSchema, type LoginValues } from "./schemas";
import { submitAuth } from "./client";

export function LoginForm() {
  const router = useRouter();
  const [serverError, setServerError] = useState("");
  const { register, handleSubmit, formState: { errors, isSubmitting } } = useForm<LoginValues>({
    resolver: zodResolver(loginSchema), mode: "onTouched", defaultValues: { email: "", password: "" },
  });
  const onSubmit = handleSubmit(async (values) => {
    setServerError("");
    try {
      await submitAuth("login", values);
      router.replace("/dashboard");
      router.refresh();
    } catch (error) {
      setServerError(error instanceof Error ? error.message : "Unable to sign in. Please try again.");
    }
  });

  return (
    <form noValidate onSubmit={onSubmit} className="space-y-5" aria-busy={isSubmitting}>
      {serverError && <Alert variant="destructive"><AlertDescription>{serverError}</AlertDescription></Alert>}
      <div className="space-y-2">
        <Label htmlFor="login-email">Email</Label>
        <Input id="login-email" type="email" autoComplete="email" autoCapitalize="none" spellCheck={false} placeholder="you@example.com" className="h-12 bg-card text-base" readOnly={isSubmitting} aria-invalid={!!errors.email} aria-describedby={errors.email ? "login-email-error" : undefined} {...register("email")} />
        {errors.email && <p id="login-email-error" className="text-sm text-destructive" role="alert">{errors.email.message}</p>}
      </div>
      <div className="space-y-2">
        <Label htmlFor="login-password">Password</Label>
        <PasswordInput id="login-password" autoComplete="current-password" readOnly={isSubmitting} aria-invalid={!!errors.password} aria-describedby={errors.password ? "login-password-error" : undefined} {...register("password")} />
        {errors.password && <p id="login-password-error" className="text-sm text-destructive" role="alert">{errors.password.message}</p>}
      </div>
      <Button type="submit" className="mt-2 h-12 w-full text-sm" disabled={isSubmitting}>
        {isSubmitting ? <><LoaderCircle className="size-4 animate-spin" aria-hidden="true" />Signing in…</> : <>Sign in<ArrowRight className="size-4" aria-hidden="true" /></>}
      </Button>
    </form>
  );
}
