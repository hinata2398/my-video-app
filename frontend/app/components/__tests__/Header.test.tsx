import { render, screen, waitFor } from "@testing-library/react";
import { describe, it, expect, vi, beforeEach } from "vitest";
import Header from "../Header";
import userEvent from "@testing-library/user-event";

const mockPush = vi.fn();
vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: mockPush }),
}));

describe("Header", () => {
  beforeEach(() => {
    localStorage.clear();
    vi.resetAllMocks();
  });

  it("未ログイン時はログインリンクが表示される", () => {
    render(<Header />);
    expect(screen.getByText("ログイン")).toBeInTheDocument();
  });

  it("ログイン済みのときはマイページ・投稿する・ログアウトが表示される", async () => {
    localStorage.setItem("token", "test-token");
    global.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({
        username: "testuser",
        avatar_url: "",
        email: "test@example.com",
      }),
    });

    render(<Header />);

    await waitFor(() => {
      expect(screen.getByText("マイページ")).toBeInTheDocument();
      expect(screen.getByText("ログアウト")).toBeInTheDocument();
    });
  });

  it("ログアウトするとログイン画面に戻る", async () => {
    // ログイン済み状態を作る
    localStorage.setItem("token", "test-token");
    global.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({
        username: "testuser",
        avatar_url: "",
        email: "test@example.com",
      }),
    });

    render(<Header />);

    // ログアウトボタンが表示されるまで待つ
    await waitFor(() => {
      expect(screen.getByText("ログアウト")).toBeInTheDocument();
    });

    // クリック
    await userEvent.click(screen.getByText("ログアウト"));

    // 検証
    expect(screen.getByText("ログイン")).toBeInTheDocument();
    expect(localStorage.getItem("token")).toBeNull();
    expect(mockPush).toHaveBeenCalledWith("/");
  });
});
