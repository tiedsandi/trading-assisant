import "server-only";
import { backendURL } from "./server";
import { sessionCookieHeader } from "./session-cookie";

const bodyLimit = 16 * 1024;

function errorResponse(status: number, code: string, message: string, allow?: string) {
  return Response.json({ error: { code, message } }, {
    status,
    headers: { "Cache-Control": "no-store", ...(allow ? { Allow: allow } : {}) },
  });
}
async function readBody(request: Request): Promise<Uint8Array<ArrayBuffer>> {
  if (Number(request.headers.get("content-length")) > bodyLimit) throw new RangeError();
  const reader = request.body?.getReader();
  if (!reader) return new Uint8Array();
  const chunks: Uint8Array[] = [];
  let length = 0;
  try {
    while (true) {
      const { done, value } = await reader.read();
      if (done) break;
      length += value.length;
      if (length > bodyLimit) throw new RangeError();
      chunks.push(value);
    }
  } finally {
    await reader.cancel();
    reader.releaseLock();
  }
  const body = new Uint8Array(length);
  let offset = 0;
  for (const chunk of chunks) {
    body.set(chunk, offset);
    offset += chunk.length;
  }
  return body;
}

export async function proxyAuth(request: Request, action: string) {
  const allowed = action === "me" ? "GET" : ["register", "login", "logout"].includes(action) ? "POST" : null;
  if (!allowed) return errorResponse(404, "not_found", "Not found.");
  if (request.method !== allowed) return errorResponse(405, "method_not_allowed", "Method not allowed.", allowed);

  try {
    const headers = new Headers({ Accept: "application/json" });
    for (const name of ["origin", "content-type"]) {
      const value = request.headers.get(name);
      if (value !== null) headers.set(name, value);
    }
    const cookie = sessionCookieHeader(request.headers.get("cookie"));
    if (cookie) headers.set("Cookie", cookie);

    const body = allowed === "POST" ? await readBody(request) : undefined;
    const upstream = await fetch(backendURL(`/auth/${action}`), {
      method: allowed,
      headers,
      body,
      cache: "no-store",
      redirect: "error",
      signal: AbortSignal.any([request.signal, AbortSignal.timeout(10_000)]),
    });
    const responseHeaders = new Headers({ "Cache-Control": "no-store" });
    const contentType = upstream.headers.get("content-type");
    if (contentType) responseHeaders.set("Content-Type", contentType);
    for (const cookie of upstream.headers.getSetCookie()) {
      if (/^(?:ta_session|__Host-ta_session)=/.test(cookie)) responseHeaders.append("Set-Cookie", cookie);
    }
    // Never expose arbitrary upstream headers, follow redirects, or accept a target URL.
    return new Response(upstream.body, { status: upstream.status, headers: responseHeaders });
  } catch (error) {
    if (error instanceof RangeError) return errorResponse(413, "request_too_large", "Request is too large.");
    return errorResponse(502, "service_unavailable", "Authentication service unavailable.");
  }
}
