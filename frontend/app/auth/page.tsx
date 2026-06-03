"use client";

import { useState } from "react";

export default function AuthPage() {
  const [mode, setMode] = useState<"login" | "register">("login");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    const endpoint = mode === "login" ? "/api/auth/login" : "/api/auth/register";
    const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}${endpoint}`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password }),
    });
    if (!res.ok) {
      setError(await res.text());
      return;
    }
    const { token } = await res.json();
    localStorage.setItem("token", token);
    window.location.href = "/";
  };

  return (
    <main style={{ maxWidth: 400, margin: "80px auto", fontFamily: "sans-serif", padding: "0 1rem" }}>
      <h1>{mode === "login" ? "ログイン" : "新規登録"}</h1>
      <form onSubmit={handleSubmit} style={{ display: "flex", flexDirection: "column", gap: "1rem" }}>
        <input
          type="email"
          placeholder="メールアドレス"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
          style={{ padding: "0.5rem", fontSize: "1rem" }}
        />
        <input
          type="password"
          placeholder="パスワード"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
          style={{ padding: "0.5rem", fontSize: "1rem" }}
        />
        {error && <p style={{ color: "red" }}>{error}</p>}
        <button type="submit" style={{ padding: "0.5rem", fontSize: "1rem", cursor: "pointer" }}>
          {mode === "login" ? "ログイン" : "登録"}
        </button>
      </form>
      <p style={{ marginTop: "1rem" }}>
        {mode === "login" ? "アカウントをお持ちでない方は" : "既にアカウントをお持ちの方は"}
        <button
          onClick={() => setMode(mode === "login" ? "register" : "login")}
          style={{ background: "none", border: "none", color: "blue", cursor: "pointer", textDecoration: "underline" }}
        >
          {mode === "login" ? "新規登録" : "ログイン"}
        </button>
      </p>
    </main>
  );
}
