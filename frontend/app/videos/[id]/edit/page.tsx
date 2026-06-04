"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";

export default function EditVideoPage() {
  const { id } = useParams();
  const router = useRouter();
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/videos/${id}`)
      .then((res) => res.json())
      .then((data) => {
        setTitle(data.title);
        setDescription(data.description);
      })
      .finally(() => setLoading(false));
  }, [id]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    setError("");
    const token = localStorage.getItem("token");
    if (!token) { router.push("/auth"); return; }

    const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/videos/${id}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json", Authorization: `Bearer ${token}` },
      body: JSON.stringify({ title, description, thumbnail_url: "", video_url: "" }),
    });
    setSaving(false);
    if (!res.ok) { setError(await res.text()); return; }
    router.push(`/videos/${id}`);
  };

  if (loading) return <main style={{ padding: "2rem" }}>読み込み中...</main>;

  return (
    <main style={{ maxWidth: 600, margin: "0 auto", padding: "2rem" }}>
      <h1>動画を編集</h1>
      <form onSubmit={handleSubmit} style={{ display: "flex", flexDirection: "column", gap: "1rem" }}>
        <div>
          <label style={{ display: "block", marginBottom: "0.25rem" }}>タイトル *</label>
          <input
            type="text" value={title} onChange={(e) => setTitle(e.target.value)} required
            style={{ width: "100%", padding: "0.5rem", fontSize: "1rem", boxSizing: "border-box" }}
          />
        </div>
        <div>
          <label style={{ display: "block", marginBottom: "0.25rem" }}>説明</label>
          <textarea
            value={description} onChange={(e) => setDescription(e.target.value)} rows={4}
            style={{ width: "100%", padding: "0.5rem", fontSize: "1rem", boxSizing: "border-box" }}
          />
        </div>
        {error && <p style={{ color: "red" }}>{error}</p>}
        <div style={{ display: "flex", gap: "1rem" }}>
          <button type="submit" disabled={saving} style={{
            padding: "0.75rem 1.5rem", background: saving ? "#ccc" : "#e00",
            color: "#fff", border: "none", borderRadius: 4, cursor: saving ? "not-allowed" : "pointer", fontSize: "1rem",
          }}>
            {saving ? "保存中..." : "保存する"}
          </button>
          <button type="button" onClick={() => router.push(`/videos/${id}`)} style={{
            padding: "0.75rem 1.5rem", background: "none", border: "1px solid #ccc",
            borderRadius: 4, cursor: "pointer", fontSize: "1rem",
          }}>
            キャンセル
          </button>
        </div>
      </form>
    </main>
  );
}
