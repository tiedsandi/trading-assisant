/** @jest-environment node */
import { headers } from "next/headers";
import { backendURL, getCurrentUser } from "./server";

jest.mock("server-only", () => ({}));
jest.mock("next/headers", () => ({ headers: jest.fn() }));

const user = { id: "user-1", email: "a@example.com", created_at: "2026-09-15T00:00:00Z" };
const originalURL = process.env.INTERNAL_API_URL;

beforeEach(() => {
  process.env.INTERNAL_API_URL = "http://backend:8080";
  jest.mocked(headers).mockResolvedValue(new Headers({ Cookie: "analytics=private; ta_session=dev_token; __Host-ta_session=prod_token" }));
});
afterAll(() => {
  if (originalURL === undefined) delete process.env.INTERNAL_API_URL;
  else process.env.INTERNAL_API_URL = originalURL;
});

it("resolves safe current-user fields via private no-store fetch with only session cookies", async () => {
  const fetchMock = jest.spyOn(globalThis, "fetch").mockResolvedValue(Response.json({ user: { ...user, password_hash: "never expose" } }));
  await expect(getCurrentUser()).resolves.toEqual(user);
  expect(fetchMock).toHaveBeenCalledWith(new URL("http://backend:8080/auth/me"), expect.objectContaining({ cache: "no-store", redirect: "error", headers: { Cookie: "ta_session=dev_token; __Host-ta_session=prod_token", Accept: "application/json" } }));
});

it("treats missing and rejected sessions as guest", async () => {
  const fetchMock = jest.spyOn(globalThis, "fetch").mockResolvedValue(new Response(null, { status: 401 }));
  jest.mocked(headers).mockResolvedValueOnce(new Headers({ Cookie: "analytics=private" }));
  await expect(getCurrentUser()).resolves.toBeNull();
  expect(fetchMock).not.toHaveBeenCalled();
  await expect(getCurrentUser()).resolves.toBeNull();
});

it("surfaces backend outages instead of redirecting an authenticated user to login", async () => {
  const fetchMock = jest.spyOn(globalThis, "fetch").mockResolvedValue(new Response(null, { status: 503 }));
  await expect(getCurrentUser()).rejects.toThrow("Authentication service unavailable.");
  fetchMock.mockRejectedValueOnce(new TypeError("unreachable"));
  await expect(getCurrentUser()).rejects.toThrow("unreachable");
});

it.each(["ftp://backend", "http://user:secret@backend", "http://backend/prefix", "http://backend?target=evil"])("rejects an invalid internal backend origin %s", (url) => {
  process.env.INTERNAL_API_URL = url;
  expect(() => backendURL("/auth/me")).toThrow("HTTP(S) origin");
});
