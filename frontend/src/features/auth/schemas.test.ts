import { emailSchema, loginSchema, registerSchema } from "./schemas";

describe("auth validation", () => {
  it("normalizes email and accepts an eight-character password without complexity rules", () => {
    expect(registerSchema.parse({ email: "  SANDI@Example.COM  ", password: "aaaaaaaa", confirmPassword: "aaaaaaaa" })).toEqual({
      email: "sandi@example.com", password: "aaaaaaaa", confirmPassword: "aaaaaaaa",
    });
  });

  it.each(["", "email", "a@localhost", "a@b.com@", "a@b.com@@extra", "a..b@example.com", ".a@example.com", "a@-example.com", "a@exam_ple.com", "é@example.com"])("rejects invalid email %j", (email) => {
    expect(emailSchema.safeParse(email).success).toBe(false);
  });

  it.each(["a+b@example.co.id", "o'connor@example.com", "a@123.example"])("accepts conventional email %j", (email) => {
    expect(emailSchema.safeParse(email).success).toBe(true);
  });

  it.each(["1234567", "a".repeat(129), "😀".repeat(7)])("rejects password outside the Unicode character limits", (password) => {
    expect(registerSchema.safeParse({ email: "a@example.com", password, confirmPassword: password }).success).toBe(false);
  });

  it("counts Unicode code points consistently with the backend", () => {
    const password = "😀".repeat(8);
    expect(registerSchema.safeParse({ email: "a@example.com", password, confirmPassword: password }).success).toBe(true);
  });

  it("attaches mismatch feedback to confirmation", () => {
    const result = registerSchema.safeParse({ email: "a@example.com", password: "password", confirmPassword: "different" });
    expect(result.success).toBe(false);
    if (!result.success) expect(result.error.issues).toContainEqual(expect.objectContaining({ path: ["confirmPassword"], message: "Passwords do not match." }));
  });

  it("requires only a nonempty password on login", () => {
    expect(loginSchema.safeParse({ email: "a@example.com", password: "a" }).success).toBe(true);
    expect(loginSchema.safeParse({ email: "a@example.com", password: "" }).success).toBe(false);
  });
});
