import { redirect } from "next/navigation";
import { getCurrentUser } from "@/features/auth/server";

export default async function Home() {
  redirect((await getCurrentUser()) ? "/dashboard" : "/login");
}
