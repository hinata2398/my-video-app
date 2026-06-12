"use client";

import Link from "next/link";
import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";

type Me = {
  username: string;
  avatar_url: string;
  email: string;
};

export default function Header() {
  const router = useRouter();
  const [me, setMe] = useState<Me | null>(null);

  useEffect(() => {
    const token = localStorage.getItem("token");
    if (!token) return;
    fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/me`, {
      headers: { Authorization: `Bearer ${token}` },
    })
      .then((res) => (res.ok ? res.json() : null))
      .then((data) => {
        if (data) setMe(data);
      })
      .catch(() => {});
  }, []);

  const handleLogout = () => {
    localStorage.removeItem("token");
    setMe(null);
    router.push("/");
  };

  return (
    <header
      style={{
        background: "#111",
        color: "#fff",
        padding: "0 2rem",
        display: "flex",
        alignItems: "center",
        justifyContent: "space-between",
        height: 56,
        position: "sticky",
        top: 0,
        zIndex: 100,
      }}
    >
      <Link
        href="/"
        style={{
          color: "#fff",
          textDecoration: "none",
          fontWeight: "bold",
          fontSize: "1.2rem",
        }}
      >
        🎬 MyVideoApp
      </Link>
      <nav style={{ display: "flex", gap: "1rem", alignItems: "center" }}>
        {me ? (
          <>
            <Link
              href="/mypage"
              style={{
                color: "#ccc",
                textDecoration: "none",
                fontSize: "0.9rem",
              }}
            >
              マイページ
            </Link>
            {/* アバター → プロフィール編集へ */}
            <Link href="/profile" style={{ textDecoration: "none" }}>
              <div
                style={{
                  width: 34,
                  height: 34,
                  borderRadius: "50%",
                  background: "#555",
                  overflow: "hidden",
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "center",
                  fontSize: "0.9rem",
                  fontWeight: "bold",
                  color: "#fff",
                  cursor: "pointer",
                  border: "2px solid #777",
                }}
              >
                {me.avatar_url ? (
                  <img
                    src={me.avatar_url}
                    alt="avatar"
                    style={{
                      width: "100%",
                      height: "100%",
                      objectFit: "cover",
                    }}
                  />
                ) : (
                  (me.username || me.email)[0].toUpperCase()
                )}
              </div>
            </Link>
            <button
              onClick={handleLogout}
              style={{
                background: "none",
                border: "1px solid #555",
                color: "#ccc",
                padding: "0.4rem 1rem",
                borderRadius: 4,
                cursor: "pointer",
                fontSize: "0.9rem",
              }}
            >
              ログアウト
            </button>
          </>
        ) : (
          <Link
            href="/auth"
            style={{
              background: "#e00",
              color: "#fff",
              textDecoration: "none",
              padding: "0.4rem 1rem",
              borderRadius: 4,
              fontSize: "0.9rem",
            }}
          >
            ログイン
          </Link>
        )}
      </nav>
    </header>
  );
}
