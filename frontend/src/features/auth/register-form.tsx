"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useForm, useWatch } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { ArrowRight, Check, Circle, LoaderCircle } from "lucide-react";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { PasswordInput } from "./password-input";
import { passwordLength, registerSchema, type RegisterValues } from "./schemas";
import { submitAuth } from "./client";

export function RegisterForm() {
  const router = useRouter();
  const [serverError, setServerError] = useState("");
  const { register, control, handleSubmit, trigger, formState: { errors, isSubmitting, touchedFields } } = useForm<RegisterValues>({
    resolver: zodResolver(registerSchema),
    mode: "onTouched",
    defaultValues: { email: "", password: "", confirmPassword: "" },
  });
  const password = useWatch({ control, name: "password" });
  const validLength = passwordLength(password) >= 8 && passwordLength(password) <= 128;
  const LengthIcon = validLength ? Check : Circle;
  useEffect(() => {
    if (touchedFields.confirmPassword) void trigger("confirmPassword");
  }, [password, touchedFields.confirmPassword, trigger]);

  const onSubmit = handleSubmit(async (values) => {
    setServerError("");
    try {
      await submitAuth("register", values);
      router.replace("/dashboard");
      router.refresh();
    } catch (error) {
      setServerError(error instanceof Error ? error.message : "Unable to create your account. Please try again.");
    }
  });

  return (
    <form noValidate onSubmit={onSubmit} className="space-y-5" aria-busy={isSubmitting}>
      {serverError && <Alert variant="destructive"><AlertDescription>{serverError}</AlertDescription></Alert>}
      <div className="space-y-2">
        <Label htmlFor="register-email">Email</Label>
        <Input id="register-email" type="email" autoComplete="email" autoCapitalize="none" spellCheck={false} className="h-12 bg-card text-base" placeholder="you@example.com" readOnly={isSubmitting} aria-invalid={!!errors.email} aria-describedby={errors.email ? "register-email-error" : undefined} {...register("email")} />
        {errors.email && <p id="register-email-error" className="text-sm text-destructive" role="alert">{errors.email.message}</p>}
      </div>
      <div className="space-y-2">
        <Label htmlFor="register-password">Password</Label>
        <PasswordInput id="register-password" autoComplete="new-password" readOnly={isSubmitting} aria-invalid={!!errors.password} aria-describedby={`password-hint${errors.password ? " register-password-error" : ""}`} {...register("password")} />
        <p id="password-hint" className={`flex items-center gap-2 text-xs ${validLength ? "text-primary" : "text-muted-foreground"}`} aria-live="polite"><LengthIcon className="size-3.5" aria-hidden="true" />{validLength ? "Password length looks good" : "Use 8–128 characters"}</p>
        {errors.password && <p id="register-password-error" className="text-sm text-destructive" role="alert">{errors.password.message}</p>}
      </div>
      <div className="space-y-2">
        <Label htmlFor="confirm-password">Confirm password</Label>
        <PasswordInput id="confirm-password" label="confirm password" autoComplete="new-password" readOnly={isSubmitting} aria-invalid={!!errors.confirmPassword} aria-describedby={errors.confirmPassword ? "confirm-password-error" : undefined} {...register("confirmPassword")} />
        {errors.confirmPassword && <p id="confirm-password-error" className="text-sm text-destructive" role="alert">{errors.confirmPassword.message}</p>}
      </div>
      <Button type="submit" className="mt-2 h-12 w-full text-sm" disabled={isSubmitting}>
        {isSubmitting ? <><LoaderCircle className="size-4 animate-spin" aria-hidden="true" />Creating account…</> : <>Create account<ArrowRight className="size-4" aria-hidden="true" /></>}
      </Button>
    </form>
  );
}
