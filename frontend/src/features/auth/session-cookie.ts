// Only application session cookies cross the frontend/backend boundary.
export function sessionCookieHeader(header: string | null): string {
  return (header ?? "").split(";").map((cookie) => cookie.trim()).filter((cookie) =>
    /^(?:ta_session|__Host-ta_session)=[A-Za-z0-9_-]{1,256}$/.test(cookie),
  ).join("; ");
}
