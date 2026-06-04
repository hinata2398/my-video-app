"use client";

import Link from "next/link";
import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";

export default function Header() {
  const router = useRouter();
  const [loggedIn, setLoggedIn] = useState(false);

  useEffect(() => {
    setLoggedIn(!!localStorage.getItem("token"));
  }, []);

  const handleLogout = () => {
    localStorage.removeItem("token");
    setLoggedIn(false);
    router.push("/");
  };

  return (
    <header style={{
      background: "#111", color: "#fff", padding: "0 2rem",
      display: "flex", alignItems: "center", justifyContent: "space-between",
      height: 56, position: "sticky", top: 0, zIndex: 100,
    }}>
      <Link href="/" style={{ color: "#fff", textDecoration: "none", fontWeight: "bold", fontSize: "1.2rem" }}>
        🎬 MyVideoApp
      </Link>
      <nav style={{ display: "flex", gap: "1rem", alignItems: "center" }}>
        {loggedIn ? (
          <>
            <Link href="/mypage" style={{ color: "#ccc", textDecoration: "none", fontSize: "0.9rem" }}>
              マイページ
            </Link>
            <Link href="/videos/new" style={{
              background: "#e00", color: "#fff", textDecoration: "none",
              padding: "0.4rem 1rem", borderRadius: 4, fontSize: "0.9rem",
            }}>
              + 投稿する
            </Link>
            <button onClick={handleLogout} style={{
              background: "none", border: "1px solid #555", color: "#ccc",
              padding: "0.4rem 1rem", borderRadius: 4, cursor: "pointer", fontSize: "0.9rem",
            }}>
              ログアウト
            </button>
          </>
        ) : (
          <Link href="/auth" style={{
            background: "#e00", color: "#fff", textDecoration: "none",
            padding: "0.4rem 1rem", borderRadius: 4, fontSize: "0.9rem",
          }}>
            ログイン
          </Link>
        )}
      </nav>
    </header>
  );
}
