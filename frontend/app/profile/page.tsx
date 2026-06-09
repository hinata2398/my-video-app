"use client";

import { useEffect, useRef, useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";

type Me = {
  id: number;
  email: string;
  username: string;
  avatar_url: string;
};

export default function ProfilePage() {
  const router = useRouter();
  const [me, setMe] = useState<Me | null>(null);
  const [username, setUsername] = useState("");
  const [avatarURL, setAvatarURL] = useState("");
  const [avatarPreview, setAvatarPreview] = useState<string | null>(null);
  const [uploading, setUploading] = useState(false);
  const [saving, setSaving] = useState(false);
  const [message, setMessage] = useState<{ type: "ok" | "err"; text: string } | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  // 自分のプロフィール取得
  useEffect(() => {
    const token = localStorage.getItem("token");
    if (!token) {
      router.push("/auth");
      return;
    }
    fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/me`, {
      headers: { Authorization: `Bearer ${token}` },
    })
      .then((res) => {
        if (!res.ok) throw new Error();
        return res.json();
      })
      .then((data: Me) => {
        setMe(data);
        setUsername(data.username);
        setAvatarURL(data.avatar_url);
      })
      .catch(() => router.push("/auth"));
  }, []);

  // アバター画像を選んだとき
  const handleAvatarChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    // プレビュー表示
    setAvatarPreview(URL.createObjectURL(file));
    setUploading(true);

    try {
      const token = localStorage.getItem("token")!;

      // 1. Presigned URL を取得
      const urlRes = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/api/me/avatar-upload-url`,
        { headers: { Authorization: `Bearer ${token}` } },
      );
      const { upload_url, avatar_url } = await urlRes.json();

      // 2. MinIO に直接アップロード
      await fetch(upload_url, {
        method: "PUT",
        headers: { "Content-Type": file.type },
        body: file,
      });

      // 3. アップロード先URLを保存
      setAvatarURL(avatar_url);
    } catch {
      setMessage({ type: "err", text: "画像のアップロードに失敗しました" });
    } finally {
      setUploading(false);
    }
  };

  // 保存
  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    setMessage(null);

    try {
      const token = localStorage.getItem("token")!;
      const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/me`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ username, avatar_url: avatarURL }),
      });
      if (!res.ok) throw new Error();
      setMessage({ type: "ok", text: "プロフィールを更新しました！" });
    } catch {
      setMessage({ type: "err", text: "更新に失敗しました" });
    } finally {
      setSaving(false);
    }
  };

  if (!me) return <main style={{ padding: "2rem" }}><p style={{ color: "#888" }}>読み込み中...</p></main>;

  const displayAvatar = avatarPreview || avatarURL;

  return (
    <main style={{ maxWidth: 480, margin: "0 auto", padding: "2rem" }}>
      <Link href="/" style={{ color: "#888", textDecoration: "none" }}>
        ← トップへ戻る
      </Link>

      <h1 style={{ marginTop: "1.5rem", marginBottom: "2rem", fontSize: "1.4rem" }}>
        プロフィール編集
      </h1>

      <form onSubmit={handleSave}>
        {/* アバター */}
        <div style={{ display: "flex", flexDirection: "column", alignItems: "center", marginBottom: "2rem" }}>
          <div
            onClick={() => fileInputRef.current?.click()}
            style={{
              width: 96,
              height: 96,
              borderRadius: "50%",
              background: "#ddd",
              overflow: "hidden",
              cursor: "pointer",
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              fontSize: "2.5rem",
              fontWeight: "bold",
              color: "#888",
              border: "2px solid #eee",
              position: "relative",
            }}
          >
            {displayAvatar ? (
              <img
                src={displayAvatar}
                alt="avatar"
                style={{ width: "100%", height: "100%", objectFit: "cover" }}
              />
            ) : (
              (username || me.email)[0].toUpperCase()
            )}
            {/* ホバー時オーバーレイ */}
            <div style={{
              position: "absolute", inset: 0, background: "rgba(0,0,0,0.35)",
              display: "flex", alignItems: "center", justifyContent: "center",
              borderRadius: "50%", opacity: 0,
              transition: "opacity 0.2s",
            }}
              onMouseEnter={(e) => (e.currentTarget.style.opacity = "1")}
              onMouseLeave={(e) => (e.currentTarget.style.opacity = "0")}
            >
              <span style={{ color: "#fff", fontSize: "0.75rem" }}>変更</span>
            </div>
          </div>
          <input
            ref={fileInputRef}
            type="file"
            accept="image/*"
            style={{ display: "none" }}
            onChange={handleAvatarChange}
          />
          {uploading && (
            <p style={{ color: "#888", fontSize: "0.8rem", marginTop: 8 }}>アップロード中...</p>
          )}
          <p style={{ color: "#aaa", fontSize: "0.8rem", marginTop: 8 }}>
            クリックして画像を変更
          </p>
        </div>

        {/* メールアドレス（読み取り専用） */}
        <div style={{ marginBottom: "1.25rem" }}>
          <label style={{ display: "block", fontSize: "0.85rem", color: "#555", marginBottom: 6 }}>
            メールアドレス
          </label>
          <input
            type="text"
            value={me.email}
            disabled
            style={{
              width: "100%",
              padding: "10px 12px",
              border: "1px solid #e0e0e0",
              borderRadius: 8,
              fontSize: "0.95rem",
              background: "#f9f9f9",
              color: "#999",
              boxSizing: "border-box",
            }}
          />
        </div>

        {/* ユーザー名 */}
        <div style={{ marginBottom: "1.5rem" }}>
          <label style={{ display: "block", fontSize: "0.85rem", color: "#555", marginBottom: 6 }}>
            ユーザー名
          </label>
          <input
            type="text"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            placeholder="ユーザー名を入力"
            style={{
              width: "100%",
              padding: "10px 12px",
              border: "1px solid #ddd",
              borderRadius: 8,
              fontSize: "0.95rem",
              boxSizing: "border-box",
            }}
          />
        </div>

        {/* メッセージ */}
        {message && (
          <div style={{
            marginBottom: "1rem",
            padding: "10px 14px",
            borderRadius: 8,
            background: message.type === "ok" ? "#e8f5e9" : "#ffebee",
            color: message.type === "ok" ? "#2e7d32" : "#c62828",
            fontSize: "0.9rem",
          }}>
            {message.text}
          </div>
        )}

        {/* 保存ボタン */}
        <button
          type="submit"
          disabled={saving || uploading}
          style={{
            width: "100%",
            padding: "12px",
            background: saving || uploading ? "#ccc" : "#e00",
            color: "#fff",
            border: "none",
            borderRadius: 8,
            fontSize: "1rem",
            cursor: saving || uploading ? "not-allowed" : "pointer",
          }}
        >
          {saving ? "保存中..." : "保存する"}
        </button>
      </form>
    </main>
  );
}
