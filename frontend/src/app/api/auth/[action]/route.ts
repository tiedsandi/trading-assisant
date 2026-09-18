import { proxyAuth } from "@/features/auth/proxy";

type Context = { params: Promise<{ action: string }> };

async function handle(request: Request, context: Context) {
  return proxyAuth(request, (await context.params).action);
}

export const GET = handle;
export const POST = handle;
