/** @jest-environment node */
import { proxyAuth } from "./proxy";
import { sessionCookieHeader } from "./session-cookie";

jest.mock("server-only", () => ({}));
const originalURL = process.env.INTERNAL_API_URL;
beforeEach(() => { process.env.INTERNAL_API_URL = "http://backend:8080"; });
afterAll(() => {
  if (originalURL === undefined) delete process.env.INTERNAL_API_URL;
  else process.env.INTERNAL_API_URL = originalURL;
});

function post(headers: HeadersInit = {}, body = "{}") {
  return new Request("http://localhost:3000/api/auth/login", { method: "POST", headers, body });
}

it("forwards the raw origin and only session cookies, preserving secure Set-Cookie and no-store", async () => {
  const tokenCookie = "__Host-ta_session=token; Path=/; HttpOnly; Secure; SameSite=Lax; Max-Age=604800";
  const fetchMock = jest.spyOn(globalThis, "fetch").mockResolvedValue(Response.json({ user: { id: "one" } }, { headers: { "Set-Cookie": tokenCookie, "X-Internal-Secret": "private" } }));
  const response = await proxyAuth(post({ Origin: "https://evil.example", "Content-Type": "application/json", Cookie: "other=secret; ta_session=valid_token", Authorization: "private" }), "login");
  expect(response.status).toBe(200);
  const [url, options] = fetchMock.mock.calls[0];
  expect(url).toEqual(new URL("http://backend:8080/auth/login"));
  const sentHeaders = new Headers(options?.headers);
  expect(sentHeaders.get("origin")).toBe("https://evil.example");
  expect(sentHeaders.get("cookie")).toBe("ta_session=valid_token");
  expect(sentHeaders.has("authorization")).toBe(false);
  expect(options).toMatchObject({ method: "POST", cache: "no-store", redirect: "error" });
  expect(response.headers.get("Set-Cookie")).toBe(tokenCookie);
  expect(response.headers.get("Cache-Control")).toBe("no-store");
  expect(response.headers.has("X-Internal-Secret")).toBe(false);
  expect(response.headers.has("Access-Control-Allow-Origin")).toBe(false);
});

it("does not fabricate Origin when a request omits it", async () => {
  const fetchMock = jest.spyOn(globalThis, "fetch").mockResolvedValue(new Response(null, { status: 403 }));
  const response = await proxyAuth(post(), "logout");
  expect(response.status).toBe(403);
  expect(new Headers(fetchMock.mock.calls[0][1]?.headers).has("origin")).toBe(false);
});

it.each([["anything", "POST", 404], ["login", "GET", 405], ["me", "POST", 405], ["https://evil.example", "GET", 404]])("restricts action %s and method %s", async (action, method, status) => {
  const fetchMock = jest.spyOn(globalThis, "fetch");
  const response = await proxyAuth(new Request("http://localhost:3000/api/auth/test", { method }), action);
  expect(response.status).toBe(status);
  expect(fetchMock).not.toHaveBeenCalled();
});

it("preserves logout cookie deletion and the empty 204 response", async () => {
  jest.spyOn(globalThis, "fetch").mockResolvedValue(new Response(null, { status: 204, headers: { "Set-Cookie": "ta_session=; Path=/; HttpOnly; Max-Age=0; SameSite=Lax" } }));
  const response = await proxyAuth(post({ Origin: "http://localhost:3000", "Content-Type": "application/json" }), "logout");
  expect(response.status).toBe(204);
  expect(await response.text()).toBe("");
  expect(response.headers.get("set-cookie")).toContain("Max-Age=0");
});

it("does not forward unrelated upstream cookies", async () => {
  jest.spyOn(globalThis, "fetch").mockResolvedValue(new Response(null, { headers: { "Set-Cookie": "other=secret" } }));
  expect((await proxyAuth(post(), "login")).headers.has("set-cookie")).toBe(false);
});

it.each([false, true])("rejects oversized bodies with and without Content-Length (declared=%s)", async (declared) => {
  const fetchMock = jest.spyOn(globalThis, "fetch");
  const response = await proxyAuth(post(declared ? { "Content-Length": "20000" } : {}, "x".repeat(20000)), "login");
  expect(response.status).toBe(413);
  expect(fetchMock).not.toHaveBeenCalled();
});

it("returns a generic no-store service error for an unreachable backend", async () => {
  jest.spyOn(globalThis, "fetch").mockRejectedValue(new Error("internal secret"));
  const response = await proxyAuth(post(), "login");
  expect(response.status).toBe(502);
  expect(await response.json()).toEqual({ error: { code: "service_unavailable", message: "Authentication service unavailable." } });
  expect(response.headers.get("cache-control")).toBe("no-store");
});

it("filters unrelated and malformed cookies", () => {
  expect(sessionCookieHeader("analytics=secret; ta_session=abc_def-123; other=1; __Host-ta_session=second")).toBe("ta_session=abc_def-123; __Host-ta_session=second");
  expect(sessionCookieHeader("ta_session=quoted token; other=1")).toBe("");
  expect(sessionCookieHeader(null)).toBe("");
});
