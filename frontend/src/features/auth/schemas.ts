import { z } from "zod";

// Conventional ASCII email addresses; normalization matches the Go authority.
export const emailSchema = z.string().trim().toLowerCase().max(254)
  .refine((email) => {
    const parts = email.split("@");
    const [local, domain] = parts;
    return parts.length === 2 && !!local && !!domain && local.length <= 64
      && /^[a-z0-9!#$%&'*+/=?^_`{|}~.-]+$/.test(local)
      && !local.startsWith(".") && !local.endsWith(".") && !local.includes("..")
      && domain.includes(".") && domain.split(".").every((label) =>
        label.length <= 63 && /^[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$/.test(label));
  }, "Enter a valid email address.");

export const passwordLength = (password: string) => Array.from(password).length;

export const registerSchema = z.object({
  email: emailSchema,
  password: z.string()
    .refine((value) => passwordLength(value) >= 8, "Use at least 8 characters.")
    .refine((value) => passwordLength(value) <= 128, "Use no more than 128 characters."),
  confirmPassword: z.string().min(1, "Confirm your password."),
}).refine((values) => values.password === values.confirmPassword, {
  message: "Passwords do not match.",
  path: ["confirmPassword"],
});

export const loginSchema = z.object({
  email: emailSchema,
  password: z.string().min(1, "Enter your password."),
});

export type RegisterValues = z.infer<typeof registerSchema>;
export type LoginValues = z.infer<typeof loginSchema>;

export const currentUserSchema = z.object({
  user: z.object({
    id: z.string().min(1),
    email: z.string().min(1),
    created_at: z.string().min(1),
  }),
});

export type AuthUser = z.infer<typeof currentUserSchema>["user"];
