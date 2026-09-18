import { act, fireEvent, render, screen, waitFor } from "@testing-library/react";
import { useRouter } from "next/navigation";
import { submitAuth } from "./client";
import { RegisterForm } from "./register-form";
import { LoginForm } from "./login-form";
import { LogoutButton } from "./logout-button";

jest.mock("next/navigation", () => ({ useRouter: jest.fn() }));
jest.mock("./client", () => ({ submitAuth: jest.fn() }));

const replace = jest.fn();
const refresh = jest.fn();
const submit = jest.mocked(submitAuth);

beforeEach(() => {
  jest.mocked(useRouter).mockReturnValue({ replace, refresh } as unknown as ReturnType<typeof useRouter>);
  submit.mockResolvedValue(undefined);
});
function fillRegister(password = "abcdefgh") {
  fireEvent.change(screen.getByLabelText("Email"), { target: { value: "Sandi@Example.com" } });
  fireEvent.change(screen.getByLabelText("Password"), { target: { value: password } });
  fireEvent.change(screen.getByLabelText("Confirm password"), { target: { value: password } });
}

describe("registration form", () => {
  it("starts without errors, then gives live feedback after interaction", async () => {
    render(<RegisterForm />);
    expect(screen.queryByRole("alert")).not.toBeInTheDocument();
    expect(screen.getByText("Use 8–128 characters")).toBeVisible();
    const email = screen.getByLabelText("Email");
    fireEvent.change(email, { target: { value: "invalid" } });
    fireEvent.blur(email);
    expect(await screen.findByText("Enter a valid email address.")).toBeVisible();
    expect(email).toHaveAttribute("aria-invalid", "true");
    fireEvent.change(email, { target: { value: "a@example.com" } });
    await waitFor(() => expect(screen.queryByText("Enter a valid email address.")).not.toBeInTheDocument());
    fireEvent.change(screen.getByLabelText("Password"), { target: { value: "abcdefgh" } });
    expect(await screen.findByText("Password length looks good")).toBeVisible();
  });

  it("checks confirmation and rechecks it when the original password changes", async () => {
    render(<RegisterForm />);
    fillRegister();
    fireEvent.change(screen.getByLabelText("Confirm password"), { target: { value: "mismatch" } });
    fireEvent.blur(screen.getByLabelText("Confirm password"));
    expect(await screen.findByText("Passwords do not match.")).toBeVisible();
    fireEvent.change(screen.getByLabelText("Password"), { target: { value: "mismatch" } });
    await waitFor(() => expect(screen.queryByText("Passwords do not match.")).not.toBeInTheDocument());
  });

  it("blocks invalid submission and focuses the first invalid field", async () => {
    render(<RegisterForm />);
    fireEvent.click(screen.getByRole("button", { name: "Create account" }));
    expect(await screen.findByText("Use at least 8 characters.")).toBeVisible();
    await waitFor(() => expect(screen.getByLabelText("Email")).toHaveFocus());
    expect(submit).not.toHaveBeenCalled();
  });

  it("normalizes values, disables repeated submission, then enters the dashboard immediately", async () => {
    let complete!: () => void;
    submit.mockImplementationOnce(() => new Promise<void>((resolve) => { complete = resolve; }));
    render(<RegisterForm />);
    fillRegister();
    fireEvent.click(screen.getByRole("button", { name: "Create account" }));
    const loading = await screen.findByRole("button", { name: "Creating account…" });
    expect(loading).toBeDisabled();
    fireEvent.click(loading);
    expect(submit).toHaveBeenCalledTimes(1);
    expect(submit).toHaveBeenCalledWith("register", { email: "sandi@example.com", password: "abcdefgh", confirmPassword: "abcdefgh" });
    await act(async () => complete());
    expect(replace).toHaveBeenCalledWith("/dashboard");
    expect(refresh).toHaveBeenCalledTimes(1);
  });

  it("shows backend failure without navigating and allows retry", async () => {
    submit.mockRejectedValueOnce(new Error("Unable to create an account with these details."));
    render(<RegisterForm />);
    fillRegister();
    fireEvent.click(screen.getByRole("button", { name: "Create account" }));
    expect(await screen.findByRole("alert")).toHaveTextContent("Unable to create an account with these details.");
    expect(replace).not.toHaveBeenCalled();
    expect(screen.getByRole("button", { name: "Create account" })).toBeEnabled();
  });

  it("toggles password visibility independently without changing values", () => {
    render(<RegisterForm />);
    fillRegister();
    fireEvent.click(screen.getByRole("button", { name: "Show password" }));
    expect(screen.getByLabelText("Password")).toHaveAttribute("type", "text");
    expect(screen.getByLabelText("Confirm password")).toHaveAttribute("type", "password");
    expect(screen.getByLabelText("Password")).toHaveValue("abcdefgh");
    fireEvent.click(screen.getByRole("button", { name: "Hide password" }));
    expect(screen.getByLabelText("Password")).toHaveAttribute("type", "password");
  });
});
describe("login form", () => {
  it("validates email and nonempty password without displaying strength rules", async () => {
    render(<LoginForm />);
    fireEvent.click(screen.getByRole("button", { name: "Sign in" }));
    expect(await screen.findByText("Enter your password.")).toBeVisible();
    expect(screen.getByText("Enter a valid email address.")).toBeVisible();
    expect(screen.queryByText(/8.*characters/)).not.toBeInTheDocument();
    expect(submit).not.toHaveBeenCalled();
  });

  it("accepts a nonempty short login password and redirects on success", async () => {
    render(<LoginForm />);
    fireEvent.change(screen.getByLabelText("Email"), { target: { value: "SANDI@example.com" } });
    fireEvent.change(screen.getByLabelText("Password"), { target: { value: "a" } });
    fireEvent.click(screen.getByRole("button", { name: "Sign in" }));
    await waitFor(() => expect(replace).toHaveBeenCalledWith("/dashboard"));
    expect(submit).toHaveBeenCalledWith("login", { email: "sandi@example.com", password: "a" });
  });

  it("shows generic credential errors and remains on login", async () => {
    submit.mockRejectedValueOnce(new Error("Invalid email or password."));
    render(<LoginForm />);
    fireEvent.change(screen.getByLabelText("Email"), { target: { value: "sandi@example.com" } });
    fireEvent.change(screen.getByLabelText("Password"), { target: { value: "incorrect" } });
    fireEvent.click(screen.getByRole("button", { name: "Sign in" }));
    expect(await screen.findByRole("alert")).toHaveTextContent("Invalid email or password.");
    expect(replace).not.toHaveBeenCalled();
  });
});

describe("logout", () => {
  it("waits for revocation before returning to login", async () => {
    let complete!: () => void;
    submit.mockImplementationOnce(() => new Promise<void>((resolve) => { complete = resolve; }));
    render(<LogoutButton />);
    fireEvent.click(screen.getByRole("button", { name: "Sign out" }));
    expect(await screen.findByRole("button", { name: "Signing out…" })).toBeDisabled();
    expect(replace).not.toHaveBeenCalled();
    await act(async () => complete());
    expect(submit).toHaveBeenCalledWith("logout");
    expect(replace).toHaveBeenCalledWith("/login");
    expect(refresh).toHaveBeenCalledTimes(1);
  });

  it("keeps the user on the page if logout fails", async () => {
    submit.mockRejectedValueOnce(new Error("Network down"));
    render(<LogoutButton />);
    fireEvent.click(screen.getByRole("button", { name: "Sign out" }));
    expect(await screen.findByRole("alert")).toHaveTextContent("Unable to sign out.");
    expect(screen.getByRole("button", { name: "Sign out" })).toBeEnabled();
    expect(replace).not.toHaveBeenCalled();
  });
});
