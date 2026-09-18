/** @jest-environment node */
import { submitAuth } from "./client";

const user = { id: "user-1", email: "a@example.com", created_at: "2026-09-15T00:00:00Z" };

it("sends same-origin JSON credentials with no confirmation field or browser token", async () => {
  const fetchMock = jest.spyOn(globalThis, "fetch").mockResolvedValue(Response.json({ user }));
  await submitAuth("register", { email: "a@example.com", password: "abcdefgh", ...{ confirmPassword: "abcdefgh" } });
  expect(fetchMock).toHaveBeenCalledWith("/api/auth/register", expect.objectContaining({
    method: "POST", credentials: "same-origin", cache: "no-store",
    body: JSON.stringify({ email: "a@example.com", password: "abcdefgh" }),
  }));
});

it("sends an empty JSON object for logout and accepts 204", async () => {
  const fetchMock = jest.spyOn(globalThis, "fetch").mockResolvedValue(new Response(null, { status: 204 }));
  await expect(submitAuth("logout")).resolves.toBeUndefined();
  expect(fetchMock).toHaveBeenCalledWith("/api/auth/logout", expect.objectContaining({ method: "POST", body: "{}" }));
});

it.each(["unknown email", "wrong password"])("shows identical invalid credential errors for %s", async (message) => {
  jest.spyOn(globalThis, "fetch").mockResolvedValue(Response.json({ error: { message } }, { status: 401 }));
  await expect(submitAuth("login", { email: "a@example.com", password: "wrong" })).rejects.toThrow("Invalid email or password.");
});

it("shows usable duplicate registration feedback", async () => {
  jest.spyOn(globalThis, "fetch").mockResolvedValue(new Response(null, { status: 409 }));
  await expect(submitAuth("register", { email: "a@example.com", password: "abcdefgh" })).rejects.toThrow("Try signing in instead.");
});

it("does not echo unsafe server failures or treat malformed success as authenticated", async () => {
  const fetchMock = jest.spyOn(globalThis, "fetch").mockResolvedValue(Response.json({ error: { message: "database secret" } }, { status: 500 }));
  await expect(submitAuth("login")).rejects.toThrow("Something went wrong. Please try again.");
  fetchMock.mockResolvedValueOnce(Response.json({ password_hash: "unsafe" }));
  await expect(submitAuth("login")).rejects.toThrow("Unable to confirm your session.");
});

it("reports network failure cleanly", async () => {
  jest.spyOn(globalThis, "fetch").mockRejectedValue(new TypeError("network failed"));
  await expect(submitAuth("login")).rejects.toThrow("Unable to connect.");
});
