import { render, screen, within } from "@testing-library/react";
import Home from "@/app/page";

describe("Home scaffold", () => {
  it("renders the starter instructions inside the main landmark", () => {
    render(<Home />);

    const main = screen.getByRole("main");
    expect(
      within(main).getByRole("heading", {
        level: 1,
        name: /to get started, edit the page\.tsx file\./i,
      }),
    ).toBeVisible();
  });

  it("provides a documentation link with a safe external target", () => {
    render(<Home />);

    const link = screen.getByRole("link", { name: "Documentation" });
    expect(link).toHaveAttribute(
      "href",
      expect.stringMatching(/^https:\/\/nextjs\.org\/docs(?:\?|$)/),
    );
    expect(link).toHaveAttribute("target", "_blank");
    expect(link.getAttribute("rel")?.split(/\s+/)).toEqual(
      expect.arrayContaining(["noopener", "noreferrer"]),
    );
  });
});
