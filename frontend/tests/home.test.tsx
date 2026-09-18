import { render, screen } from "@testing-library/react";
import { redirect } from "next/navigation";
import { getCurrentUser } from "@/features/auth/server";
import Home from "@/app/page";
import LoginPage from "@/app/login/page";
import RegisterPage from "@/app/register/page";
import DashboardPage from "@/app/dashboard/page";

jest.mock("@/features/auth/server", () => ({ getCurrentUser: jest.fn() }));
jest.mock("next/navigation", () => ({
  redirect: jest.fn((path: string) => { throw new Error(`REDIRECT:${path}`); }),
  useRouter: () => ({ replace: jest.fn(), refresh: jest.fn() }),
}));

const user = { id: "user-1", email: "sandi@example.com", created_at: "2026-09-15T00:00:00Z" };
const currentUser = jest.mocked(getCurrentUser);

beforeEach(() => { currentUser.mockResolvedValue(null); });

describe("home navigation", () => {
  it("directs a guest to login", async () => {
    await expect(Home()).rejects.toThrow("REDIRECT:/login");
    expect(redirect).toHaveBeenCalledWith("/login");
  });

  it("directs an authenticated user to their dashboard", async () => {
    currentUser.mockResolvedValue(user);
    await expect(Home()).rejects.toThrow("REDIRECT:/dashboard");
  });
});

describe("guest pages", () => {
  it("shows accessible login and registration navigation", async () => {
    render(await LoginPage());
    expect(screen.getByRole("main")).toBeVisible();
    expect(screen.getByRole("heading", { level: 1, name: "Sign in to your space" })).toBeVisible();
    expect(screen.getByLabelText("Email")).toBeVisible();
    expect(screen.getByRole("link", { name: "Create an account" })).toHaveAttribute("href", "/register");
  });

  it("shows registration with frontend password confirmation", async () => {
    render(await RegisterPage());
    expect(screen.getByRole("heading", { level: 1, name: "Create your account" })).toBeVisible();
    expect(screen.getByLabelText("Confirm password")).toBeVisible();
    expect(screen.getByRole("link", { name: "Sign in" })).toHaveAttribute("href", "/login");
  });

  it.each([LoginPage, RegisterPage])("redirects authenticated visitors from guest pages", async (Page) => {
    currentUser.mockResolvedValue(user);
    await expect(Page()).rejects.toThrow("REDIRECT:/dashboard");
  });
});

describe("authenticated dashboard", () => {
  it("rejects guests before rendering authenticated content", async () => {
    await expect(DashboardPage()).rejects.toThrow("REDIRECT:/login");
  });

  it("shows safe account information and logout for authenticated users", async () => {
    currentUser.mockResolvedValue(user);
    render(await DashboardPage());
    expect(screen.getByRole("heading", { level: 1, name: "Welcome to Trading Assistant." })).toBeVisible();
    expect(screen.getByText(user.email)).toBeVisible();
    expect(screen.getByRole("button", { name: "Sign out" })).toBeVisible();
    expect(screen.queryByText(user.id)).not.toBeInTheDocument();
  });

  it("surfaces an auth outage without pretending the user is a guest", async () => {
    currentUser.mockRejectedValue(new Error("Authentication service unavailable."));
    await expect(DashboardPage()).rejects.toThrow("Authentication service unavailable.");
    expect(redirect).not.toHaveBeenCalled();
  });
});
