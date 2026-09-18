import "server-only";
import { headers } from "next/headers";
import { currentUserSchema, type AuthUser } from "./schemas";
import { sessionCookieHeader } from "./session-cookie";

export function backendURL(path: string): URL {
  const configured = process.env.INTERNAL_API_URL;
  if (!configured) throw new Error("INTERNAL_API_URL is required.");
  const url = new URL(configured);
  if (!["http:", "https:"].includes(url.protocol) || url.username || url.password || url.pathname !== "/" || url.search || url.hash) {
    throw new Error("INTERNAL_API_URL must be an HTTP(S) origin.");
  }
  return new URL(path, url);
}

export async function getCurrentUser(): Promise<AuthUser | null> {
  const cookie = sessionCookieHeader((await headers()).get("cookie"));
  if (!cookie) return null;

  const response = await fetch(backendURL("/auth/me"), {
    headers: { Cookie: cookie, Accept: "application/json" },
    cache: "no-store",
    redirect: "error",
    signal: AbortSignal.timeout(10_000),
  });
  if (response.status === 401) return null;
  // An outage is an error, not a signed-out user; route error.tsx provides retry.
  if (!response.ok) throw new Error("Authentication service unavailable.");
  return currentUserSchema.parse(await response.json()).user;
}
