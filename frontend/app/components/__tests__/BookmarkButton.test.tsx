import { render, screen, waitFor } from "@testing-library/react";
import { describe, it, expect, vi, beforeEach } from "vitest";
import userEvent from "@testing-library/user-event";
import BookmarkButton from "../BookmarkButton";

describe("BookmarkButton", () => {
  beforeEach(() => {
    localStorage.clear();
    vi.resetAllMocks();
  });

  it("未ブックマーク時は空のしおりアイコンが表示される", () => {
    render(<BookmarkButton videoID={1} initialBookmarked={false} />);
    expect(
      screen.getByRole("button", { name: "ブックマーク" }),
    ).toBeInTheDocument();
  });

  it("ブックマーク済み時は塗りつぶしのしおりアイコンが表示される", () => {
    render(<BookmarkButton videoID={1} initialBookmarked={true} />);
    expect(
      screen.getByRole("button", { name: "ブックマーク済み" }),
    ).toBeInTheDocument();
  });

  it("クリックするとトグルされる", async () => {
    localStorage.setItem("token", "test-token");
    render(<BookmarkButton videoID={1} initialBookmarked={false} />);
    global.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ bookmarked: true }),
    });
    await userEvent.click(screen.getByRole("button", { name: "ブックマーク" }));
    await waitFor(() => {
      expect(
        screen.getByRole("button", { name: "ブックマーク済み" }),
      ).toBeInTheDocument();
    });
  });
});
