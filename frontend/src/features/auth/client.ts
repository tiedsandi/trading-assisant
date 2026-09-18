import { currentUserSchema, type LoginValues } from "./schemas";

export async function submitAuth(action: "register" | "login" | "logout", values?: LoginValues) {
  let response: Response;
  try {
    response = await fetch(`/api/auth/${action}`, {
      method: "POST",
      credentials: "same-origin",
      cache: "no-store",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(values ? { email: values.email, password: values.password } : {}),
      signal: AbortSignal.timeout(15_000),
    });
  } catch {
    throw new Error("Unable to connect. Please try again.");
  }

  if (!response.ok) {
    if (action === "login" && response.status === 401) {
      throw new Error("Invalid email or password.");
    }
    if (action === "register" && response.status === 409) {
      throw new Error("Unable to create an account with these details. Try signing in instead.");
    }
    if (response.status === 429) throw new Error("Too many attempts. Please try again shortly.");
    if (response.status === 400) throw new Error("Please check your details and try again.");
    if (response.status === 403) throw new Error("Your request could not be verified. Refresh the page and try again.");
    throw new Error("Something went wrong. Please try again.");
  }

  if (action !== "logout") {
    try {
      currentUserSchema.parse(await response.json());
    } catch {
      throw new Error("Unable to confirm your session. Please try again.");
    }
  }
}
